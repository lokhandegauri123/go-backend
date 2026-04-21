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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		http.Error(w, "error hashing password", 500)
		return
	}

	user.Password = string(HashedPassword)

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
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7600,
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

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method is not post", 405)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	collection := config.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"email": user.Email,
	}

	result := collection.FindOne(ctx, filter)

	var foundUser models.User

	err = result.Decode(&foundUser)

	if err != nil {
		http.Error(w, "user not found", 404)
		return
	}

	// if user.Password != foundUser.Password {
	// 	http.Error(w, "invalid password", 401)
	// 	return
	// }
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "invalid password", 401)
	}

	token, err := utils.GenerateToken(foundUser.ID.Hex(), foundUser.Email)

	if err != nil {
		log.Println("Token error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   7600,
	}

	http.SetCookie(w, cookie)

	response := map[string]interface{}{
		"message": "user logged in successfully",
		"user_Id": foundUser.ID.Hex(),
		"token":   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", 405)
		return
	}

	user := r.Context().Value("user")

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "profile data retrieved",
		"user":    user,
	}
	json.NewEncoder(w).Encode(response)
}
