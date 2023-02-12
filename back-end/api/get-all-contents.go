package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
)

func AllContents(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

	if len(r.Header["Referer"]) == 0 || r.Header["Referer"][0] != "https://localhost/fileshare" || len(r.Header["X-Csrf-Token"]) == 0 {
		w.WriteHeader(401)
		fmt.Println("CSRF")
		return
	} else if t.After(time.Date(yyyy, fullmmmm, dd, hhhh, mm, ss, 0, time.UTC)) || csrfToken == "" || r.Header["X-Csrf-Token"][0] != csrfToken {
		w.WriteHeader(401)
		fmt.Println("CSRF")
		return
	}

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
	var file_id string
	var file_name string
	var file_permission int
	var file_parent_id sql.NullString
	// var folder_id string
	// var folder_name string
	// var folder_permission int
	// var folder_parent_id sql.NullString

	myslice := []string{}

	// Get all files user have access to
	filerows, err := db.Query("SELECT file.id, file.name, account_file.permission, file.folder_id FROM account_file INNER JOIN file ON file.id = account_file.file_id WHERE account_file.account_id::text = $1", JWTClaims.UUID)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Close query result
	defer filerows.Close()

	// Iterate through the rows
	for filerows.Next() {
		// Extract row data
		err := filerows.Scan(&file_id, &file_name, &file_permission, &file_parent_id)
		if err != nil {
			fmt.Println(err)
			return
		}

		folderTraverse := []string{}
		filepath := []string{}

		// Get all folders the files are in
		for {
			if file_parent_id.String != "" {

				var file_directory_name string
				var file_directory_permission int

				file_current_id := file_parent_id.String

				fileDirectoryErr := db.QueryRow("SELECT folder.name, account_folder.permission, folder.parent_id FROM account_folder INNER JOIN folder ON folder.id = account_folder.folder_id WHERE account_folder.account_id::text = $1 AND folder.id = $2", JWTClaims.UUID, file_parent_id.String).Scan(&file_directory_name, &file_directory_permission, &file_parent_id)

				if fileDirectoryErr != nil {
					fmt.Println(fileDirectoryErr)
					// return
					break
				}

				fileDetails := map[string]interface{}{
					// "Count":			numberCount[0],
					"Name":          file_directory_name,
					"ID":            file_current_id,
					"Parent Folder": file_parent_id.String,
					"Permission":    file_directory_permission,
					"Type":          "Folder",
					"Owner":         userEmail,
				}
				fileDetailsJSON, err := json.Marshal(fileDetails)
				if err != nil {
					panic(err)
				}

				filepath = append([]string{string(file_directory_name)}, filepath...)
				folderTraverse = append([]string{string(fileDetailsJSON)}, folderTraverse...)
			} else {

				break
			}
		}
		fileDetails := map[string]interface{}{
			"Document Name":  file_name,
			"File ID":        file_id,
			"Parent Folder":  file_parent_id,
			"Permission":     file_permission,
			"Type":           "File",
			"User Email":     userEmail,
			"File Directory": folderTraverse,
			"File Path":      strings.Join(filepath, "/") + "/" + file_name,
		}
		fileDetailsJSON, err := json.Marshal(fileDetails)
		if err != nil {
			// logging.ErrorLogger(err, r)
			panic(err)
		}

		myslice = append(myslice, string(fileDetailsJSON))
	}

	// // Get all folders user have access to
	// folderrows, err := db.Query("SELECT folder.id, folder.name, account_folder.permission, folder.parent_id  FROM account_folder INNER JOIN folder ON folder.id = account_folder.folder_id WHERE account_folder.account_id::text = $1", JWTClaims.UUID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // Close query result
	// defer folderrows.Close()

	// // Iterate through the rows
	// for folderrows.Next() {
	// 	// Extract row data
	// 	err := folderrows.Scan(&folder_id, &folder_name, &folder_permission, &folder_parent_id)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	folderTraverse := []string{}
	// 	folderpath := []string{}

	// 	// for {
	// 	// 	if folder_parent_id.String != "" {

	// 	// 		var folder_directory_name string
	// 	// 		var folder_directory_permission int

	// 	// 		folder_current_id := folder_parent_id.String

	// 	// 		folderDirectoryErr := db.QueryRow("SELECT folder.name, account_folder.permission, folder.parent_id FROM account_folder INNER JOIN folder ON folder.id = account_folder.folder_id WHERE account_folder.account_id::text = $1 AND folder.id = $2", JWTClaims.UUID, folder_parent_id.String).Scan(&folder_directory_name, &folder_directory_permission, &file_parent_id)

	// 	// 		if folderDirectoryErr != nil {
	// 	// 			fmt.Println(folderDirectoryErr)
	// 	// 			return
	// 	// 		}

	// 	// 		fileDetails := map[string]interface{}{
	// 	// 			"Document Name":    folder_directory_name,
	// 	// 			"File ID":    		folder_current_id,
	// 	// 			"Parent Folder":	folder_parent_id.String,
	// 	// 			"Permission":   	folder_directory_permission,
	// 	// 			"Type": 			"Folder",
	// 	// 			"User Email": 		userEmail,
	// 	// 		}
	// 	// 		fileDetailsJSON, err := json.Marshal(fileDetails)
	// 	// 		if err != nil {
	// 	// 			panic(err)
	// 	// 		}

	// 	// 		folderpath = append(folderpath, folder_directory_name)
	// 	// 		folderTraverse = append(folderTraverse, string(fileDetailsJSON))
	// 	// 	} else {
	// 	// 		folderDetails := map[string]interface{}{
	// 	// 			"Document Name":    folder_name,
	// 	// 			"File ID":    		folder_id,
	// 	// 			"Parent Folder":	folder_parent_id,
	// 	// 			"Permission":    	folder_permission,
	// 	// 			"Type": 			"Folder",
	// 	// 			"User Email": 		userEmail,
	// 	// 			"File Directory":	folderTraverse,
	// 	// 			"File Path":		strings.Join(folderpath,"\\")+"\\"+folder_name,
	// 	// 		}
	// 	// 		folderDetailsJSON, err := json.Marshal(folderDetails)
	// 	// 		if err != nil {
	// 	// 			panic(err)
	// 	// 		}

	// 	// 		myslice = append(myslice, string(folderDetailsJSON))
	// 	// 		break
	// 	// 	}
	// 	// }
	// 	folderDetails := map[string]interface{}{
	// 		"Document Name":    folder_name,
	// 		"File ID":    		folder_id,
	// 		"Parent Folder":	folder_parent_id,
	// 		"Permission":    	folder_permission,
	// 		"Type": 			"Folder",
	// 		"User Email": 		userEmail,
	// 		"File Directory":	folderTraverse,
	// 		"File Path":		strings.Join(folderpath,"\\")+"\\"+folder_name,
	// 	}
	// 	folderDetailsJSON, err := json.Marshal(folderDetails)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	myslice = append(myslice, string(folderDetailsJSON))
	// }

	// Sort the Slice
	sort.Sort(sort.StringSlice(myslice))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(myslice)
}
