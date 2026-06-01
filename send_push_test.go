package customerio_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/customerio/go-customerio/v3"
)

func TestSendPush(t *testing.T) {
	pushRequest := &customerio.SendPushRequest{
		Identifiers: map[string]string{
			"id": "customer_1",
		},
		To:      "customer@example.com",
		Title:   "hello, {{ trigger.name }}",
		Message: "hello from the Customer.io {{ trigger.client }} client",
		MessageData: map[string]any{
			"client": "Go",
			"name":   "gopher",
		},
	}
	d, err := customerio.NewDevice("device-id", "ios", map[string]any{"attr1": "value1"})
	if err != nil {
		t.Error(err)
	}
	pushRequest.Device = d

	var verify = func(request []byte) {
		var body customerio.SendPushRequest
		if err := json.Unmarshal(request, &body); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(&body, pushRequest) {
			t.Errorf("Request differed, want: %#v, got: %#v", request, body)
		}
	}

	api, srv := transactionalServer(t, verify)
	defer srv.Close()

	resp, err := api.SendPush(context.Background(), pushRequest)
	if err != nil {
		t.Error(err)
	}

	expect := &customerio.SendPushResponse{
		TransactionalResponse: customerio.TransactionalResponse{
			DeliveryID: testDeliveryID,
			QueuedAt:   time.Unix(int64(testQueuedAt), 0),
		},
	}

	if !reflect.DeepEqual(resp, expect) {
		t.Errorf("Expect: %#v, Got: %#v", expect, resp)
	}
}

func TestSendPushError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(502)
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	resp, err := api.SendPush(context.Background(), &customerio.SendPushRequest{
		Identifiers: map[string]string{
			"id": "customer_1",
		},
		To:      "customer@example.com",
		Title:   "hello",
		Message: "hello from Go",
	})
	if err == nil {
		t.Errorf("Expected error, got: %#v", resp)
	}

	if _, ok := err.(*customerio.TransactionalError); !ok {
		t.Errorf("Expected TransactionalError, got: %#v", err)
	}
}
