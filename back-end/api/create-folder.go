package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
)

func CreateFolder(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var newFolder types.Folder
	json.NewDecoder(r.Body).Decode(&newFolder)

	// Validate parent id is acceptable
	if newFolder.Parent != "root" {
		matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, newFolder.Parent)
		if !matched {
			// throw 400 bad request
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("request parameters are incorrect"))
			return
		}
	}

	// Extract JWT from cookie
	cookie_token, err := token.ExtractToken(r)
	if err != nil {
		// throw 401 error if no token
		logging.ErrorLogger(err, r)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("you are not authorised"))
		return
	}

	// Get JWT claims
	JWTClaims, err := token.Claims(cookie_token)
	if err != nil {
		logging.ErrorLogger(err, r)
		fmt.Println(err)
		return
	}

	var csrfreturns string
	t := time.Now().UTC()

	tokenerr := db.QueryRow("SELECT csrf_check($1)", JWTClaims.UUID).Scan(&csrfreturns)
	if tokenerr != nil {
		logging.ErrorLogger(tokenerr, r)
		fmt.Println(tokenerr)
	}

	csrfToken := strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[0]
	yyyy, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[1])
	mmmm, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[2])
	dd, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[3])
	hhhh, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[4])
	mm, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[5])
	ss, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[6])
	fullmmmm := time.Month(mmmm)

	if len(r.Header["Origin"]) == 0 || r.Header["Origin"][0] != "https://localhost" || len(r.Header["Referer"]) == 0 || r.Header["Referer"][0] != "https://localhost/fileshare" || len(r.Header["X-Csrf-Token"]) == 0 {
		w.WriteHeader(401)
		fmt.Println("CSRF")
	} else if t.After(time.Date(yyyy, fullmmmm, dd, hhhh, mm, ss, 0, time.UTC)) || csrfToken == "" || r.Header["X-Csrf-Token"][0] != csrfToken {
		w.WriteHeader(401)
		fmt.Println("CSRF")
	}

	if newFolder.Parent == "root" {
		// Initialise variable to hold new folder id
		var newFolderUUID string

		createerr := db.QueryRow("SELECT create_folder($1, $2)", JWTClaims.UUID, newFolder.Name).Scan(&newFolderUUID)
		if createerr != nil {
			// Return 500 Server Internal Error
			logging.ErrorLogger(createerr, r)
			fmt.Println(createerr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var email string
		accounterr := db.QueryRow("SELECT get_user_email($1)", JWTClaims.UUID).Scan(&email)
		if accounterr != nil {
			logging.ErrorLogger(accounterr, r)
			fmt.Println(accounterr)
		}

		folderDetails := map[string]interface{}{
			"Email":         email,
			"Folder":        newFolder.Name,
			"Parent Folder": newFolder.Parent,
		}
		folderDetailsJSON, err := json.Marshal(folderDetails)
		if err != nil {
			logging.ErrorLogger(err, r)
			panic(err)
		}
		logging.Logger(string(folderDetailsJSON), 201, r)

		// Return 201 Created
		w.WriteHeader(http.StatusCreated)
	} else {

		// Check if directory exists
		var count int

		existenceerr := db.QueryRow("SELECT check_folder_existence($1)", newFolder.Parent).Scan(&count)
		if existenceerr != nil {
			fmt.Println(existenceerr)
		}

		if count == 0 {
			// throw 404 not found
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("folder not found"))
			return
		}

		// Check if user has permissions to directory
		var permisisonLevel int

		permissionerr := db.QueryRow("SELECT permission FROM account_folder WHERE account_id=$1 AND folder_id=$2", JWTClaims.UUID, newFolder.Parent).Scan(&permisisonLevel)
		if permissionerr != nil {
			// throw 500 internal server error
			w.WriteHeader(http.StatusInternalServerError)
			logging.ErrorLogger(permissionerr, r)
			fmt.Println(permissionerr)
			return
		}

		if permisisonLevel < 2 {
			// throw 401 unauthorised
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("you are not authorised to make changes to this folder"))
			return
		}

		// Initialise variable to hold new folder id
		var newFolderUUID string

		createerr := db.QueryRow("SELECT create_folder($1, $2, $3)", JWTClaims.UUID, newFolder.Name, newFolder.Parent).Scan(&newFolderUUID)
		if createerr != nil {
			// Return 500 Server Internal Error
			logging.ErrorLogger(createerr, r)
			fmt.Println(createerr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var email string
		accounterr := db.QueryRow("SELECT get_user_email($1)", JWTClaims.UUID).Scan(&email)
		if accounterr != nil {
			logging.ErrorLogger(accounterr, r)
			fmt.Println(accounterr)
		}

		folderDetails := map[string]interface{}{
			"Email":         email,
			"Folder Name":   newFolder.Name,
			"Parent Folder": newFolder.Parent,
		}
		folderDetailsJSON, err := json.Marshal(folderDetails)
		if err != nil {
			logging.ErrorLogger(err, r)
			panic(err)
		}
		logging.Logger(string(folderDetailsJSON), 201, r)

		// Return 201 Created
		w.WriteHeader(http.StatusCreated)
	}
}
