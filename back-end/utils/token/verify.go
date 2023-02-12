package token

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
)

func Verify(endpointHandler func(writer http.ResponseWriter, request *http.Request), rdb *redis.Client) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		// get token
		cookie_token, err := ExtractToken(request)
		if err != nil {
			// throw 401 error if no token found
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		// verify cookie is not blacklisted
		isBlacklisted := rdb.Exists(cookie_token)
		if isBlacklisted.Err() != nil {
			fmt.Println("error")
		}

		if isBlacklisted.Val() > 0 {
			// throw 401 unauthorised error if token is not valid
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
			if err != nil {
				return
			}
		}

		// verify token signature
		token, err := jwt.Parse(cookie_token,
			func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC) // define signing method as HMAC family

				// throw 401 unauthorised if token signature is invalid
				if !ok {
					fmt.Println("signature error")
					writer.WriteHeader(http.StatusUnauthorized)
					_, err := writer.Write([]byte("You're Unauthorized!"))
					if err != nil {
						return nil, err
					}
				}

				return sampleSecretKey, nil
			},
		)

		// throw 401 unauthorised error if value in cookie "Token" fails to parse
		// this could be due to a malformed cookie which is not in the correct JWT format
		if err != nil {
			fmt.Println("parsing error")
			writer.WriteHeader(http.StatusUnauthorized)
			_, err2 := writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
			if err2 != nil {
				return
			}
		}

		// check if token is valid
		if token.Valid {
			// parse request on to router if valid
			endpointHandler(writer, request)
		} else {
			// throw 401 unauthorised error if token is not valid
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
			if err != nil {
				return
			}
		}
	})
}
