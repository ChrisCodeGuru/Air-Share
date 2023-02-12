package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"github.com/gorilla/mux"
)

func DeleteFile(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Get route variables
	file_id := mux.Vars(r)["id"]
	// Validate parent id is acceptable
	if file_id != "root" {
		matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, file_id)
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
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("you are not authorised to delete this file"))
		return
	}

	// Get JWT claims
	JWTClaims, err := token.Claims(cookie_token)
	if err != nil {
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
		return
	} else if t.After(time.Date(yyyy, fullmmmm, dd, hhhh, mm, ss, 0, time.UTC)) || csrfToken == "" || r.Header["X-Csrf-Token"][0] != csrfToken {
		w.WriteHeader(401)
		fmt.Println("CSRF")
		return
	}

	// Check if file exists
	var count int

	counterr := db.QueryRow("SELECT COUNT(*) FROM file WHERE id=$1", file_id).Scan(&count)
	if counterr != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(counterr)
		return
	}

	if count != 1 {
		// throw 404 not found
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("file not found"))
		return
	}

	// Check if user has permissions to delete the file (is owner)
	var permission int

	filepermissionerr := db.QueryRow("SELECT permission FROM account_file WHERE account_id=$1 AND file_id=$2", JWTClaims.UUID, file_id).Scan(&permission)
	if filepermissionerr != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(filepermissionerr)
		return
	}

	if permission != 4 {
		// throw 401 unauthorised
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("you are not authorised to delete this file"))
		return
	}

	// Delete file from database
	var filename string

	deleteerr := db.QueryRow("SELECT delete_file($1)", file_id).Scan(&filename)
	if deleteerr != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(deleteerr)
		return
	}

	// Delete file from directory
	directorydeleteerr := os.Remove(fmt.Sprintf("objects/%s", file_id))
	if directorydeleteerr != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(directorydeleteerr)
		return
	}

	// Delete censored file from directory
	os.Remove(fmt.Sprintf("objects/%s-censored", file_id))

	var email string
	accounterr := db.QueryRow("SELECT email FROM account WHERE id = $1", JWTClaims.UUID).Scan(&email)
	if accounterr != nil {
		logging.ErrorLogger(accounterr, r)
		fmt.Println(accounterr)
	}

	fileDetails := map[string]interface{}{
		"Email":     email,
		"File Name": filename,
	}
	fileDetailsJSON, err := json.Marshal(fileDetails)
	if err != nil {
		logging.ErrorLogger(err, r)
		panic(err)
	}
	logging.Logger(string(fileDetailsJSON), 201, r)

	// return 200 ok
	w.WriteHeader(http.StatusOK)
}
