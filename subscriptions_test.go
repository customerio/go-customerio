package customerio_test

import (
	"context"
	"testing"
)

func TestListSubscriptionTopics(t *testing.T) {
	// The list endpoint may return topic IDs as an integer or a string; both
	// must decode via the lenient json.Number type on SubscriptionTopic.ID.
	api, rec := apiServer(t, 200, map[string]any{"topics": []map[string]any{
		{"id": 4, "identifier": "topic_4", "name": "Product Updates", "subscribed_by_default": false},
		{"id": "7", "identifier": "topic_7", "name": "Newsletter", "subscribed_by_default": true},
	}})
	got, err := api.ListSubscriptionTopics(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/subscription_topics")
	if len(got.Topics) != 2 || got.Topics[0].Name != "Product Updates" {
		t.Errorf("unexpected response: %#v", got)
	}
	if got.Topics[0].ID.String() != "4" || got.Topics[1].ID.String() != "7" {
		t.Errorf("unexpected topic IDs: %q, %q", got.Topics[0].ID, got.Topics[1].ID)
	}
}

func TestListSubscriptionChannels(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"channels": []map[string]any{
		{"id": 1, "type": "email", "name": "Email", "subscribed_by_default": true},
	}})
	got, err := api.ListSubscriptionChannels(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/subscription_channels")
	if len(got.Channels) != 1 || got.Channels[0].Type != "email" {
		t.Errorf("unexpected response: %#v", got)
	}
}
