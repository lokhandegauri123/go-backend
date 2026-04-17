package handlers

import (
	"context"
	"encoding/json"

	// "go/token"
	"log"
	"net/http"
	"time"

	"go-backend/config"
	"go-backend/models"
	"go-backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", 405)
		return
	}

	// request body read

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	collection := config.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// insert into MongoDB
	result, err := collection.InsertOne(ctx, user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// taking MongoDB generated ID
	userId := result.InsertedID.(primitive.ObjectID)

	idStr := userId.Hex()

	// generate token
	token, err := utils.GenerateToken(idStr, user.Email)
	if err != nil {
		log.Println("TOKEN ERROR:", err) // 🔥 add this
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setting token into cookie
	cookie := &http.Cookie{
		Name : "token",
		Value: token,
		HttpOnly: true,
		Path: "/",
		MaxAge: 7600,
	}
	http.SetCookie(w, cookie)

	// sending response
	response := map[string]interface{}{
		"message": "user Registered successfully",
		"user_Id": idStr,
		"token":   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
