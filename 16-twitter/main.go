package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	pathToEnv = "./16-twitter/.env"
	consumerKey = "CONSUMER_KEY"
	consumerSecret = "CONSUMER_SECRET"
)

func main() {
	if err := godotenv.Load(pathToEnv); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	key := os.Getenv(consumerKey)
	secret := os.Getenv(consumerSecret)
	fmt.Println(key, secret)
}
