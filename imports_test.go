package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestCreateImport(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateImport(context.Background(), customerio.ImportInput{Type: customerio.ImportTypePeople}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "dataFileURL")
	}
	if _, err := api.CreateImport(context.Background(), customerio.ImportInput{DataFileURL: "https://example.com/f.csv"}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "type")
	}

	api, rec := apiServer(t, 200, map[string]any{"import": map[string]any{"id": 1, "state": "preprocessing"}})
	input := customerio.ImportInput{DataFileURL: "https://example.com/f.csv", Type: customerio.ImportTypePeople}
	got, err := api.CreateImport(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/imports")
	assertJSONEqual(t, rec.body, map[string]any{"import": input})
	if got.Import.ID != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetImport(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetImport(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "importID")
	}

	api, rec := apiServer(t, 200, map[string]any{"import": map[string]any{"id": 1, "state": "imported"}})
	got, err := api.GetImport(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/imports/1")
	if got.Import.State != "imported" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestBatchUpdateAttributes(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if err := api.BatchUpdateAttributes(context.Background(), nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "attributes")
	}

	api, rec := apiServer(t, 204, nil)
	attrs := []customerio.DataIndexAttribute{{Name: "plan"}}
	if err := api.BatchUpdateAttributes(context.Background(), attrs); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/data_index/attributes")
	assertJSONEqual(t, rec.body, map[string]any{"attributes": attrs})
}

func TestBatchUpdateEvents(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if err := api.BatchUpdateEvents(context.Background(), nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "events")
	}

	api, rec := apiServer(t, 204, nil)
	events := []customerio.DataIndexEvent{{Name: "purchased"}}
	if err := api.BatchUpdateEvents(context.Background(), events); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/data_index/events")
	assertJSONEqual(t, rec.body, map[string]any{"events": events})
}
