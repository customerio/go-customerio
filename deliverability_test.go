package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestSearchSuppression(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.SearchSuppression(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "email")
	}

	api, rec := apiServer(t, 200, map[string]any{"category": "bounces", "suppressions": []map[string]any{
		{"email": "a@example.com", "reason": "hard bounce"},
	}})
	got, err := api.SearchSuppression(context.Background(), "a@example.com")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/esp/search_suppression/a@example.com")
	if got.Category != "bounces" || len(got.Suppressions) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetSuppressions(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetSuppressions(context.Background(), "", customerio.SuppressionsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "suppressionType")
	}

	api, rec := apiServer(t, 200, map[string]any{"suppressions": []map[string]any{{"email": "a@example.com"}}})
	got, err := api.GetSuppressions(context.Background(), customerio.SuppressionTypeBounces, customerio.SuppressionsOptions{
		Limit: 50, Offset: 10, Domain: "example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/esp/suppression/bounces?domain=example.com&limit=50&offset=10")
	if len(got.Suppressions) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetDomainSuppressions(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetDomainSuppressions(context.Background(), "", "bounces", customerio.DomainSuppressionsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "domainName")
	}
	if _, err := api.GetDomainSuppressions(context.Background(), "example.com", "", customerio.DomainSuppressionsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "suppressionType")
	}

	api, rec := apiServer(t, 200, map[string]any{"suppressions": []any{}, "next": "cursor2"})
	got, err := api.GetDomainSuppressions(context.Background(), "example.com", customerio.SuppressionTypeSpamReports, customerio.DomainSuppressionsOptions{Start: "cursor1"})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/esp/domains/example.com/suppression/spam_reports?start=cursor1")
	if got.Next != "cursor2" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateDeleteSuppression(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if err := api.CreateSuppression(context.Background(), "", "a@example.com"); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "suppressionType")
	}
	if err := api.CreateSuppression(context.Background(), customerio.SuppressionTypeBlocks, ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "email")
	}

	api, rec := apiServer(t, 204, nil)
	if err := api.CreateSuppression(context.Background(), customerio.SuppressionTypeInvalidEmails, "a@example.com"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/esp/suppression/invalid_emails/a@example.com")

	if err := api.DeleteSuppression(context.Background(), customerio.SuppressionTypeInvalidEmails, "a@example.com"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/esp/suppression/invalid_emails/a@example.com")
}
