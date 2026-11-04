package customerio_test

import (
	"context"
	"testing"
)

func TestListSubscriptionTopics(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"topics": []map[string]any{
		{"id": 4, "identifier": "topic_4", "name": "Product Updates", "subscribed_by_default": false},
	}})
	got, err := api.ListSubscriptionTopics(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/subscription_topics")
	if len(got.Topics) != 1 || got.Topics[0].Name != "Product Updates" {
		t.Errorf("unexpected response: %#v", got)
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
