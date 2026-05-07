package customerio_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

const expectedBroadcastResponseID = 999

func boolPtr(b bool) *bool { return &b }

func broadcastServer(t *testing.T, verify func(method, path string, body []byte)) (*customerio.APIClient, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
		}
		defer req.Body.Close()
		verify(req.Method, req.URL.Path, b)
		fmt.Fprintf(w, `{"id":%d}`, expectedBroadcastResponseID)
	}))

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	return api, srv
}

func TestTriggerBroadcastSegment(t *testing.T) {
	broadcastID := 123
	data := map[string]interface{}{"name": "Joe"}
	recipients := customerio.BroadcastRecipients{
		Segment: map[string]interface{}{"id": float64(1)},
	}

	api, srv := broadcastServer(t, func(method, path string, body []byte) {
		if method != "POST" {
			t.Errorf("expected POST, got %s", method)
		}
		if path != "/v1/campaigns/123/triggers" {
			t.Errorf("expected /v1/campaigns/123/triggers, got %s", path)
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}

		if _, ok := payload["recipients"]; !ok {
			t.Error("expected top-level recipients key for segment-based broadcast")
		}
		if _, ok := payload["data"]; !ok {
			t.Error("expected top-level data key")
		}
		if _, ok := payload["ids"]; ok {
			t.Error("ids should not be present in segment-based payload")
		}

		wantRecipients := map[string]interface{}{
			"segment": map[string]interface{}{"id": float64(1)},
		}
		if !reflect.DeepEqual(payload["recipients"], wantRecipients) {
			t.Errorf("recipients mismatch: want %#v got %#v", wantRecipients, payload["recipients"])
		}
	})
	defer srv.Close()

	resp, err := api.TriggerBroadcast(context.Background(), broadcastID, data, recipients)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != expectedBroadcastResponseID {
		t.Errorf("unexpected response ID: %d", resp.ID)
	}
}

func TestTriggerBroadcastIDs(t *testing.T) {
	data := map[string]interface{}{"promo": "SAVE10"}
	recipients := customerio.BroadcastRecipients{
		Ids:             []string{"c1", "c2"},
		IDIgnoreMissing: boolPtr(true),
		// These should NOT appear in the output for the ids path.
		EmailIgnoreMissing: boolPtr(true),
		EmailAddDuplicates: boolPtr(false),
	}

	api, srv := broadcastServer(t, func(_, _ string, body []byte) {
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}

		if _, ok := payload["recipients"]; ok {
			t.Error("recipients key must not appear in direct ids payload")
		}
		ids, ok := payload["ids"]
		if !ok {
			t.Fatal("expected ids key")
		}
		wantIDs := []interface{}{"c1", "c2"}
		if !reflect.DeepEqual(ids, wantIDs) {
			t.Errorf("ids mismatch: want %v got %v", wantIDs, ids)
		}
		if payload["id_ignore_missing"] != true {
			t.Errorf("expected id_ignore_missing=true, got %v", payload["id_ignore_missing"])
		}
		if _, ok := payload["email_ignore_missing"]; ok {
			t.Error("email_ignore_missing must not appear in ids payload")
		}
		if _, ok := payload["email_add_duplicates"]; ok {
			t.Error("email_add_duplicates must not appear in ids payload")
		}
	})
	defer srv.Close()

	if _, err := api.TriggerBroadcast(context.Background(), 42, data, recipients); err != nil {
		t.Fatal(err)
	}
}

func TestTriggerBroadcastEmails(t *testing.T) {
	data := map[string]interface{}{}
	recipients := customerio.BroadcastRecipients{
		Emails:             []string{"a@example.com", "b@example.com"},
		EmailIgnoreMissing: boolPtr(true),
		EmailAddDuplicates: boolPtr(false),
	}

	api, srv := broadcastServer(t, func(_, _ string, body []byte) {
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}

		if _, ok := payload["recipients"]; ok {
			t.Error("recipients key must not appear in direct emails payload")
		}
		emails, ok := payload["emails"]
		if !ok {
			t.Fatal("expected emails key")
		}
		wantEmails := []interface{}{"a@example.com", "b@example.com"}
		if !reflect.DeepEqual(emails, wantEmails) {
			t.Errorf("emails mismatch: want %v got %v", wantEmails, emails)
		}
		if payload["email_ignore_missing"] != true {
			t.Errorf("expected email_ignore_missing=true, got %v", payload["email_ignore_missing"])
		}
		if payload["email_add_duplicates"] != false {
			t.Errorf("expected email_add_duplicates=false, got %v", payload["email_add_duplicates"])
		}
		if _, ok := payload["id_ignore_missing"]; ok {
			t.Error("id_ignore_missing must not appear in emails payload")
		}
	})
	defer srv.Close()

	if _, err := api.TriggerBroadcast(context.Background(), 7, data, recipients); err != nil {
		t.Fatal(err)
	}
}

func TestTriggerBroadcastPerUserData(t *testing.T) {
	data := map[string]interface{}{"campaign": "spring"}
	recipients := customerio.BroadcastRecipients{
		PerUserData: []map[string]interface{}{
			{"id": "u1", "data": map[string]interface{}{"first_name": "Alice"}},
		},
		IDIgnoreMissing:    boolPtr(false),
		EmailIgnoreMissing: boolPtr(true),
		EmailAddDuplicates: boolPtr(true),
	}

	api, srv := broadcastServer(t, func(_, _ string, body []byte) {
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}

		if _, ok := payload["per_user_data"]; !ok {
			t.Fatal("expected per_user_data key")
		}
		if _, ok := payload["recipients"]; ok {
			t.Error("recipients key must not appear in per_user_data payload")
		}
		if payload["id_ignore_missing"] != false {
			t.Errorf("expected id_ignore_missing=false, got %v", payload["id_ignore_missing"])
		}
		if payload["email_ignore_missing"] != true {
			t.Errorf("expected email_ignore_missing=true, got %v", payload["email_ignore_missing"])
		}
		if payload["email_add_duplicates"] != true {
			t.Errorf("expected email_add_duplicates=true, got %v", payload["email_add_duplicates"])
		}
	})
	defer srv.Close()

	if _, err := api.TriggerBroadcast(context.Background(), 99, data, recipients); err != nil {
		t.Fatal(err)
	}
}

func TestTriggerBroadcastDataFileURL(t *testing.T) {
	data := map[string]interface{}{}
	recipients := customerio.BroadcastRecipients{
		DataFileURL:        "s3://mybucket/users.csv",
		IDIgnoreMissing:    boolPtr(true),
		EmailIgnoreMissing: boolPtr(false),
	}

	api, srv := broadcastServer(t, func(_, path string, body []byte) {
		if path != "/v1/campaigns/5/triggers" {
			t.Errorf("unexpected path: %s", path)
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}

		if payload["data_file_url"] != "s3://mybucket/users.csv" {
			t.Errorf("unexpected data_file_url: %v", payload["data_file_url"])
		}
		if _, ok := payload["recipients"]; ok {
			t.Error("recipients key must not appear in data_file_url payload")
		}
		if payload["id_ignore_missing"] != true {
			t.Errorf("expected id_ignore_missing=true, got %v", payload["id_ignore_missing"])
		}
		if payload["email_ignore_missing"] != false {
			t.Errorf("expected email_ignore_missing=false, got %v", payload["email_ignore_missing"])
		}
		if _, ok := payload["email_add_duplicates"]; ok {
			t.Error("email_add_duplicates should be absent when not set")
		}
	})
	defer srv.Close()

	if _, err := api.TriggerBroadcast(context.Background(), 5, data, recipients); err != nil {
		t.Fatal(err)
	}
}

func TestTriggerBroadcastIDZero(t *testing.T) {
	api := customerio.NewAPIClient("myKey")

	_, err := api.TriggerBroadcast(context.Background(), 0, nil, customerio.BroadcastRecipients{})
	checkParamError(t, err, "broadcastID")
}

func TestTriggerBroadcastIDNegative(t *testing.T) {
	api := customerio.NewAPIClient("myKey")

	_, err := api.TriggerBroadcast(context.Background(), -1, nil, customerio.BroadcastRecipients{})
	checkParamError(t, err, "broadcastID")
}

func TestTriggerBroadcastError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"detail":"broadcast with id 1 does not exist","status":"404"}]}`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.TriggerBroadcast(context.Background(), 1, nil, customerio.BroadcastRecipients{})
	if err == nil {
		t.Fatal("expected error")
	}
	te, ok := err.(*customerio.TransactionalError)
	if !ok {
		t.Fatalf("expected *TransactionalError, got %T", err)
	}
	if te.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, te.StatusCode)
	}
	if te.Err != "broadcast with id 1 does not exist" {
		t.Errorf("expected detail in Err, got %q", te.Err)
	}
}

func TestTriggerBroadcastErrorUnparseableBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`upstream issue`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.TriggerBroadcast(context.Background(), 1, nil, customerio.BroadcastRecipients{})
	te, ok := err.(*customerio.TransactionalError)
	if !ok {
		t.Fatalf("expected *TransactionalError, got %T", err)
	}
	if te.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status %d, got %d", http.StatusBadGateway, te.StatusCode)
	}
	if te.Err != "upstream issue" {
		t.Errorf("expected raw body in Err, got %q", te.Err)
	}
}

func TestTriggerBroadcastCompanionOptionsIsolation(t *testing.T) {
	// When ids is the direct field, email_* options must be excluded even if set.
	recipients := customerio.BroadcastRecipients{
		Ids:                []string{"u1"},
		IDIgnoreMissing:    boolPtr(true),
		EmailIgnoreMissing: boolPtr(true),
		EmailAddDuplicates: boolPtr(true),
	}

	api, srv := broadcastServer(t, func(_, _ string, body []byte) {
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}
		if _, ok := payload["email_ignore_missing"]; ok {
			t.Error("email_ignore_missing must not be included in ids payload")
		}
		if _, ok := payload["email_add_duplicates"]; ok {
			t.Error("email_add_duplicates must not be included in ids payload")
		}
		if payload["id_ignore_missing"] != true {
			t.Errorf("id_ignore_missing should be true, got %v", payload["id_ignore_missing"])
		}
	})
	defer srv.Close()

	if _, err := api.TriggerBroadcast(context.Background(), 1, nil, recipients); err != nil {
		t.Fatal(err)
	}
}
