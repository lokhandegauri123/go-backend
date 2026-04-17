package config

import (
	"context"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(os.Getenv("MONGO_URI")),
	)

	if err != nil {
		log.Fatal("Db connection error: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("DB not responding:", err)
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}
	if dbName == "" {
		parsedURI, parseErr := url.Parse(os.Getenv("MONGO_URI"))
		if parseErr == nil {
			dbName = strings.TrimPrefix(parsedURI.Path, "/")
		}
	}
	if dbName == "" {
		log.Fatal("Missing database name: set MONGO_DB or DB_NAME, or include it in MONGO_URI")
	}

	DB = client.Database(dbName)
	log.Println("DB connected", dbName)

}
