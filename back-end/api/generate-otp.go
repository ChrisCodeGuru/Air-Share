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

func GenerateOTP(w http.ResponseWriter, r *http.Request, db *sql.DB) {

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

	// Generate TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "fileshare.com",
		AccountName: JWTClaims.UUID,
		SecretSize:  15,
	})
	if err != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update database with key data
	res, err := db.Exec("UPDATE account SET otp_secret=$1, otp_auth_url=$2 WHERE id=$3", key.Secret(), key.URL(), JWTClaims.UUID)
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

	otp_response := types.Otp_response{
		Base32:      key.Secret(),
		Otpauth_url: key.URL(),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(otp_response)
}
