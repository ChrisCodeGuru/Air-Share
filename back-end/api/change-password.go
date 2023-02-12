package api

import (
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"golang.org/x/crypto/argon2"
)

// Change Password API
func ChangePassword(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var Passwords types.Passwords
	json.NewDecoder(r.Body).Decode(&Passwords)

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

	if len(r.Header["Origin"]) == 0 || r.Header["Origin"][0] != "https://localhost" || len(r.Header["Referer"]) == 0 || r.Header["Referer"][0] != "https://localhost/settings" || len(r.Header["X-Csrf-Token"]) == 0 {
		w.WriteHeader(401)
		fmt.Println("CSRF")
	} else if t.After(time.Date(yyyy, fullmmmm, dd, hhhh, mm, ss, 0, time.UTC)) || csrfToken == "" || r.Header["X-Csrf-Token"][0] != csrfToken {
		w.WriteHeader(401)
		fmt.Println("CSRF")
	} else if Passwords.OPassword == "" || Passwords.NPassword == "" || Passwords.CPassword == "" {
		w.WriteHeader(206)
		w.Write([]byte("Please fill in all the information"))
		return
	} else if len(Passwords.NPassword) < 8 {
		w.WriteHeader(206)
		w.Write([]byte("Passwords must be at least 8 Characters long"))
		return
	} else if Passwords.NPassword != Passwords.CPassword {
		w.WriteHeader(206)
		w.Write([]byte("Passwords Do Not Match"))
		return
	} else {
		// Check if user is using Google Login (no password)
		var userCount int
		counterr := db.QueryRow("SELECT account_password_check($1)", JWTClaims.UUID).Scan(&userCount)
		if counterr != nil {
			logging.ErrorLogger(counterr, r)
			fmt.Println(counterr)
		}
		if userCount == 0 {
			w.WriteHeader(401)
			w.Write([]byte("Incorrect Password"))
			return
		}

		// Get database email & password
		var email string
		var hashed_password string

		passworderr := db.QueryRow("SELECT email, password FROM account WHERE id = $1;", JWTClaims.UUID).Scan(&email, &hashed_password)
		if passworderr != nil {
			logging.ErrorLogger(passworderr, r)
			fmt.Println(passworderr)
		}

		// Verify Hash
		vals := strings.Split(hashed_password, "$")

		// Check Version
		var version int
		_, err := fmt.Sscanf(vals[2], "v=%d", &version)
		if err != nil {
			logging.ErrorLogger(err, r)
			fmt.Println(err)
		} else if version != argon2.Version {
			fmt.Println("Incorrect Version")
		}

		// Get input parameters
		var memory uint32
		var iterations uint32
		var parallelism uint8
		_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
		if err != nil {
			logging.ErrorLogger(err, r)
			fmt.Println(err)
		}

		// Get salt & og hash
		originalArgon2Hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
		if err != nil {
			logging.ErrorLogger(err, r)
			fmt.Println(err)
		}

		salt2, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
		if err != nil {
			logging.ErrorLogger(err, r)
			fmt.Println(err)
		}

		// Create New Hash
		newArgon2Hash := argon2.IDKey([]byte(Passwords.OPassword), salt2, iterations, memory, parallelism, 32)

		if subtle.ConstantTimeCompare(originalArgon2Hash, newArgon2Hash) == 1 {
			// Hash Password
			var memory uint32 = 64 * 1024
			var iterations uint32 = 2
			var parallelism uint8 = 4

			argon2Hash := argon2.IDKey([]byte(Passwords.NPassword), salt2, iterations, memory, parallelism, 32)
			b64Salt := base64.RawStdEncoding.EncodeToString(salt2)
			b64Argon2Hash := base64.RawStdEncoding.EncodeToString(argon2Hash)
			hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, iterations, parallelism, b64Salt, b64Argon2Hash)

			_, err := db.Exec("UPDATE account SET password = $1 WHERE id = $2", hash, JWTClaims.UUID)
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

			fmt.Println("Change Password: Success")
			w.WriteHeader(200)
			return
		} else {
			accountDetails := map[string]interface{}{
				"Email": email,
			}
			accountDetailsJSON, err := json.Marshal(accountDetails)
			if err != nil {
				logging.ErrorLogger(err, r)
				panic(err)
			}
			logging.Logger(string(accountDetailsJSON), 401, r)

			fmt.Println("Change Password: Incorrect Password")
			w.WriteHeader(401)
			w.Write([]byte("Incorrect Password"))
		}
	}
}
