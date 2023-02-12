package token

import (
	"fmt"

	"github.com/IsaacKoh88/infosecurity_project/back-end/types"
	"github.com/golang-jwt/jwt/v4"
)

func Claims(tokenString string) (types.Token, error) {

	// verify token signature
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC) // define signing method as HMAC (HS256 is HMAC)

			// print error if token signature is invalid
			if !ok {
				return nil, fmt.Errorf("there's an error with the signing method")
			}
			return sampleSecretKey, nil
		},
	)

	// throw error if token fails to parse
	if err != nil {
		return types.Token{}, err
	}

	// extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uuid := claims["uuid"].(string) // extract uuid claim

		// return token claims
		return types.Token{
			UUID: uuid,
		}, nil
	}

	// if unable to extract claims
	return types.Token{}, fmt.Errorf("unable to extract claims")
}
