package main

import (
	// "fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"go-backend/config"
	"go-backend/routes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, loading environment variables from OS")
	}

	//connect DB
	config.ConnectDB()

	// create router (mux)
	mux := http.NewServeMux()

	// register routes
	routes.AuthRoutes(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server running with MongoDB 🚀"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
	
}
