package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	API_Key := os.Getenv("APCA_API_KEY_ID")
	API_Secret := os.Getenv("APCA_API_SECRET_KEY")

	if API_Key == "" || API_Secret == "" {
		log.Fatal("API credentials not set in environment variables")
	}

	client := &http.Client{}

	// Create request with authentication headers
	req, err := http.NewRequest("GET", "https://paper-api.alpaca.markets/v2/account", nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Add required authentication headers
	req.Header.Set("APCA-API-KEY-ID", API_Key)
	req.Header.Set("APCA-API-SECRET-KEY", API_Secret)

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response body to see the exact error
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Printf("Response time: %vms, Status: %d\n", duration, resp.StatusCode)
	fmt.Printf("Response body: %s\n", string(body))
}
