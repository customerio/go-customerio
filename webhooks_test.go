package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestListReportingWebhooks(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"reporting_webhooks": []map[string]any{{"id": 1, "endpoint": "https://example.com/hook"}}})
	got, err := api.ListReportingWebhooks(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/reporting_webhooks")
	if len(got.ReportingWebhooks) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateReportingWebhook(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateReportingWebhook(context.Background(), customerio.ReportingWebhookInput{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "endpoint")
	}

	api, rec := apiServer(t, 200, map[string]any{"reporting_webhook": map[string]any{"id": 1}})
	input := customerio.ReportingWebhookInput{Endpoint: "https://example.com/hook", Events: []string{"email_sent"}}
	got, err := api.CreateReportingWebhook(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/reporting_webhooks")
	assertJSONEqual(t, rec.body, input)
	if got.ReportingWebhook.ID != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetUpdateDeleteReportingWebhook(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetReportingWebhook(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "webhookID")
	}

	api, rec := apiServer(t, 200, map[string]any{"reporting_webhook": map[string]any{"id": 1}})
	if _, err := api.GetReportingWebhook(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/reporting_webhooks/1")

	data := map[string]any{"disabled": true}
	if _, err := api.UpdateReportingWebhook(context.Background(), "1", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/reporting_webhooks/1")
	assertJSONEqual(t, rec.body, data)

	if err := api.DeleteReportingWebhook(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/reporting_webhooks/1")
}
