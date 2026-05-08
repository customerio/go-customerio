package customerio_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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
		defer func() { _ = req.Body.Close() }()
		verify(req.Method, req.URL.Path, b)
		fmt.Fprintf(w, `{"id":%d}`, expectedBroadcastResponseID)
	}))

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	return api, srv
}

func TestTriggerBroadcastSegment(t *testing.T) {
	broadcastID := 123
	data := map[string]any{"name": "Joe"}
	recipients := customerio.BroadcastRecipients{
		Segment: map[string]any{"id": float64(1)},
	}

	api, srv := broadcastServer(t, func(method, path string, body []byte) {
		if method != "POST" {
			t.Errorf("expected POST, got %s", method)
		}
		if path != "/v1/campaigns/123/triggers" {
			t.Errorf("expected /v1/campaigns/123/triggers, got %s", path)
		}

		var payload map[string]any
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

		wantRecipients := map[string]any{
			"segment": map[string]any{"id": float64(1)},
		}
		if !reflect.DeepEqual(payload["recipients"], wantRecipients) {
			t.Errorf("recipients mismatch: want %#v got %#v", wantRecipients, payload["recipients"])
		}
	})
	defer srv.Close()

	resp, err := api.TriggerBroadcast(context.Background(), broadcastID, data, recipients, customerio.BroadcastOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != expectedBroadcastResponseID {
		t.Errorf("unexpected response ID: %d", resp.ID)
	}
}

func TestTriggerBroadcastIDs(t *testing.T) {
	data := map[string]any{"promo": "SAVE10"}
	recipients := customerio.BroadcastRecipients{
		Ids: []string{"c1", "c2"},
	}
	opts := customerio.BroadcastOptions{
		IDIgnoreMissing: boolPtr(true),
		// These should NOT appear in the output for the ids path.
		EmailIgnoreMissing: boolPtr(true),
		EmailAddDuplicates: boolPtr(false),
	}

	api, srv := broadcastServer(t, func(_, _ string, body []byte) {
		var payload map[string]any
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
		wantIDs := []any{"c1", "c2"}
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

	if _, err := api.TriggerBroadcast(context.Background(), 42, data, recipients, opts); err != nil {
		t.Fatal(err)
	}
}

func TestTriggerBroadcastEmails(t *testing.T) {
	data := map[string]any{}
	recipients := customerio.BroadcastRecipients{
		Emails: []string{"a@example.com", "b@example.com"},
	}
	opts := customerio.BroadcastOptions{
		EmailIgnoreMissing: boolPtr(true),
		EmailAddDuplicates: boolPtr(false),
	}

	api, srv := broadcastServer(t, func(_, _ string, body []byte) {
		var payload map[string]any
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
		wantEmails := []any{"a@example.com", "b@example.com"}
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

	if _, err := api.TriggerBroadcast(context.Background(), 7, data, recipients, opts); err != nil {
		t.Fatal(err)
	}
}

func TestTriggerBroadcastPerUserData(t *testing.T) {
	data := map[string]any{"campaign": "spring"}
	recipients := customerio.BroadcastRecipients{
		PerUserData: []map[string]any{
			{"id": "u1", "data": map[string]any{"first_name": "Alice"}},
		},
	}
	opts := customerio.BroadcastOptions{
		IDIgnoreMissing:    boolPtr(false),
		EmailIgnoreMissing: boolPtr(true),
		EmailAddDuplicates: boolPtr(true),
	}

	api, srv := broadcastServer(t, func(_, _ string, body []byte) {
		var payload map[string]any
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

	if _, err := api.TriggerBroadcast(context.Background(), 99, data, recipients, opts); err != nil {
		t.Fatal(err)
	}
}

func TestTriggerBroadcastDataFileURL(t *testing.T) {
	data := map[string]any{}
	recipients := customerio.BroadcastRecipients{
		DataFileURL: "s3://mybucket/users.csv",
	}
	opts := customerio.BroadcastOptions{
		IDIgnoreMissing:    boolPtr(true),
		EmailIgnoreMissing: boolPtr(false),
	}

	api, srv := broadcastServer(t, func(_, path string, body []byte) {
		if path != "/v1/campaigns/5/triggers" {
			t.Errorf("unexpected path: %s", path)
		}
		var payload map[string]any
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

	if _, err := api.TriggerBroadcast(context.Background(), 5, data, recipients, opts); err != nil {
		t.Fatal(err)
	}
}

func TestTriggerBroadcastIDZero(t *testing.T) {
	api := customerio.NewAPIClient("myKey")

	_, err := api.TriggerBroadcast(context.Background(), 0, nil, customerio.BroadcastRecipients{}, customerio.BroadcastOptions{})
	checkParamError(t, err, "broadcastID")
}

func TestTriggerBroadcastIDNegative(t *testing.T) {
	api := customerio.NewAPIClient("myKey")

	_, err := api.TriggerBroadcast(context.Background(), -1, nil, customerio.BroadcastRecipients{}, customerio.BroadcastOptions{})
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

	_, err := api.TriggerBroadcast(context.Background(), 1, nil, customerio.BroadcastRecipients{}, customerio.BroadcastOptions{})
	if err == nil {
		t.Fatal("expected error")
	}
	cioErr, ok := err.(*customerio.CustomerIOError)
	if !ok {
		t.Fatalf("expected *CustomerIOError, got %T", err)
	}
	if cioErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, cioErr.StatusCode)
	}
	if !strings.Contains(string(cioErr.Body), "broadcast with id 1 does not exist") {
		t.Errorf("expected detail in Body, got %q", string(cioErr.Body))
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

	_, err := api.TriggerBroadcast(context.Background(), 1, nil, customerio.BroadcastRecipients{}, customerio.BroadcastOptions{})
	cioErr, ok := err.(*customerio.CustomerIOError)
	if !ok {
		t.Fatalf("expected *CustomerIOError, got %T", err)
	}
	if cioErr.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status %d, got %d", http.StatusBadGateway, cioErr.StatusCode)
	}
	if string(cioErr.Body) != "upstream issue" {
		t.Errorf("expected raw body in Body, got %q", string(cioErr.Body))
	}
}

func TestTriggerBroadcastCompanionOptionsIsolation(t *testing.T) {
	// When ids is the direct field, email_* options must be excluded even if set.
	recipients := customerio.BroadcastRecipients{
		Ids: []string{"u1"},
	}
	opts := customerio.BroadcastOptions{
		IDIgnoreMissing:    boolPtr(true),
		EmailIgnoreMissing: boolPtr(true),
		EmailAddDuplicates: boolPtr(true),
	}

	api, srv := broadcastServer(t, func(_, _ string, body []byte) {
		var payload map[string]any
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

	if _, err := api.TriggerBroadcast(context.Background(), 1, nil, recipients, opts); err != nil {
		t.Fatal(err)
	}
}
