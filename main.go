package main

import (
	"fmt"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func main() {
	client := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    "APCA-API-KEY-ID",
		APISecret: "APCA-API-SECRET-KEY",
		BaseURL:   "https://paper-api.alpaca.markets",
	})
	acct, err := client.GetAccount()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *acct)
}
