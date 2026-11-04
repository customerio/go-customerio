package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestListCampaigns(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"campaigns": []map[string]any{{"id": 1, "name": "Onboarding"}}})
	got, err := api.ListCampaigns(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns")
	if len(got.Campaigns) != 1 || got.Campaigns[0].Name != "Onboarding" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetCampaign(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCampaign(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "campaignID")
	}

	api, rec := apiServer(t, 200, map[string]any{"campaign": map[string]any{"id": 1, "state": "running"}})
	got, err := api.GetCampaign(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns/1")
	if got.Campaign.State != "running" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetCampaignActions(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCampaignActions(context.Background(), "", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "campaignID")
	}

	api, rec := apiServer(t, 200, map[string]any{"actions": []map[string]any{{"id": 5, "type": "email"}}})
	got, err := api.GetCampaignActions(context.Background(), "1", "cursor1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns/1/actions?start=cursor1")
	if len(got.Actions) != 1 || got.Actions[0].Type != "email" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetCampaignAction(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCampaignAction(context.Background(), "", "5"); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "campaignID")
	}
	if _, err := api.GetCampaignAction(context.Background(), "1", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "actionID")
	}

	api, rec := apiServer(t, 200, map[string]any{"action": map[string]any{"id": 5}})
	if _, err := api.GetCampaignAction(context.Background(), "1", "5"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns/1/actions/5")
}

func TestUpdateCampaignAction(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.UpdateCampaignAction(context.Background(), "", "5", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "campaignID")
	}
	if _, err := api.UpdateCampaignAction(context.Background(), "1", "", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "actionID")
	}

	api, rec := apiServer(t, 200, map[string]any{"action": map[string]any{"id": 5}})
	data := map[string]any{"subject": "new subject"}
	if _, err := api.UpdateCampaignAction(context.Background(), "1", "5", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/campaigns/1/actions/5")
	assertJSONEqual(t, rec.body, data)
}

func TestCampaignActionLanguage(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCampaignActionLanguage(context.Background(), "", "5", "en"); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "campaignID")
	}
	if _, err := api.GetCampaignActionLanguage(context.Background(), "1", "", "en"); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "actionID")
	}
	if _, err := api.GetCampaignActionLanguage(context.Background(), "1", "5", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "language")
	}

	api, rec := apiServer(t, 200, map[string]any{"action": map[string]any{"id": 5, "language": "en"}})
	if _, err := api.GetCampaignActionLanguage(context.Background(), "1", "5", "en"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns/1/actions/5/language/en")

	data := map[string]any{"subject": "hi"}
	if _, err := api.UpdateCampaignActionLanguage(context.Background(), "1", "5", "en", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/campaigns/1/actions/5/language/en")
	assertJSONEqual(t, rec.body, data)
}

func TestGetCampaignActionMetrics(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCampaignActionMetrics(context.Background(), "", "5", customerio.CampaignActionMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "campaignID")
	}
	if _, err := api.GetCampaignActionMetrics(context.Background(), "1", "", customerio.CampaignActionMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "actionID")
	}

	api, rec := apiServer(t, 200, map[string]any{"metric": map[string]any{"series": map[string]any{"sent": []int{1}}}})
	got, err := api.GetCampaignActionMetrics(context.Background(), "1", "5", customerio.CampaignActionMetricsOptions{
		Version:                     customerio.CampaignMetricsVersion2,
		Res:                         "days",
		TZ:                          "UTC",
		Start:                       100,
		End:                         200,
		TransactionalMetricsOptions: customerio.TransactionalMetricsOptions{Period: customerio.MetricsPeriodDays, Steps: 3},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns/1/actions/5/metrics?end=200&period=days&res=days&start=100&steps=3&tz=UTC&version=2")
	if len(got.Metric.Series.Sent) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetCampaignActionMetricsLinks(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"links": []map[string]any{}})
	unique := true
	if _, err := api.GetCampaignActionMetricsLinks(context.Background(), "1", "5", customerio.TransactionalLinkMetricsOptions{
		TransactionalMetricsOptions: customerio.TransactionalMetricsOptions{Period: customerio.MetricsPeriodWeeks, Steps: 2},
		Unique:                      &unique,
	}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns/1/actions/5/metrics/links?period=weeks&steps=2&unique=true")
}

func TestGetCampaignMetrics(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCampaignMetrics(context.Background(), "", customerio.CampaignMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "campaignID")
	}

	api, rec := apiServer(t, 200, map[string]any{"metric": map[string]any{"series": map[string]any{}}})
	if _, err := api.GetCampaignMetrics(context.Background(), "1", customerio.CampaignMetricsOptions{Type: customerio.MetricTypeEmail}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns/1/metrics?type=email")
}

func TestGetCampaignMetricsLinks(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCampaignMetricsLinks(context.Background(), "", customerio.TransactionalLinkMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "campaignID")
	}

	api, rec := apiServer(t, 200, map[string]any{"links": []any{}})
	if _, err := api.GetCampaignMetricsLinks(context.Background(), "1", customerio.TransactionalLinkMetricsOptions{}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns/1/metrics/links")
}

func TestGetCampaignJourneyMetrics(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCampaignJourneyMetrics(context.Background(), "", customerio.JourneyMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "campaignID")
	}
	if _, err := api.GetCampaignJourneyMetrics(context.Background(), "1", customerio.JourneyMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "start")
	}
	if _, err := api.GetCampaignJourneyMetrics(context.Background(), "1", customerio.JourneyMetricsOptions{Start: 100}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "end")
	}
	if _, err := api.GetCampaignJourneyMetrics(context.Background(), "1", customerio.JourneyMetricsOptions{Start: 100, End: 200}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "resolution")
	}

	api, rec := apiServer(t, 200, map[string]any{"journey_metric": map[string]any{"started": []int{1, 2}}})
	got, err := api.GetCampaignJourneyMetrics(context.Background(), "1", customerio.JourneyMetricsOptions{Start: 100, End: 200, Resolution: "days"})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns/1/journey_metrics?end=200&resolution=days&start=100")
	if len(got.JourneyMetric.Started) != 2 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetCampaignMessages(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCampaignMessages(context.Background(), "", customerio.CampaignMessagesOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "campaignID")
	}

	drafts := true
	api, rec := apiServer(t, 200, map[string]any{"messages": []map[string]any{{"id": "m1"}}})
	got, err := api.GetCampaignMessages(context.Background(), "1", customerio.CampaignMessagesOptions{
		Type:              customerio.MetricTypeEmail,
		Drafts:            &drafts,
		PaginationOptions: customerio.PaginationOptions{Limit: 10},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/campaigns/1/messages?drafts=true&limit=10&type=email")
	if len(got.Messages) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}
