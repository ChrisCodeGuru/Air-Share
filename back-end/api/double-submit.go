package api

import (
	"fmt"
	"database/sql"
	"encoding/base64"
	"math/rand"
	"net/http"
	"time"
	"crypto/hmac"
	"crypto/sha256"
)

func DoubleSubmit(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if len(r.Header["Referer"]) == 0 || r.Header["Referer"][0] != "https://localhost/login" {
		w.WriteHeader(401)
		fmt.Println("CSRF")
		return
	}

	doubleSubmitToken := make([]byte, 32)
	for i := 0; i < 32; i++ {
		charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQESTUVWXYZ1234567890"
		doubleSubmitToken[i] = charset[rand.Intn(len(charset))]
	}

	mac := hmac.New(sha256.New, []byte("JIvygVYT*Y*GTY{YGVGfvtF&cFC&CF&T"))
	mac.Write([]byte(doubleSubmitToken))
	messageMAC := mac.Sum(nil)

	DoubleSubmitCookie := http.Cookie{}
	DoubleSubmitCookie.Name = "DoubleSubmitToken"
	DoubleSubmitCookie.Value = base64.StdEncoding.EncodeToString([]byte(messageMAC))
	DoubleSubmitCookie.Expires = time.Now().Add(365 * 24 * time.Hour)
	DoubleSubmitCookie.Secure = true
	DoubleSubmitCookie.Path = "/"

	http.SetCookie(w, &DoubleSubmitCookie)
	w.WriteHeader(201)
	w.Write([]byte(doubleSubmitToken))
}
