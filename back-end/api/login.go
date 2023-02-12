package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"golang.org/x/crypto/argon2"
)

// Login API
func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var loginUser types.Login
	json.NewDecoder(r.Body).Decode(&loginUser)

	doubleSubmitToken, cookieErr := r.Cookie("DoubleSubmitToken")
	if cookieErr != nil {
		fmt.Println(cookieErr)
	}

	var re1 = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)

	_, emailErr := mail.ParseAddress(loginUser.Username)

	if len(r.Header["X-Csrf-Token"]) == 0 || cookieErr != nil {
		w.WriteHeader(401)
		fmt.Println("CSRF")
	} else if len(r.Header["X-Csrf-Token"]) != 0 && cookieErr == nil {
		mac := hmac.New(sha256.New, []byte("JIvygVYT*Y*GTY{YGVGfvtF&cFC&CF&T"))
		mac.Write([]byte(r.Header["X-Csrf-Token"][0]))
		messageMAC := mac.Sum(nil)
		originalMAC, err := base64.StdEncoding.DecodeString(doubleSubmitToken.Value)
		if err != nil {
			logging.ErrorLogger(err, r)
			fmt.Println(err)
		}

		if !hmac.Equal(messageMAC, originalMAC) {
			w.WriteHeader(401)
			fmt.Println("CSRF")
		}
	} else if loginUser.Username == "" || loginUser.Password == "" {
		w.WriteHeader(206)
		w.Write([]byte("Please fill in all the information"))
		return
	} else if strings.Contains(loginUser.Username, "@") && emailErr != nil {
		w.WriteHeader(206)
		w.Write([]byte("Incorrect Username or Password"))
		return
	} else if !strings.Contains(loginUser.Username, "@") && !re1.MatchString(loginUser.Username) {
		w.WriteHeader(206)
		w.Write([]byte("Incorrect Username or Password"))
		return
	}

	// Check database for user
	var uuid string
	var hashed_password string
	var verified bool
	var mfaenabled bool

	passworderr := db.QueryRow("SELECT id, password, verified, otp_enabled FROM account WHERE email=$1 OR username=$1;", loginUser.Username).Scan(&uuid, &hashed_password, &verified, &mfaenabled)
	if passworderr != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		logging.ErrorLogger(passworderr, r)
		fmt.Println(passworderr)
		return
	}

	// If account does not exist
	if uuid == "" {
		loginDetails := map[string]interface{}{
			"loginID": loginUser.Username,
		}
		loginDetailsJSON, err := json.Marshal(loginDetails)
		if err != nil {
			logging.ErrorLogger(err, r)
			panic(err)
		}
		logging.Logger(string(loginDetailsJSON), 401, r)
		fmt.Println("account does not exist")
		w.WriteHeader(401)
		w.Write([]byte("Incorrect Username or Password"))
		return
	}

	// If account is not verified
	if !verified {
		// throw 401 not authorised
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please verify email first"))
		return
	}

	// Verify Hash
	vals := strings.Split(hashed_password, "$")
	// fmt.Println("vals", vals)
	// Check Version
	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		logging.ErrorLogger(err, r)
		fmt.Println(err)
		return
	} else if version != argon2.Version {
		fmt.Println("Incorrect Version")
		return
	}

	// Get input parameters
	var memory uint32
	var iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		logging.ErrorLogger(err, r)
		fmt.Println(err)
		return
	}
	// fmt.Println(memory, iterations, parallelism)

	// Get salt & og hash
	originalArgon2Hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		logging.ErrorLogger(err, r)
		fmt.Println(err)
		return
	}
	salt2, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		logging.ErrorLogger(err, r)
		fmt.Println(err)
		return
	}

	// Create New Hash
	newArgon2Hash := argon2.IDKey([]byte(loginUser.Password), salt2, iterations, memory, parallelism, 32)
	// newb64Salt := base64.RawStdEncoding.EncodeToString(salt2)
	// newb64Argon2Hash := base64.RawStdEncoding.EncodeToString(newArgon2Hash)
	// newhash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, 64*1024, 1, 4, newb64Salt, newb64Argon2Hash)

	// vals2 := strings.Split(newhash, "$")
	// newArgon2Hash2, err := base64.RawStdEncoding.Strict().DecodeString(vals2[5])

	// fmt.Println("ReceivedHash:", RetrievedHash)
	// fmt.Println("NewHash:", newhash)
	// fmt.Println(subtle.ConstantTimeCompare(originalArgon2Hash, newArgon2Hash2))
	if subtle.ConstantTimeCompare(originalArgon2Hash, newArgon2Hash) == 1 {

		// redirect to 2fa page if mfa is enabled
		if mfaenabled {
			w.WriteHeader(http.StatusPartialContent)
			return
		}

		jwt, err := token.GenerateJWT(uuid)
		if err != nil {
			logging.ErrorLogger(err, r)
			fmt.Println(err)
			w.WriteHeader(301)
			return
		}

		// Create token cookie
		cookie := token.GenerateTokenCookie(jwt)

		csrfToken := make([]byte, 32)
		for i := 0; i < 32; i++ {
			charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQESTUVWXYZ1234567890"
			// insert ASCII code into Array
			csrfToken[i] = charset[rand.Intn(len(charset))]
		}

		t := time.Now().UTC()
		n := t.AddDate(0, 0, 1)
		dateTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second())

		var totalCount int
		tokenerr := db.QueryRow("SELECT COUNT(*) FROM csrf WHERE user_id = $1", uuid).Scan(&totalCount)
		if tokenerr != nil {
			logging.ErrorLogger(tokenerr, r)
			fmt.Println(tokenerr)
		}

		if totalCount >= 1 {
			_, err := db.Exec("UPDATE csrf SET csrf_token = $1, expire_datetime = $2 WHERE user_id = $3", string(csrfToken), dateTime, uuid)
			if err != nil {
				logging.ErrorLogger(err, r)
				fmt.Println(err)
			}
		} else {
			_, err := db.Exec("INSERT INTO csrf (user_id, csrf_token, expire_datetime) VALUES($1, $2, $3)", uuid, string(csrfToken), dateTime)
			if err != nil {
				logging.ErrorLogger(err, r)
				fmt.Println(err)
			}
		}

		// Set csrf cookie
		csrfcookie := http.Cookie{}
		csrfcookie.Name = "csrf"
		csrfcookie.Value = string(csrfToken)
		csrfcookie.Expires = time.Now().Add(365 * 24 * time.Hour)
		csrfcookie.Secure = false
		csrfcookie.Path = "/"

		loginDetails := map[string]interface{}{
			"loginID": loginUser.Username,
		}
		loginDetailsJSON, err := json.Marshal(loginDetails)
		if err != nil {
			logging.ErrorLogger(err, r)
			panic(err)
		}
		logging.Logger(string(loginDetailsJSON), 200, r)

		http.SetCookie(w, &csrfcookie)
		http.SetCookie(w, &cookie)
		w.WriteHeader(200)
		return
	}

	loginDetails := map[string]interface{}{
		"loginID": loginUser.Username,
	}
	loginDetailsJSON, err := json.Marshal(loginDetails)
	if err != nil {
		logging.ErrorLogger(err, r)
		panic(err)
	}
	logging.Logger(string(loginDetailsJSON), 401, r)

	fmt.Println("Incorrect Username or Password")
	w.WriteHeader(401)
	w.Write([]byte("Incorrect Username or Password"))
}
