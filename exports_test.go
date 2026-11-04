package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestListExports(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"exports": []map[string]any{{"id": 1, "type": "customers"}}})
	got, err := api.ListExports(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/exports")
	if len(got.Exports) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetExport(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetExport(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "exportID")
	}

	api, rec := apiServer(t, 200, map[string]any{"export": map[string]any{"id": 1, "status": "done"}})
	got, err := api.GetExport(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/exports/1")
	if got.Export.Status != "done" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestDownloadExport(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.DownloadExport(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "exportID")
	}

	api, rec := apiServer(t, 200, map[string]any{"url": "https://example.com/export.csv"})
	got, err := api.DownloadExport(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/exports/1/download")
	if got.URL != "https://example.com/export.csv" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateCustomersExport(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"export": map[string]any{"id": 1, "type": "customers"}})
	filter := customerio.FilterByAttribute("plan", customerio.FilterOperatorEq, "gold")
	got, err := api.CreateCustomersExport(context.Background(), filter)
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/exports/customers")
	assertJSONEqual(t, rec.body, map[string]any{"filters": filter})
	if got.Export.ID != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateDeliveriesExport(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateDeliveriesExport(context.Background(), 0, customerio.DeliveryExportOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 200, map[string]any{"export": map[string]any{"id": 1, "type": "deliveries"}})
	got, err := api.CreateDeliveriesExport(context.Background(), 42, customerio.DeliveryExportOptions{
		Metric: customerio.DeliveryExportMetricOpened,
		Start:  100,
		End:    200,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/exports/deliveries")
	assertJSONEqual(t, rec.body, map[string]any{"newsletter_id": 42, "metric": "opened", "start": 100, "end": 200})
	if got.Export.ID != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}
