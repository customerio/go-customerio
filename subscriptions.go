package customerio

import (
	"context"
	"encoding/json"
)

// SubscriptionTopic is a topic customers can subscribe to (e.g. "Product Updates").
//
// ID is decoded leniently (json.Number): the App API returns subscription topic
// IDs as either a JSON string or an integer across endpoints, so this mirrors
// SubscriptionTopicPreference.ID rather than pinning a single numeric type.
type SubscriptionTopic struct {
	ID                  json.Number `json:"id"`
	Identifier          string      `json:"identifier,omitempty"`
	Name                string      `json:"name"`
	Description         string      `json:"description,omitempty"`
	SubscribedByDefault bool        `json:"subscribed_by_default"`
}

// ListSubscriptionTopicsResponse is the decoded shape of GET /v1/subscription_topics.
type ListSubscriptionTopicsResponse struct {
	Topics []SubscriptionTopic `json:"topics"`
}

// ListSubscriptionTopics returns every subscription topic defined in the workspace.
// See https://docs.customer.io/api/app/#operation/listSubscriptionTopics
func (c *APIClient) ListSubscriptionTopics(ctx context.Context) (*ListSubscriptionTopicsResponse, error) {
	var resp ListSubscriptionTopicsResponse
	if err := c.doJSON(ctx, "GET", "/v1/subscription_topics", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SubscriptionChannel is a delivery channel (email, push, SMS, ...)
// customers can opt in or out of.
type SubscriptionChannel struct {
	ID                  int    `json:"id"`
	Type                string `json:"type"`
	Name                string `json:"name"`
	Description         string `json:"description,omitempty"`
	SubscribedByDefault bool   `json:"subscribed_by_default"`
}

// ListSubscriptionChannelsResponse is the decoded shape of GET /v1/subscription_channels.
type ListSubscriptionChannelsResponse struct {
	Channels []SubscriptionChannel `json:"channels"`
}

// ListSubscriptionChannels returns every subscription channel defined in the workspace.
// See https://docs.customer.io/api/app/#operation/listSubscriptionChannels
func (c *APIClient) ListSubscriptionChannels(ctx context.Context) (*ListSubscriptionChannelsResponse, error) {
	var resp ListSubscriptionChannelsResponse
	if err := c.doJSON(ctx, "GET", "/v1/subscription_channels", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
