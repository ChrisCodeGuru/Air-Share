package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
)

func Authentication(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

	// Get password state and otp verification
	var content types.Authentication
	var returns string

	authenticationerr := db.QueryRow("SELECT account_authentication_check($1)", JWTClaims.UUID).Scan(&returns)
	if authenticationerr != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(authenticationerr)
		return
	}

	switch returns {
	case "(t,t)":
		content.Password, content.OTP = true, true
	case "(t,f)":
		content.Password, content.OTP = true, false
	case "(f,t)":
		content.Password, content.OTP = false, true
	case "(f,f)":
		content.Password, content.OTP = false, false
	}

	// Return content
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}
