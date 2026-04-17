package models

import (
	"go-backend/config"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

func InsertUser(user User) (primitive.ObjectID, error){
	collection := config.DB.Collection("users")

	result,err := collection.InsertOne(context.TODO(),user)
	if err!= nil{
		return primitive.NilObjectID,err
	}

	id := result.InsertedID.(primitive.ObjectID)

	return id ,nil
}






































































