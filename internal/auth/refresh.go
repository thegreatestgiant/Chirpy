package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func MakeRefresh(userID int, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 24 * 60)),
		Subject:   fmt.Sprintf("%d", userID),
	})

	return token.SignedString([]byte(tokenSecret))
}
