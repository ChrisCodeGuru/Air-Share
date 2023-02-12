package token

import (
	"fmt"
	"net/http"
)

func ExtractToken(r *http.Request) (string, error) {

	// get token contents
	token, err := r.Cookie("Token")
	if err != nil {
		// handle extraction error
		fmt.Println("no token")
	}

	// return JWT string
	return token.Value, err
}
