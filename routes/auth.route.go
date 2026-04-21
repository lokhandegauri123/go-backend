package routes

import (
	"go-backend/handlers"
	"go-backend/middleware"
	"net/http"
)

func AuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)

	mux.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.Profile)))
}
