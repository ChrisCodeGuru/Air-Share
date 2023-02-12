package token

import (
	"net/http"
	"time"
)

func GenerateTokenCookie(token string) http.Cookie {
	cookie := http.Cookie{}
	cookie.Name = "Token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(365 * 24 * time.Hour)
	cookie.Secure = false
	cookie.HttpOnly = true
	cookie.Path = "/"

	return cookie
}
