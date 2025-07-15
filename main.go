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
		log.Fatal("API credentials not set in the enviroment variables")
	}

	client := &http.Client{}

	// create request authentication headers
	req, err := http.NewRequest("Get", "https://paper-api.alpaca.markets/v2/account", nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Add required authentication headers
	req.Header.Set("APCA-API-KEY-ID", API_Key)
	req.Header.Set("APCA-API-SECRET-KEY", API_Secret)

	start := time.Now()
	resp, err := client.Do(req)
	duration :=
		time.Since(start).Milliseconds()

	if err != nil {
		fmt.Printf("Error:  %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response time: %vms, Status: %d\n", duration, resp.StatusCode)
}
