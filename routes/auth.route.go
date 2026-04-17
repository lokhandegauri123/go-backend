package routes

import(
	"go-backend/handlers"
	"net/http"
)

func AuthRoutes(mux *http.ServeMux){
	mux.HandleFunc("/register",handlers.Register)
}