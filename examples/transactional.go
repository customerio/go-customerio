package main

import (
	"context"
	"fmt"
	"os"

	"github.com/customerio/go-customerio/v3"
)

func main() {

	ctx := context.Background()

	client := customerio.NewAPIClient("<your-key-here>", customerio.WithRegion(customerio.RegionUS))

	emailReq := customerio.SendEmailRequest{
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

	if err := emailReq.Attach("sample.pdf", f); err != nil {
		panic(err)
	}

	emailResp, err := client.SendEmail(ctx, &emailReq)
	if err != nil {
		panic(err)
	}

	fmt.Println(emailResp)

	// Create and configure your transactional messages in the Customer.io UI.
	transactionalMessageID := "push_message_id"

	pushReq := customerio.SendPushRequest{
		TransactionalMessageID: transactionalMessageID,
		Identifiers: map[string]string{
			"id": "customer_1",
		},
		Title:   "hello, {{ trigger.name }}",
		Message: "hello from the Customer.io {{ trigger.client }} client",
	}

	pushResp, err := client.SendPush(ctx, &pushReq)
	if err != nil {
		panic(err)
	}

	fmt.Println(pushResp)
}
