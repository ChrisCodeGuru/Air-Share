package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
)

// Delete Account API
func DeleteAccount(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var Account types.DelAccount
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
	} else if Account.Del == "" {
		w.WriteHeader(206)
		w.Write([]byte("Please fill in all the information"))
		return
	} else if strings.EqualFold(("delete " + username), Account.Del) {
		_, err := db.Exec("DELETE FROM account WHERE id = $1", JWTClaims.UUID)
		if err != nil {
			logging.ErrorLogger(err, r)
			fmt.Println(err)
		}

		accountDetails := map[string]interface{}{
			"Email": email,
		}
		accountDetailsJSON, err := json.Marshal(accountDetails)
		if err != nil {
			logging.ErrorLogger(err, r)
			panic(err)
		}
		logging.Logger(string(accountDetailsJSON), 200, r)

		w.WriteHeader(200)
	} else {
		w.WriteHeader(206)
		w.Write([]byte("Incorrect Input"))
	}
}
