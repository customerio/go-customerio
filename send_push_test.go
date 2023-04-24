package customerio_test

import (
	"context"
	"encoding/json"
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
		MessageData: map[string]interface{}{
			"client": "Go",
			"name":   "gopher",
		},
	}
	d, err := customerio.NewDevice("device-id", "ios", map[string]interface{}{"attr1": "value1"})
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
