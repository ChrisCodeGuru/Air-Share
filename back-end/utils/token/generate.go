package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(uuid string) (string, error) {

	// adding claims to JWT token
	atClaims := jwt.MapClaims{}
	atClaims["uuid"] = uuid                                 // setting user UUID claim
	atClaims["exp"] = time.Now().Add(time.Hour * 12).Unix() // setting expiry time to 15 minutes from token generation time

	// JWT token signing configuration
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	// signing JWT token
	token, err := at.SignedString(sampleSecretKey)
	// handle signing errors
	if err != nil {
		return "", err
	}

	// return token as string
	return token, nil
}
