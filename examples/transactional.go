package main

import (
	"context"
	"fmt"

	"github.com/customerio/go-customerio"
)

func main() {

	ctx := context.Background()

	client := customerio.NewAPIClient("<your-key-here>")

	req := customerio.SendEmailRequest{
		CustomerID: "customer_1",
		To:         "customer@example.com",
		From:       "business@example.com",
		Subject:    "hello, {{ trigger.name }}",
		Body:       "hello from the Customer.io {{ trigger.client }} client",
		MessageData: map[string]interface{}{
			"client": "Go",
			"name":   "gopher",
		},
	}

	resp, err := client.SendEmail(ctx, &req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
