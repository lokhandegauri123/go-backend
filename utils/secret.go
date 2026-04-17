package utils

import(
	"crypto/rand"
	"encoding/hex"
	"log"
)

var SECRET_KEY []byte

func InitSecret(){
	b:=make([]byte,32)
	_, err := rand.Read(b)
	if err != nil{
		log.Fatal("failed to generate secret key")

	}

	SECRET_KEY = []byte(hex.EncodeToString(b))
	log.Println("SECRET GENERATED:", string(SECRET_KEY))
}