package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestGetSenderIdentities(t *testing.T) {
	hidden := false
	api, rec := apiServer(t, 200, map[string]any{"sender_identities": []map[string]any{{"id": 1, "email": "team@example.com"}}})
	got, err := api.GetSenderIdentities(context.Background(), customerio.SenderIdentitiesOptions{
		Sort:   customerio.SortAscending,
		Hidden: &hidden,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/sender_identities?hidden=false&sort=asc")
	if len(got.SenderIdentities) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetSenderIdentity(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetSenderIdentity(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "senderID")
	}

	api, rec := apiServer(t, 200, map[string]any{"sender_identity": map[string]any{"id": 1, "email": "team@example.com"}})
	got, err := api.GetSenderIdentity(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/sender_identities/1")
	if got.SenderIdentity.Email != "team@example.com" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetSenderIdentityUsedBy(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetSenderIdentityUsedBy(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "senderID")
	}

	api, rec := apiServer(t, 200, map[string]any{"sender_id": 1, "used_by": map[string]any{"campaigns": []int{1, 2}}})
	got, err := api.GetSenderIdentityUsedBy(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/sender_identities/1/used_by")
	if got.SenderID != 1 || got.UsedBy["campaigns"] == nil {
		t.Errorf("unexpected response: %#v", got)
	}
}
