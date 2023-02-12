package api

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/cryptography"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"github.com/gorilla/mux"
)

func DownloadFile(w http.ResponseWriter, r *http.Request, db *sql.DB) {

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
		w.Write([]byte("you are not authorised to download this file"))
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

	if len(r.Header["Referer"]) == 0 || r.Header["Referer"][0] != "https://localhost/fileshare" || len(r.Header["X-Csrf-Token"]) == 0 {
		w.WriteHeader(401)
		fmt.Println("CSRF")
		return
	} else if t.After(time.Date(yyyy, fullmmmm, dd, hhhh, mm, ss, 0, time.UTC)) || csrfToken == "" || r.Header["X-Csrf-Token"][0] != csrfToken {
		w.WriteHeader(401)
		fmt.Println("CSRF")
		return
	}

	// Check if file exists
	var fileName string
	var encryptedKey []byte
	var sensitive bool
	var asymmetric_key string

	counterr := db.QueryRow("SELECT name, encrypted_key, asymmetric_key, sensitive FROM file WHERE id=$1", file_id).Scan(&fileName, &encryptedKey, &asymmetric_key, &sensitive)
	if counterr != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(counterr)
		return
	}

	if fileName == "" {
		// throw 404 not found
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("file not found"))
		return
	}

	// Check if user has permissions to download
	var permission int

	filepermissionerr := db.QueryRow("SELECT permission FROM account_file WHERE account_id=$1 AND file_id=$2", JWTClaims.UUID, file_id).Scan(&permission)
	if filepermissionerr != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(filepermissionerr)
		return
	}

	if permission == 0 {
		// throw 401 unauthorised
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("you are not authorised to download this file"))
		return
	}

	// Read file
	var path string
	if permission == 4 || !sensitive {
		path = fmt.Sprintf("objects/%s", file_id)
	} else {
		path = fmt.Sprintf("objects/%s-censored", file_id)
	}
	encryptedFileData, err := os.ReadFile(path)
	if err != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// decrypt AES encryption key
	aesKey, err := cryptography.Decrypt(encryptedKey, asymmetric_key)
	if err != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// create AES cipher
	newcipher, err := aes.NewCipher(aesKey.Plaintext)
	if err != nil {
		// throw 500 interal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	gcm, err := cipher.NewGCM(newcipher)
	if err != nil {
		// throw 500 interal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	nonceSize := gcm.NonceSize()
	if len(encryptedFileData) < nonceSize {
		// throw 500 interal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	nonce, ciphertext := encryptedFileData[:nonceSize], encryptedFileData[nonceSize:]

	// decrypt file
	decryptedFileData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	var email string
	accounterr := db.QueryRow("SELECT get_user_email($1)", JWTClaims.UUID).Scan(&email)
	if accounterr != nil {
		logging.ErrorLogger(accounterr, r)
		fmt.Println(accounterr)
	}

	fileDetails := map[string]interface{}{
		"Email":     email,
		"File Name": fileName,
	}
	fileDetailsJSON, err := json.Marshal(fileDetails)
	if err != nil {
		logging.ErrorLogger(err, r)
		panic(err)
	}
	logging.Logger(string(fileDetailsJSON), 201, r)

	// send file data
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(decryptedFileData)
}
