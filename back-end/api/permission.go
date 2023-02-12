package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"github.com/gorilla/mux"
)

func Permissions(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Get route variables
	objectType := mux.Vars(r)["objectType"]
	id := mux.Vars(r)["id"]
	// Validate object type is acceptable
	if objectType != "file" && objectType != "folder" {
		// throw 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request parameters are incorrect"))
		return
	}
	// Validate UUID is valid format
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, id)
	if !matched {
		// throw 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request parameters are incorrect"))
		return
	}

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

	// Validate object exists
	var objectCount int

	counterr := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id=$1", objectType), id).Scan(&objectCount)
	if counterr != nil {
		fmt.Println(counterr)
	}

	if objectCount == 0 {
		// throw 404 not found
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("object not found")
		return
	}

	// Check if user has permissions to view permissions
	var permissionLevel int

	permissionerr := db.QueryRow(fmt.Sprintf("SELECT permission FROM account_%[1]s WHERE account_id=$1 AND %[1]s_id=$2", objectType), JWTClaims.UUID, id).Scan(&permissionLevel)
	if permissionerr != nil {
		fmt.Println(permissionerr)
	}

	if permissionLevel == 0 {
		// throw 401 unauthorised
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("unauthorised to view object permissions")
		return
	}

	// Temporary variables
	var returns string
	var accountsPermission types.AccountsPermission

	// Get object permissions
	accountrows, err := db.Query(fmt.Sprintf("SELECT get_%s_permission($1)", objectType), id)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Close query result
	defer accountrows.Close()

	// Iterate throuhg the rows
	for accountrows.Next() {
		// Extract row data
		err := accountrows.Scan(&returns)
		if err != nil {
			fmt.Println(err)
			return
		}
		permission, _ := strconv.Atoi(strings.Split(returns[1:len(returns)-1], ",")[2])

		// Append to content variable
		accountsPermission = append(
			accountsPermission,
			types.AccountPermission{
				ID:         strings.Split(returns[1:len(returns)-1], ",")[0],
				Email:      strings.Split(returns[1:len(returns)-1], ",")[1],
				Permission: permission,
			},
		)
	}

	// Check for errors after interating
	err = accountrows.Err()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Return accounts
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accountsPermission)
}
