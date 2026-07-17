package customerio

import (
	"context"
)

type SendWhatsAppRequest struct {
	MessageData             map[string]any    `json:"message_data,omitempty"`
	TransactionalMessageID  string            `json:"transactional_message_id,omitempty"`
	Identifiers             map[string]string `json:"identifiers"`
	DisableMessageRetention *bool             `json:"disable_message_retention,omitempty"`
	SendToUnsubscribed      *bool             `json:"send_to_unsubscribed,omitempty"`
	QueueDraft              *bool             `json:"queue_draft,omitempty"`
	SendAt                  *int64            `json:"send_at,omitempty"`
	Language                *string           `json:"language,omitempty"`

	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	Tracked *bool  `json:"tracked,omitempty"`
}

type SendWhatsAppResponse struct {
	TransactionalResponse
}

// SendWhatsApp sends a single transactional WhatsApp message using the Customer.io transactional API
func (c *APIClient) SendWhatsApp(ctx context.Context, req *SendWhatsAppRequest) (*SendWhatsAppResponse, error) {
	resp, err := c.sendTransactional(ctx, TransactionalTypeWhatsApp, req)
	if err != nil {
		return nil, err
	}

	return &SendWhatsAppResponse{
		*resp,
	}, nil
}
