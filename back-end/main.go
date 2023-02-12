package main

import (
	"database/sql"
	_ "errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IsaacKoh88/infosecurity_project/back-end/api"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gorilla/mux"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://localhost/api/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// Define server routes and start server
func main() {
	// Open connection to postgresql server
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@ispj.c4pvcrvnysxe.ap-southeast-1.rds.amazonaws.com/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB")))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Configure connection to redis server
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379",
		Password: "",
		DB:       0,
	})

	r := mux.NewRouter()

	// Signup API
	r.HandleFunc("/api/createUser",
		func(w http.ResponseWriter, r *http.Request) {
			api.Signup(w, r, db)
		},
	).Methods("POST")

	// Verify account Email API
	r.HandleFunc("/api/verify/{code}",
		func(w http.ResponseWriter, r *http.Request) {
			api.VerifyEmail(w, r, db)
		},
	).Methods("GET")

	// Login API
	r.HandleFunc("/api/login",
		func(w http.ResponseWriter, r *http.Request) {
			api.Login(w, r, db)
		},
	).Methods("POST")

	// Login API
	r.HandleFunc("/api/login/mfa",
		func(w http.ResponseWriter, r *http.Request) {
			api.LoginMFA(w, r, db)
		},
	).Methods("POST")

	// OAuth Google Login
	r.HandleFunc("/api/google-login",
		func(w http.ResponseWriter, r *http.Request) {
			api.GoogleLogin(w, r, googleOauthConfig)
		},
	).Methods("GET")

	// OAuth Google Callback
	r.HandleFunc("/api/callback",
		func(w http.ResponseWriter, r *http.Request) {
			api.GoogleCallback(w, r, googleOauthConfig, db)
		},
	).Methods("GET")

	// Logout API
	r.HandleFunc("/api/logout",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.Logout(w, r, rdb)
			},
			rdb,
		),
	).Methods("GET")

	// Generating 2FA link and secret
	r.HandleFunc("/api/otp/generate",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.GenerateOTP(w, r, db)
			},
			rdb,
		),
	).Methods("GET")

	// Verifying 2FA token
	r.HandleFunc("/api/otp/verify",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.VerifyOTP(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Delete 2FA token
	r.HandleFunc("/api/otp/delete",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.DeleteOTP(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Get contents of a directory
	r.HandleFunc("/api/files/{directory}",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.Contents(w, r, db)
			},
			rdb,
		),
	).Methods("GET")

	// Get permissions of an object
	r.HandleFunc("/api/{objectType}/permissions/{id}",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.Permissions(w, r, db)
			},
			rdb,
		),
	).Methods("GET")

	// Handle Folder Creation
	r.HandleFunc("/api/create-folder",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.CreateFolder(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Handle Object Editing
	r.HandleFunc("/api/edit/{objectType}/{id}",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.ObjectEdit(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Handle Folder Deletion
	r.HandleFunc("/api/delete-folder/{id}",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.DeleteFolder(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Sharing & Permissions API
	r.HandleFunc("/api/share/{objectType}/{id}",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.Share(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Handle File Uploading API
	r.HandleFunc("/api/upload-file/{directory}",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.UploadFile(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Handle File Delete API
	r.HandleFunc("/api/delete-file/{id}",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.DeleteFile(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Handle File Download API
	r.HandleFunc("/api/download-file/{id}",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.DownloadFile(w, r, db)
			},
			rdb,
		),
	).Methods("GET")

	// Change Password API
	r.HandleFunc("/api/changePassword",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.ChangePassword(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Delete Account API
	r.HandleFunc("/api/deleteAccount",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.DeleteAccount(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Update Account API
	r.HandleFunc("/api/updateAccount",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.UpdateAccount(w, r, db)
			},
			rdb,
		),
	).Methods("POST")

	// Get Double Submit Cookie API
	r.HandleFunc("/api/doubleSubmit",
		func(w http.ResponseWriter, r *http.Request) {
			api.DoubleSubmit(w, r, db)
		},
	).Methods("GET")

	// Get Username API
	r.HandleFunc("/api/username",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.Username(w, r, db)
			},
			rdb,
		),
	).Methods("GET")

	// Get UseInfo API
	r.HandleFunc("/api/userinfo",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.UserInfo(w, r, db)
			},
			rdb,
		),
	).Methods("GET")

	// Get UseInfo API
	r.HandleFunc("/api/authentication",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.Authentication(w, r, db)
			},
			rdb,
		),
	).Methods("GET")

	r.HandleFunc("/api/contents",
		token.Verify(
			func(w http.ResponseWriter, r *http.Request) {
				api.AllContents(w, r, db)
			},
			rdb,
		),
	).Methods("GET")

	// Check if log folder exists
	_, err2 := os.Stat("./log")
	if err2 != nil {
		logFolder := os.Mkdir("./log", os.ModeDir)
		fmt.Println(logFolder)
	}

	// Check is error-log.log file exists
	_, err3 := os.Stat("./log/error-log.log")
	if err3 != nil {
		errorLogFile, _ := os.Create("./log/error-log.log")
		fmt.Println(errorLogFile)
	}

	// Check is general-log.log file exists
	_, err4 := os.Stat("./log/general-log.log")
	if err4 != nil {
		errorLogFile, _ := os.Create("./log/general-log.log")
		fmt.Println(errorLogFile)
	}

	fmt.Println("server is running on localhost:4000")

	log.Fatal(http.ListenAndServe("0.0.0.0:4000", r))
}
