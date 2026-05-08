package customerio

import (
	"context"
)

type SendInboxMessageRequest struct {
	MessageData             map[string]any    `json:"message_data,omitempty"`
	TransactionalMessageID  string            `json:"transactional_message_id,omitempty"`
	Identifiers             map[string]string `json:"identifiers"`
	DisableMessageRetention *bool             `json:"disable_message_retention,omitempty"`
	QueueDraft              *bool             `json:"queue_draft,omitempty"`
	SendAt                  *int64            `json:"send_at,omitempty"`
	Language                *string           `json:"language,omitempty"`
}

type SendInboxMessageResponse struct {
	TransactionalResponse
}

// SendInboxMessage sends a single transactional inbox message using the Customer.io transactional API
func (c *APIClient) SendInboxMessage(ctx context.Context, req *SendInboxMessageRequest) (*SendInboxMessageResponse, error) {
	resp, err := c.sendTransactional(ctx, TransactionalTypeInboxMessage, req)
	if err != nil {
		return nil, err
	}

	return &SendInboxMessageResponse{
		*resp,
	}, nil
}
