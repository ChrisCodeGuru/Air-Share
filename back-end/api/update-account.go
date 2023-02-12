package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
)

// Delete Account API
func UpdateAccount(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var Account types.User
	json.NewDecoder(r.Body).Decode(&Account)

	// Extract JWT from cookie
	cookie_token, err := token.ExtractToken(r)
	if err != nil {
		// throw 401 error if no token
		w.WriteHeader(http.StatusUnauthorized)
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

	var email string
	var username string
	accounterr := db.QueryRow("SELECT email, username FROM account WHERE id = $1", JWTClaims.UUID).Scan(&email, &username)

	var re1 = regexp.MustCompile(`^[a-zA-Z ]+$`)
	var re2 = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
	_, emailErr := mail.ParseAddress(Account.Email)

	if accounterr != nil {
		logging.ErrorLogger(accounterr, r)
		fmt.Println(accounterr)
	}

	if len(r.Header["Origin"]) == 0 || r.Header["Origin"][0] != "https://localhost" || len(r.Header["Referer"]) == 0 || r.Header["Referer"][0] != "https://localhost/settings" || len(r.Header["X-Csrf-Token"]) == 0 {
		w.WriteHeader(401)
		fmt.Println("CSRF")
	} else if t.After(time.Date(yyyy, fullmmmm, dd, hhhh, mm, ss, 0, time.UTC)) || csrfToken == "" || r.Header["X-Csrf-Token"][0] != csrfToken {
		w.WriteHeader(401)
		fmt.Println("CSRF")
	} else if Account.Fname == "" || Account.Lname == "" || Account.Username == "" || Account.Email == "" {
		w.WriteHeader(206)
		w.Write([]byte("Please fill in all the information"))
		return
	} else if !re1.MatchString(Account.Fname) || !re1.MatchString(Account.Lname) {
		w.WriteHeader(206)
		w.Write([]byte("Please fill in all the information"))
		return
	} else if !re2.MatchString(Account.Username) {
		w.WriteHeader(206)
		w.Write([]byte("Username cannot contain Special Characters"))
		return
	} else if emailErr != nil {
		w.WriteHeader(206)
		w.Write([]byte("Invalid Email."))
		return
	} else {
		_, err := db.Exec("UPDATE account SET username = $1, fname = $2, lname = $3 WHERE id = $4 AND email = $5", Account.Username, Account.Fname, Account.Lname, JWTClaims.UUID, Account.Email)
		if err != nil {
			logging.ErrorLogger(err, r)
			fmt.Println(err)
		}

		accountDetails := map[string]interface{}{
			"fname":    Account.Fname,
			"lname":    Account.Lname,
			"username": Account.Username,
			"email":    Account.Email,
		}
		accountDetailsJSON, err := json.Marshal(accountDetails)
		if err != nil {
			logging.ErrorLogger(err, r)
			panic(err)
		}
		logging.Logger(string(accountDetailsJSON), 200, r)
		w.WriteHeader(200)
	}
}
