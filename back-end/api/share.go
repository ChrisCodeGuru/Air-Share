package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"github.com/gorilla/mux"
)

func Share(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var shareParams types.Share
	json.NewDecoder(r.Body).Decode(&shareParams)

	// Get route variables
	objectType := mux.Vars(r)["objectType"]
	objectId := mux.Vars(r)["id"]
	// Validate object type is acceptable
	if objectType != "file" && objectType != "folder" {
		// throw 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request parameters are incorrect"))
		return
	}
	// Validate UUID is valid format
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, objectId)
	if !matched {
		// throw 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request parameters are incorrect"))
		return
	}
	// Validate permission level
	if shareParams.Permission > 4 || shareParams.Permission <= 0 {
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

	// Check if object exists
	var count int

	existenceerr := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id=$1", objectType), objectId).Scan(&count)
	if existenceerr != nil {
		fmt.Println(existenceerr)
	}

	if count == 0 {
		// throw 404 not found
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("object not found"))
		return
	}

	// Check if user has permissions to share the object (is adnimistrator or above)
	var permissionLevel int

	permissionerr := db.QueryRow(fmt.Sprintf("SELECT permission FROM account_%[1]s WHERE account_id=$1 AND %[1]s_id=$2", objectType), JWTClaims.UUID, objectId).Scan(&permissionLevel)
	if permissionerr != nil {
		fmt.Println(permissionerr)
	}

	if permissionLevel < 3 {
		// throw 401 unauthorised
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("you are not authorised to share this object"))
		return
	}

	// Check if target exisits
	var targetId string

	checkexistenceerr := db.QueryRow("SELECT id FROM account WHERE email=$1", shareParams.Email).Scan(&targetId)
	if checkexistenceerr != nil {
		fmt.Println(checkexistenceerr)
	}

	if targetId == "" {
		// throw 404 not found
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	// Check if target already has permissions to the file/folder
	var targetpermission int

	checktargetpermissionerr := db.QueryRow(fmt.Sprintf("SELECT permission FROM account_%[1]s WHERE account_id=$1 AND %[1]s_id=$2", objectType), targetId, objectId).Scan(&targetpermission)
	if checktargetpermissionerr != nil {
		fmt.Println(checktargetpermissionerr)
	}

	// Target already has the permissions
	if targetpermission == shareParams.Permission {
		// throw 409 conflict error
		fmt.Println(targetpermission, shareParams.Permission)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("user already has the requested permissions"))
		return
	}

	// Target is owner
	if targetpermission == 4 {
		// throw 403 forbidden
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("you cannot change the owner permissions"))
		return
	}

	// Target has more permissions
	if targetpermission > shareParams.Permission {
		var parentobjectpermission int

		checkparentobjecttargetpermissionerr := db.QueryRow(fmt.Sprintf("SELECT get_%s_parent_permission($1, $2)", objectType), targetId, objectId).Scan(&parentobjectpermission)
		if checkparentobjecttargetpermissionerr != nil {
			fmt.Println(checkparentobjecttargetpermissionerr)
		}

		// Parent object has more permission than requested permissions
		if parentobjectpermission > shareParams.Permission {
			// throw 403 forbidden
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("user has more permissions to parent folder"))
			return
		}
	}

	// Grant target permissions
	if shareParams.Permission != 0 {
		var success bool

		if targetpermission == 0 {
			createpermissionerr := db.QueryRow(fmt.Sprintf("SELECT create_%s_permission($1, $2, $3)", objectType), targetId, objectId, shareParams.Permission).Scan(&success)
			if createpermissionerr != nil {
				fmt.Println(createpermissionerr)
			}
		} else {
			createpermissionerr := db.QueryRow(fmt.Sprintf("SELECT edit_%s_permission($1, $2, $3)", objectType), targetId, objectId, shareParams.Permission).Scan(&success)
			if createpermissionerr != nil {
				fmt.Println(createpermissionerr)
			}
		}
	} else {
		var success bool

		createpermissionerr := db.QueryRow(fmt.Sprintf("SELECT delete_%s_permission($1, $2)", objectType), objectId).Scan(&success)
		if createpermissionerr != nil {
			fmt.Println(createpermissionerr)
		}
	}

	w.WriteHeader(201)
}
