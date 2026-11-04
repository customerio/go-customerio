package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestListCollections(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"collections": []map[string]any{{"id": 1, "name": "Holidays"}}})
	got, err := api.ListCollections(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/collections")
	if len(got.Collections) != 1 || got.Collections[0].Name != "Holidays" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateCollection(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateCollection(context.Background(), customerio.CollectionInput{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "name")
	}

	api, rec := apiServer(t, 200, map[string]any{"collection": map[string]any{"id": 1, "name": "Holidays"}})
	input := customerio.CollectionInput{Name: "Holidays", Data: []map[string]any{{"date": "2026-12-25"}}}
	got, err := api.CreateCollection(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/collections")
	assertJSONEqual(t, rec.body, input)
	if got.Collection.ID != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetUpdateDeleteCollection(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCollection(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "collectionID")
	}

	api, rec := apiServer(t, 200, map[string]any{"collection": map[string]any{"id": 1}})
	if _, err := api.GetCollection(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/collections/1")

	update := customerio.CollectionUpdate{URL: "https://example.com/data.csv"}
	if _, err := api.UpdateCollection(context.Background(), "1", update); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/collections/1")
	assertJSONEqual(t, rec.body, update)

	if err := api.DeleteCollection(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/collections/1")
}

func TestGetCollectionContent(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCollectionContent(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "collectionID")
	}

	api, rec := apiServer(t, 200, map[string]any{"eventName": "christmas", "presents": 2})
	got, err := api.GetCollectionContent(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/collections/1/content")
	if got["eventName"] != "christmas" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestUpdateCollectionContent(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if err := api.UpdateCollectionContent(context.Background(), "", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "collectionID")
	}
	if err := api.UpdateCollectionContent(context.Background(), "1", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "rows")
	}

	api, rec := apiServer(t, 204, nil)
	rows := []map[string]any{{"date": "2026-12-25"}, {"date": "2026-01-01"}}
	if err := api.UpdateCollectionContent(context.Background(), "1", rows); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/collections/1/content")
	gotRows, ok := rec.body.([]any)
	if !ok || len(gotRows) != 2 {
		t.Errorf("expected 2 rows sent as a top-level array, got %#v", rec.body)
	}
}
