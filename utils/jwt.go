package utils

import (
	// "go/token"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// var SECRET_KEY = []byte("mysecretkey")

func GenerateToken(userId string, email string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(SECRET_KEY)
}	
