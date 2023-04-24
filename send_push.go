package customerio

import (
	"context"
	"encoding/json"
)

type SendPushRequest struct {
	MessageData             map[string]interface{} `json:"message_data,omitempty"`
	TransactionalMessageID  string                 `json:"transactional_message_id,omitempty"`
	Identifiers             map[string]string      `json:"identifiers"`
	To                      string                 `json:"to,omitempty"`
	DisableMessageRetention *bool                  `json:"disable_message_retention,omitempty"`
	SendToUnsubscribed      *bool                  `json:"send_to_unsubscribed,omitempty"`
	QueueDraft              *bool                  `json:"queue_draft,omitempty"`
	SendAt                  *int64                 `json:"send_at,omitempty"`
	Language                *string                `json:"language,omitempty"`

	Title         string          `json:"title,omitempty"`
	Message       string          `json:"message,omitempty"`
	ImageURL      string          `json:"image_url,omitempty"`
	Link          string          `json:"link,omitempty"`
	CustomData    json.RawMessage `json:"custom_data,omitempty"`
	CustomPayload json.RawMessage `json:"custom_payload,omitempty"`
	Device        *deviceV2       `json:"custom_device,omitempty"`
}

type SendPushResponse struct {
	TransactionalResponse
}

// SendPush sends a single transactional push using the Customer.io transactional API
func (c *APIClient) SendPush(ctx context.Context, req *SendPushRequest) (*SendPushResponse, error) {
	resp, err := c.sendTransactional(ctx, TransactionalTypePush, req)
	if err != nil {
		return nil, err
	}

	return &SendPushResponse{
		*resp,
	}, nil
}
