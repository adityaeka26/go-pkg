package helper

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(privateKey string, claims jwt.MapClaims) (*string, error) {
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(1 * 24 * time.Hour).Unix()

	privateKey = strings.ReplaceAll(privateKey, "\\n", "\n")

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return nil, err
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
