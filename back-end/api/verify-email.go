package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/email"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/gorilla/mux"
)

func VerifyEmail(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Get route variables
	code := mux.Vars(r)["code"]

	// Encode code to get verification code
	verfication_code := email.Encode(code)

	// Get uuid for verification code
	var uuid string

	accountverificationerr := db.QueryRow("SELECT id FROM account WHERE verification_code=$1", verfication_code).Scan(&uuid)
	if accountverificationerr != nil {
		fmt.Println(accountverificationerr)
	}

	// If user not found
	if uuid == "" {
		// throw 404 conflict
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("something went wrong"))
		return
	}

	// Update user verification
	res, err := db.Exec("UPDATE account SET verified=$1, verification_code=$2 WHERE id=$3", true, sql.NullString{}, uuid)
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

	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
}
