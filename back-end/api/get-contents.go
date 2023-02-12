package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"github.com/gorilla/mux"
)

func Contents(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Get route variables
	directory := mux.Vars(r)["directory"]
	// Validate directory is accepted
	if directory != "root" && directory != "shared" {
		matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, directory)
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
		w.Write([]byte("you are not authorised"))
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

	if directory == "root" {

		// Get user email
		var userEmail string
		emailerr := db.QueryRow("SELECT email FROM account WHERE id=$1", JWTClaims.UUID).Scan(&userEmail)
		if emailerr != nil {
			// throw 500 internal server error
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(emailerr)
			return
		}

		// Temporary variables
		var returns string
		var Content types.Content

		// Get user root folder contents folders
		folderrows, err := db.Query("SELECT get_owned_contents_folder($1)", JWTClaims.UUID)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Close query result
		defer folderrows.Close()

		// Iterate throuhg the rows
		for folderrows.Next() {
			// Extract row data
			err := folderrows.Scan(&returns)
			if err != nil {
				fmt.Println(err)
				return
			}
			// Append to content variable
			Content.Folders = append(
				Content.Folders,
				types.FolderObject{
					ID:         strings.Split(returns[1:len(returns)-1], ",")[0],
					Name:       strings.Split(returns[1:len(returns)-1], ",")[1],
					Permission: 4,
				},
			)
		}

		// Check for errors after interating
		err = folderrows.Err()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Get user root folder contents files
		filerows, err := db.Query("SELECT get_owned_contents_file($1)", JWTClaims.UUID)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Close query result
		defer filerows.Close()

		// Iterate throuhg the rows
		for filerows.Next() {
			// Extract row data
			err := filerows.Scan(&returns)
			if err != nil {
				fmt.Println(err)
				return
			}
			// Append to content variable
			Content.Files = append(
				Content.Files,
				types.FileObject{
					ID:         strings.Split(returns[1:len(returns)-1], ",")[0],
					Name:       strings.Split(returns[1:len(returns)-1], ",")[1],
					Permission: 4,
					Sensitive:  strings.Split(returns[1:len(returns)-1], ",")[2],
					Hash:       strings.Split(returns[1:len(returns)-1], ",")[3],
				},
			)
		}

		// Check for errors after interating
		err = filerows.Err()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Append permission level of directory
		Content.Permission = 4

		// Return content
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Content)

	} else if directory == "shared" {

		// Temporary variables
		var returns string
		var Content types.Content

		// Get user root folder contents folders
		folderrows, err := db.Query("SELECT get_shared_contents_folder($1)", JWTClaims.UUID)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Close query result
		defer folderrows.Close()

		// Iterate throuhg the rows
		for folderrows.Next() {
			// Extract row data
			err := folderrows.Scan(&returns)
			if err != nil {
				fmt.Println(err)
				return
			}
			permission, _ := strconv.Atoi(strings.Split(returns[1:len(returns)-1], ",")[2])

			// Append to content variable
			Content.Folders = append(
				Content.Folders,
				types.FolderObject{
					ID:         strings.Split(returns[1:len(returns)-1], ",")[0],
					Name:       strings.Split(returns[1:len(returns)-1], ",")[1],
					Permission: permission,
				},
			)
		}

		// Check for errors after interating
		err = folderrows.Err()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Get user root folder contents files
		filerows, err := db.Query("SELECT get_shared_contents_file($1)", JWTClaims.UUID)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Close query result
		defer filerows.Close()

		// Iterate throuhg the rows
		for filerows.Next() {
			// Extract row data
			err := filerows.Scan(&returns)
			if err != nil {
				fmt.Println(err)
				return
			}
			permission, _ := strconv.Atoi(strings.Split(returns[1:len(returns)-1], ",")[2])

			// Append to content variable
			Content.Files = append(
				Content.Files,
				types.FileObject{
					ID:         strings.Split(returns[1:len(returns)-1], ",")[0],
					Name:       strings.Split(returns[1:len(returns)-1], ",")[1],
					Permission: permission,
					Sensitive:  strings.Split(returns[1:len(returns)-1], ",")[3],
					Hash:       strings.Split(returns[1:len(returns)-1], ",")[4],
				},
			)
		}

		// Check for errors after interating
		err = filerows.Err()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Append permission level of directory
		Content.Permission = 0

		// Return content
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Content)

	} else {

		// Check if folder exists
		var count int

		existenceerr := db.QueryRow("SELECT COUNT(*) FROM folder WHERE id=$1", directory).Scan(&count)
		if existenceerr != nil {
			fmt.Println(existenceerr)
		}

		if count == 0 {
			// throw 404 not found
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("folder not found"))
			return
		}

		// Check if user has permissions to folder
		var permissionLevel int

		permissionerr := db.QueryRow("SELECT permission FROM account_folder WHERE account_id=$1 AND folder_id=$2", JWTClaims.UUID, directory).Scan(&permissionLevel)
		if permissionerr != nil {
			fmt.Println(permissionerr)
		}

		if permissionLevel == 0 {
			// throw 401 unauthorised
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("you are not authorised to view this folder"))
			return
		}

		// Temporary variables
		var returns string
		var Content types.Content

		// Get user root folder contents folders
		folderrows, err := db.Query("SELECT get_contents_folder($1, $2)", directory, JWTClaims.UUID)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Close query result
		defer folderrows.Close()

		// Iterate throuhg the rows
		for folderrows.Next() {
			// Extract row data
			err := folderrows.Scan(&returns)
			if err != nil {
				fmt.Println(err)
				return
			}
			permission, _ := strconv.Atoi(strings.Split(returns[1:len(returns)-1], ",")[2])

			// Append to content variable
			Content.Folders = append(
				Content.Folders,
				types.FolderObject{
					ID:         strings.Split(returns[1:len(returns)-1], ",")[0],
					Name:       strings.Split(returns[1:len(returns)-1], ",")[1],
					Permission: permission,
				},
			)
		}

		// Check for errors after interating
		err = folderrows.Err()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Get user root folder contents files
		filerows, err := db.Query("SELECT get_contents_file($1, $2)", directory, JWTClaims.UUID)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Close query result
		defer filerows.Close()

		// Iterate throuhg the rows
		for filerows.Next() {
			// Extract row data
			err := filerows.Scan(&returns)
			if err != nil {
				fmt.Println(err)
				return
			}
			permission, _ := strconv.Atoi(strings.Split(returns[1:len(returns)-1], ",")[2])

			// Append to content variable
			Content.Files = append(
				Content.Files,
				types.FileObject{
					ID:         strings.Split(returns[1:len(returns)-1], ",")[0],
					Name:       strings.Split(returns[1:len(returns)-1], ",")[1],
					Permission: permission,
					Sensitive:  strings.Split(returns[1:len(returns)-1], ",")[3],
					Hash:       strings.Split(returns[1:len(returns)-1], ",")[4],
				},
			)
		}

		// Check for errors after interating
		err = filerows.Err()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Append permission level of directory
		Content.Permission = permissionLevel

		// Return content
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Content)
	}
}
