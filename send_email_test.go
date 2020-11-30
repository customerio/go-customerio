package customerio_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/customerio/go-customerio"
)

func TestSendEmail(t *testing.T) {
	emailRequest := &customerio.SendEmailRequest{
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

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
		}
		defer req.Body.Close()

		var body customerio.SendEmailRequest
		if err := json.Unmarshal(b, &body); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(&body, emailRequest) {
			t.Errorf("Request differed, want: %#v, got: %#v", emailRequest, body)
		}

		w.Write([]byte(`{
			"delivery_id": "ABCDEFG",
			"queued_at": 1500111111,
			"recipient": "12345"
		  }`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	resp, err := api.SendEmail(context.Background(), emailRequest)
	if err != nil {
		t.Error(err)
	}

	expect := &customerio.SendEmailResponse{
		TransactionalResponse: customerio.TransactionalResponse{
			DeliveryID: "ABCDEFG",
			QueuedAt:   time.Unix(1500111111, 0),
			Recipient:  "12345",
		},
	}

	if !reflect.DeepEqual(resp, expect) {
		t.Errorf("Expect: %#v, Got: %#v", expect, resp)
	}
}

func TestSendEmailError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(502)
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	resp, err := api.SendEmail(context.Background(), &customerio.SendEmailRequest{
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
	})
	if err == nil {
		t.Errorf("Expected error, got: %#v", resp)
	}

	if e, ok := err.(*customerio.TransactionalError); !ok {
		t.Errorf("Expected TransactionalError, got: %#v", e)
	}
}
