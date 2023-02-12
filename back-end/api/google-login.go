package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"golang.org/x/oauth2"
)

var (
	oauthStateString = "7A9LKthF1ZFhcpjb"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request, g *oauth2.Config) {
	url := g.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request, g *oauth2.Config, db *sql.DB) {
	content, err := GetGoogleUserInfo(r.FormValue("state"), r.FormValue("code"), g)
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Convert byte data to json
	var contentjson types.GoogleOAuthContent
	json.Unmarshal(content, &contentjson)

	var uuid string

	// Check if user exists in account database
	accountexistanceerr := db.QueryRow("SELECT id FROM account WHERE email=$1;", contentjson.Email).Scan(&uuid)
	if accountexistanceerr != nil {
		fmt.Println(accountexistanceerr)
	}

	// If user does not exist create user in database
	if uuid == "" {
		res, err := db.Exec("INSERT INTO account(email, username, verified) VALUES($1, $2, $3)", contentjson.Email, contentjson.ID, true)
		if err != nil {
			fmt.Println(err)
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			fmt.Println(err)
		}
		rowCnt, err := res.RowsAffected()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

		// Get user UUID
		accounterr := db.QueryRow("SELECT id FROM account WHERE email=$1;", contentjson.Email).Scan(&uuid)
		if accounterr != nil {
			fmt.Println(accounterr)
		}
	}

	// Generate JWT token
	jwt, err := token.GenerateJWT(uuid)
	if err != nil {
		// throw 500 internal server error
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create token cookie
	cookie := token.GenerateTokenCookie(jwt)

	// Create CSRF token
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
	csrfCookie := http.Cookie{}
	csrfCookie.Name = "csrf"
	csrfCookie.Value = string(csrfToken)
	csrfCookie.Expires = time.Now().Add(365 * 24 * time.Hour)
	csrfCookie.Secure = false
	csrfCookie.Path = "/"

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &csrfCookie)
	http.Redirect(w, r, "/fileshare", http.StatusPermanentRedirect)
}

func GetGoogleUserInfo(state string, code string, g *oauth2.Config) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := g.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}
