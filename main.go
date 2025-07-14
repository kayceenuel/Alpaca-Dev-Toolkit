package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	alpacaAPI_Key := os.Getenv("APCA-API-KEY-ID")
	alpacaSecret_Key := os.Getenv("APCA-API-SECRET-KEY")
}
