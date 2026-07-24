package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestListTransactionalMessages(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"messages": []map[string]any{
		{"id": 1, "name": "Welcome"},
	}})
	got, err := api.ListTransactionalMessages(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/transactional")
	if len(got.Messages) != 1 || got.Messages[0].Name != "Welcome" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetTransactionalMessage(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetTransactionalMessage(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "transactionalID")
	}

	api, rec := apiServer(t, 200, map[string]any{"message": map[string]any{"id": 1, "name": "Welcome"}})
	got, err := api.GetTransactionalMessage(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/transactional/1")
	if got.Message.Name != "Welcome" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetTransactionalMessageContents(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetTransactionalMessageContents(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "transactionalID")
	}

	api, rec := apiServer(t, 200, map[string]any{"contents": []map[string]any{
		{"id": 1, "language": "en", "subject": "hi"},
	}})
	got, err := api.GetTransactionalMessageContents(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/transactional/1/contents")
	if len(got.Contents) != 1 || got.Contents[0].Subject != "hi" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetTransactionalMessageLanguage(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetTransactionalMessageLanguage(context.Background(), "", "en"); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "transactionalID")
	}
	if _, err := api.GetTransactionalMessageLanguage(context.Background(), "1", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "language")
	}

	api, rec := apiServer(t, 200, map[string]any{"content": []map[string]any{{"language": "en", "subject": "hi"}}})
	got, err := api.GetTransactionalMessageLanguage(context.Background(), "1", "en")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/transactional/1/language/en")
	if len(got.Content) != 1 || got.Content[0].Language != "en" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestUpdateTransactionalMessageLanguage(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.UpdateTransactionalMessageLanguage(context.Background(), "", "en", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "transactionalID")
	}
	if _, err := api.UpdateTransactionalMessageLanguage(context.Background(), "1", "", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "language")
	}

	api, rec := apiServer(t, 200, map[string]any{"content": []map[string]any{{"language": "en"}}})
	data := map[string]any{"subject": "updated"}
	if _, err := api.UpdateTransactionalMessageLanguage(context.Background(), "1", "en", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/transactional/1/language/en")
	assertJSONEqual(t, rec.body, data)
}

func TestGetTransactionalMessageDeliveries(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetTransactionalMessageDeliveries(context.Background(), "", customerio.TransactionalDeliveriesOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "transactionalID")
	}

	tracked := true
	api, rec := apiServer(t, 200, map[string]any{"messages": []map[string]any{{"id": "m1", "subject": "hi"}}, "next": "cursor2"})
	got, err := api.GetTransactionalMessageDeliveries(context.Background(), "1", customerio.TransactionalDeliveriesOptions{
		Metric:              "opened",
		StartTS:             100,
		EndTS:               200,
		GetTrackedResponses: &tracked,
		PaginationOptions:   customerio.PaginationOptions{Limit: 10},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/transactional/1/messages?end_ts=200&get_tracked_responses=true&limit=10&metric=opened&start_ts=100")
	if len(got.Messages) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
	if got.Next != "cursor2" {
		t.Errorf("unexpected next cursor: %#v", got)
	}
}

func TestGetTransactionalMessageMetrics(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetTransactionalMessageMetrics(context.Background(), "", customerio.TransactionalMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "transactionalID")
	}

	api, rec := apiServer(t, 200, map[string]any{"metric": map[string]any{"series": map[string]any{"sent": []int{1, 2, 3}}}})
	got, err := api.GetTransactionalMessageMetrics(context.Background(), "1", customerio.TransactionalMetricsOptions{
		Period: customerio.MetricsPeriodDays,
		Steps:  7,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/transactional/1/metrics?period=days&steps=7")
	if len(got.Metric.Series.Sent) != 3 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetTransactionalMessageLinkMetrics(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetTransactionalMessageLinkMetrics(context.Background(), "", customerio.TransactionalLinkMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "transactionalID")
	}

	unique := true
	api, rec := apiServer(t, 200, map[string]any{"links": []map[string]any{
		{"link": map[string]any{"id": "l1", "href": "https://example.com"}, "metric": map[string]any{"series": map[string]any{"clicked": []int{1}}}},
	}})
	got, err := api.GetTransactionalMessageLinkMetrics(context.Background(), "1", customerio.TransactionalLinkMetricsOptions{
		TransactionalMetricsOptions: customerio.TransactionalMetricsOptions{Period: customerio.MetricsPeriodWeeks, Steps: 4},
		Unique:                      &unique,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/transactional/1/metrics/links?period=weeks&steps=4&unique=true")
	if len(got.Links) != 1 || got.Links[0].Link.Href != "https://example.com" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestUpdateTransactionalMessageContent(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.UpdateTransactionalMessageContent(context.Background(), "", "c1", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "transactionalID")
	}
	if _, err := api.UpdateTransactionalMessageContent(context.Background(), "1", "", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "contentID")
	}

	api, rec := apiServer(t, 200, map[string]any{"content": []map[string]any{{"id": 5}}})
	data := map[string]any{"body": "updated"}
	if _, err := api.UpdateTransactionalMessageContent(context.Background(), "1", "5", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/transactional/1/content/5")
	assertJSONEqual(t, rec.body, data)
}
