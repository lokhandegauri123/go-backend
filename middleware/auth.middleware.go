package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	// "golang.org/x/crypto/bcrypt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.Split(authHeader, " ")[1]

		PublicKeyData, err := os.ReadFile("certs/public.pem")
		if err != nil {
			http.Error(w, "error reading public key", 500)
			return
		}

		PublicKey, err := jwt.ParseRSAPublicKeyFromPEM(PublicKeyData)

		if err != nil {
			http.Error(w, "invalid public key", 500)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return PublicKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
