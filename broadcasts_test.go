package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestListBroadcasts(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"broadcasts": []map[string]any{{"id": 1, "name": "Flash Sale"}}})
	got, err := api.ListBroadcasts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts")
	if len(got.Broadcasts) != 1 || got.Broadcasts[0].Name != "Flash Sale" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetBroadcast(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetBroadcast(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "broadcastID")
	}

	api, rec := apiServer(t, 200, map[string]any{"broadcast": map[string]any{"id": 1}})
	if _, err := api.GetBroadcast(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts/1")
}

func TestGetBroadcastActions(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetBroadcastActions(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "broadcastID")
	}

	api, rec := apiServer(t, 200, map[string]any{"actions": []map[string]any{{"id": 5}}})
	got, err := api.GetBroadcastActions(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts/1/actions")
	if len(got.Actions) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestBroadcastAction(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetBroadcastAction(context.Background(), "", "5"); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "broadcastID")
	}
	if _, err := api.GetBroadcastAction(context.Background(), "1", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "actionID")
	}

	api, rec := apiServer(t, 200, map[string]any{"action": map[string]any{"id": 5}})
	if _, err := api.GetBroadcastAction(context.Background(), "1", "5"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts/1/actions/5")

	data := map[string]any{"subject": "hi"}
	if _, err := api.UpdateBroadcastAction(context.Background(), "1", "5", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/broadcasts/1/actions/5")
	assertJSONEqual(t, rec.body, data)
}

func TestBroadcastActionLanguage(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"action": map[string]any{"id": 5, "language": "en"}})
	if _, err := api.GetBroadcastActionLanguage(context.Background(), "1", "5", "en"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts/1/actions/5/language/en")

	data := map[string]any{"subject": "hi"}
	if _, err := api.UpdateBroadcastActionLanguage(context.Background(), "1", "5", "en", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/broadcasts/1/actions/5/language/en")
	assertJSONEqual(t, rec.body, data)

	if _, err := api.GetBroadcastActionLanguage(context.Background(), "1", "5", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "language")
	}
}

func TestGetBroadcastActionMetrics(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"metric": map[string]any{"series": map[string]any{}}})
	if _, err := api.GetBroadcastActionMetrics(context.Background(), "1", "5", customerio.TransactionalMetricsOptions{Period: customerio.MetricsPeriodHours, Steps: 6}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts/1/actions/5/metrics?period=hours&steps=6")

	if _, err := api.GetBroadcastActionMetrics(context.Background(), "", "5", customerio.TransactionalMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "broadcastID")
	}
}

func TestGetBroadcastActionMetricsLinks(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"links": []any{}})
	if _, err := api.GetBroadcastActionMetricsLinks(context.Background(), "1", "5", customerio.TransactionalLinkMetricsOptions{}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts/1/actions/5/metrics/links")
}

func TestGetBroadcastMetrics(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetBroadcastMetrics(context.Background(), "", customerio.BroadcastMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "broadcastID")
	}

	api, rec := apiServer(t, 200, map[string]any{"metric": map[string]any{"series": map[string]any{}}})
	if _, err := api.GetBroadcastMetrics(context.Background(), "1", customerio.BroadcastMetricsOptions{Type: customerio.MetricTypePush}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts/1/metrics?type=push")
}

func TestGetBroadcastMetricsLinks(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"links": []any{}})
	if _, err := api.GetBroadcastMetricsLinks(context.Background(), "1", customerio.TransactionalLinkMetricsOptions{}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts/1/metrics/links")
}

func TestGetBroadcastMessages(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetBroadcastMessages(context.Background(), "", customerio.BroadcastMessagesOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "broadcastID")
	}

	api, rec := apiServer(t, 200, map[string]any{"messages": []map[string]any{{"id": "m1"}}})
	got, err := api.GetBroadcastMessages(context.Background(), "1", customerio.BroadcastMessagesOptions{Metric: "opened"})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts/1/messages?metric=opened")
	if len(got.Messages) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetBroadcastTriggers(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetBroadcastTriggers(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "broadcastID")
	}

	api, rec := apiServer(t, 200, map[string]any{"triggers": []map[string]any{{"id": 9, "broadcast_id": 1}}})
	got, err := api.GetBroadcastTriggers(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/broadcasts/1/triggers")
	if len(got.Triggers) != 1 || got.Triggers[0].ID != 9 {
		t.Errorf("unexpected response: %#v", got)
	}
}
