package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestGetMessages(t *testing.T) {
	drafts := false
	api, rec := apiServer(t, 200, map[string]any{"messages": []map[string]any{
		{"id": "m1", "type": "email", "campaign_id": 4},
	}})
	got, err := api.GetMessages(context.Background(), customerio.MessagesOptions{
		Drafts:            &drafts,
		Type:              "email",
		CampaignID:        "4",
		PaginationOptions: customerio.PaginationOptions{Limit: 50},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/messages?campaign_id=4&drafts=false&limit=50&type=email")
	if len(got.Messages) != 1 || got.Messages[0].CampaignID != 4 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetMessagesNoOptions(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"messages": []any{}})
	if _, err := api.GetMessages(context.Background(), customerio.MessagesOptions{}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/messages")
}

func TestGetMessage(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetMessage(context.Background(), "", customerio.MessageOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "messageID")
	}

	associations := true
	api, rec := apiServer(t, 200, map[string]any{"message": map[string]any{"id": "m1", "subject": "hi"}})
	got, err := api.GetMessage(context.Background(), "m1", customerio.MessageOptions{Associations: &associations})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/messages/m1?associations=true")
	if got.Message.Subject != "hi" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetArchivedMessage(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetArchivedMessage(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "messageID")
	}

	api, rec := apiServer(t, 200, map[string]any{"archived_message": map[string]any{"id": "m1", "body": "<html></html>"}})
	got, err := api.GetArchivedMessage(context.Background(), "m1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/messages/m1/archived_message")
	if got.ArchivedMessage.Body != "<html></html>" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetArchivedMessageRateLimited(t *testing.T) {
	api, _ := apiServer(t, 429, map[string]any{"errors": []map[string]any{{"detail": "rate limited"}}})
	_, err := api.GetArchivedMessage(context.Background(), "m1")
	cioErr, ok := err.(*customerio.CustomerIOError)
	if !ok {
		t.Fatalf("expected *CustomerIOError, got %T", err)
	}
	if cioErr.StatusCode() != 429 {
		t.Errorf("expected status 429, got %d", cioErr.StatusCode())
	}
}
