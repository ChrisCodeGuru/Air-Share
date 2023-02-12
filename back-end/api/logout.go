package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"github.com/go-redis/redis"
)

func Logout(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {

	// Extract JWT from cookie
	cookie_token, err := token.ExtractToken(r)
	if err != nil {
		// throw 401 error if no token
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("you are not authorised"))
		return
	}

	// Add JWT to redis
	added := rdb.Set(cookie_token, cookie_token, 0)
	if added.Err() != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Remove Authentication cookie
	cookie := http.Cookie{}
	cookie.Name = "Token"
	cookie.Value = "delete"
	cookie.Expires = time.Now()
	cookie.Secure = false
	cookie.Path = "/"

	// Remove CSRF cookie
	csrfCookie := http.Cookie{}
	csrfCookie.Name = "csrf"
	csrfCookie.Value = "delete"
	csrfCookie.Expires = time.Now()
	csrfCookie.Secure = false
	csrfCookie.Path = "/"

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &csrfCookie)
	w.WriteHeader(http.StatusOK)
}
