package main

import (
	"fmt"

	"github.com/customerio/go-customerio"
)

func main() {

	client := customerio.NewTransactionalClient("<your-key-here>")

	email := customerio.Email{
		CustomerID: "customer_1",
		To:         "customer@example.com",
		From:       "business@example.com",
		Subject:    "hello world",
		Body:       "hello from the Customer.io Go Client",
	}

	resp, err := client.SendEmail(email)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
