package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestGetSnippets(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"snippets": []map[string]any{{"name": "greeting", "value": "Hi {{name}}"}}})
	got, err := api.GetSnippets(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/snippets")
	if len(got.Snippets) != 1 || got.Snippets[0].Name != "greeting" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateUpdateSnippet(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateSnippet(context.Background(), customerio.SnippetInput{Value: "x"}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "name")
	}
	if _, err := api.CreateSnippet(context.Background(), customerio.SnippetInput{Name: "x"}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "value")
	}

	api, rec := apiServer(t, 200, map[string]any{"snippet": map[string]any{"name": "greeting", "value": "Hi"}})
	input := customerio.SnippetInput{Name: "greeting", Value: "Hi"}
	if _, err := api.CreateSnippet(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/snippets")
	assertJSONEqual(t, rec.body, input)

	if _, err := api.UpdateSnippet(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/snippets")

	if _, err := api.UpdateSnippet(context.Background(), customerio.SnippetInput{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "name")
	}
}

func TestDeleteSnippet(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if err := api.DeleteSnippet(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "name")
	}

	api, rec := apiServer(t, 204, nil)
	if err := api.DeleteSnippet(context.Background(), "greeting"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/snippets/greeting")
}
