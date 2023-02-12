package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"github.com/pquerna/otp/totp"
)

func DeleteOTP(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var otpdata types.Otp_payload
	json.NewDecoder(r.Body).Decode(&otpdata)

	// Extract JWT from cookie
	cookie_token, err := token.ExtractToken(r)
	if err != nil {
		// throw 401 error if no token
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("you are not authorised"))
		return
	}

	// Get JWT claims
	JWTClaims, err := token.Claims(cookie_token)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get OTP secret
	var otp_secret string

	getotpsecreterr := db.QueryRow("SELECT otp_secret FROM account WHERE id=$1", JWTClaims.UUID).Scan(&otp_secret)
	if getotpsecreterr != nil {
		fmt.Println(getotpsecreterr)
	}

	// Validate OTP token
	valid := totp.Validate(otpdata.Token, otp_secret)
	if !valid {
		// throw 401
		fmt.Println(otpdata.Token, otp_secret)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("2FA code incorrect"))
		return
	}

	// Disable otp token
	res, err := db.Exec("UPDATE account SET otp_enabled=$1, otp_secret=NULL, otp_auth_url=NULL WHERE id=$2", false, JWTClaims.UUID)
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

	w.WriteHeader(http.StatusCreated)
}
