package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// load evn file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	API_Key := os.Getenv("APCA_API_KEY_ID")
	API_Secret := os.Getenv("APCA_API_SECRET_KEY")

	if API_Key == "" || API_Secret == "" {
		log.Fatal("Not set in the enviroment variables")
	}

	client := &http.Client{}

	start := time.Now()
	resp, err := client.Get("https://paper-api.alpaca.markets/v2/account")
	duration :=
		time.Since(start).Milliseconds()

	if err != nil {
		fmt.Printf("Error:  %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response time: %v, Status: %d\n", duration, resp.StatusCode)
}
