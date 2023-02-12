package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/email"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/argon2"
)

func Signup(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var newUser types.User
	json.NewDecoder(r.Body).Decode(&newUser)

	// Check if Email or Username is taken
	var existing string

	userdetailserr := db.QueryRow("SELECT check_existing_user($1, $2)", newUser.Email, newUser.Username).Scan(&existing)
	if userdetailserr != nil {
		fmt.Println(userdetailserr)
	}

	if strings.Split(existing[1:len(existing)-1], ",")[0] != "0" {
		fmt.Println("email is taken:", newUser.Email)
		w.WriteHeader(409)
		w.Write([]byte("Email is taken"))
		return
	}

	if strings.Split(existing[1:len(existing)-1], ",")[1] != "0" {
		fmt.Println("username is taken:", newUser.Username)
		w.WriteHeader(409)
		w.Write([]byte("Username is taken"))
		return
	}

	if strings.Contains(newUser.Username, "@") {
		// throw 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username cannot contain '@'"))
		fmt.Println("username cannot contain '@'")
		return
	}

	if len(r.Header["Origin"]) == 0 || r.Header["Origin"][0] != "https://localhost" || len(r.Header["Referer"]) == 0 || r.Header["Referer"][0] != "https://localhost/signup" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("CSRF")
		return
	}

	// Check if payload is valid
	var re1 = regexp.MustCompile(`^[a-zA-Z ]+$`)
	var re2 = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
	_, emailErr := mail.ParseAddress(newUser.Email)

	if newUser.Fname == "" || newUser.Lname == "" || newUser.Username == "" || newUser.Email == "" || newUser.Password == "" || newUser.Cpassword == "" {
		w.WriteHeader(206)
		w.Write([]byte("Please fill in all the information"))
		return
	} else if !re1.MatchString(newUser.Fname) || !re1.MatchString(newUser.Lname) {
		w.WriteHeader(206)
		w.Write([]byte("Name cannot contain Special Characters or Numbers"))
		return
	} else if !re2.MatchString(newUser.Username) {
		w.WriteHeader(206)
		w.Write([]byte("Username cannot contain Special Characters"))
		return
	} else if emailErr != nil {
		w.WriteHeader(206)
		w.Write([]byte("Invalid Email."))
		return
	} else if len(newUser.Password) < 8 {
		w.WriteHeader(206)
		w.Write([]byte("Passwords must be at least 8 Characters long"))
		return
	} else if newUser.Password != newUser.Cpassword {
		w.WriteHeader(206)
		w.Write([]byte("Passwords Do Not Match"))
		return
	}

	// Create 32 characters salt
	salt := make([]byte, 32)
	for i := 0; i < 32; i++ {
		charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQESTUVWXYZ1234567890"
		// insert ASCII code into Array
		salt[i] = charset[rand.Intn(len(charset))]
	}

	// Hash the password+salt (password, salt, iterations/execution time, memory, threads/parallelism, keyLength)
	var memory uint32 = 64 * 1024
	var iterations uint32 = 2
	var parallelism uint8 = 4

	argon2Hash := argon2.IDKey([]byte(newUser.Password), salt, iterations, memory, parallelism, 32)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Argon2Hash := base64.RawStdEncoding.EncodeToString(argon2Hash)
	hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, iterations, parallelism, b64Salt, b64Argon2Hash)
	fmt.Println(hash)

	// Generate verification code
	code := randstr.String(20)
	verification_code := email.Encode(code)

	res, err := db.Exec("INSERT INTO account(email, username, password, verified, verification_code, fname, lname) VALUES($1, $2, $3, $4, $5, $6, $7)", newUser.Email, newUser.Username, hash, false, verification_code, newUser.Fname, newUser.Lname)
	if err != nil {
		logging.ErrorLogger(err, r)
		fmt.Println(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		logging.ErrorLogger(err, r)
		fmt.Println(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		logging.ErrorLogger(err, r)
		fmt.Println(err)
	}

	fmt.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	// Send verification email
	emailData := types.EmailData{
		URL:       "https://localhost" + "/api/verify/" + code,
		FirstName: newUser.Fname,
	}
	email.SendEmail(newUser.Email, &emailData)

	userDetails := map[string]interface{}{
		"fname":    newUser.Fname,
		"lname":    newUser.Lname,
		"username": newUser.Username,
		"email":    newUser.Email,
	}
	userDetailsJSON, err := json.Marshal(userDetails)
	if err != nil {
		logging.ErrorLogger(err, r)
		panic(err)
	}
	logging.Logger(string(userDetailsJSON), 201, r)
	w.WriteHeader(201)
}
