package auth

import (
	"errors"
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

func RefreshToken(token, secret string) (string, error) {
	issuer, id, err := ValidateJWT(token, secret)
	if err != nil {
		return "", err
	}

	if issuer != "chirpy-refresh" {
		return "", errors.New("not a refresh token")
	}

	accessToken, err := MakeRefresh(id, token)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
