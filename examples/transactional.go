package main

import (
	"context"
	"fmt"
	"os"

	"github.com/customerio/go-customerio"
)

func main() {

	ctx := context.Background()

	client := customerio.NewAPIClient("<your-key-here>", customerio.WithRegion(customerio.RegionUS))

	req := customerio.SendEmailRequest{
		Identifiers: map[string]string{
			"id": "customer_1",
		},
		To:      "customer@example.com",
		From:    "business@example.com",
		Subject: "hello, {{ trigger.name }}",
		Body:    "hello from the Customer.io {{ trigger.client }} client",
		MessageData: map[string]interface{}{
			"client": "Go",
			"name":   "gopher",
		},
	}

	f, err := os.Open("<path to file>")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := req.Attach("sample.pdf", f); err != nil {
		panic(err)
	}

	resp, err := client.SendEmail(ctx, &req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
