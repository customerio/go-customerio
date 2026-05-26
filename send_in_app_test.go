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

func TestSendInApp(t *testing.T) {
	req := &customerio.SendInAppRequest{
		TransactionalMessageID: "123456",
		Identifiers: map[string]string{
			"id": "customer_1",
		},
		MessageData: map[string]any{
			"token": "123456",
		},
	}

	var verify = func(request []byte) {
		var body customerio.SendInAppRequest
		if err := json.Unmarshal(request, &body); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(&body, req) {
			t.Errorf("Request differed, want: %#v, got: %#v", request, body)
		}
	}

	api, srv := transactionalServer(t, verify)
	defer srv.Close()

	resp, err := api.SendInApp(context.Background(), req)
	if err != nil {
		t.Error(err)
	}

	expect := &customerio.SendInAppResponse{
		TransactionalResponse: customerio.TransactionalResponse{
			DeliveryID: testDeliveryID,
			QueuedAt:   time.Unix(int64(testQueuedAt), 0),
		},
	}

	if !reflect.DeepEqual(resp, expect) {
		t.Errorf("Expect: %#v, Got: %#v", expect, resp)
	}
}

func TestSendInAppError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(502)
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	resp, err := api.SendInApp(context.Background(), &customerio.SendInAppRequest{
		TransactionalMessageID: "123456",
		Identifiers: map[string]string{
			"id": "customer_1",
		},
	})
	if err == nil {
		t.Errorf("Expected error, got: %#v", resp)
	}

	if _, ok := err.(*customerio.TransactionalError); !ok {
		t.Errorf("Expected TransactionalError, got: %#v", err)
	}
}
