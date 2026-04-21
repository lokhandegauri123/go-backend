package utils

import (
	// "go/token"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// var SECRET_KEY = []byte("mysecretkey")

func GenerateToken(userId string, email string) (string, error) {

	PrivateKeyData, _ := os.ReadFile("certs/private.pem")
	PrivateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(PrivateKeyData)

	claims := jwt.MapClaims{
		"userId": userId,
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(PrivateKey)
}
