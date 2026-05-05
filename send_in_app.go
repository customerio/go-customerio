package customerio

import (
	"context"
)

type SendInAppRequest struct {
	MessageData             map[string]any    `json:"message_data,omitempty"`
	TransactionalMessageID  string            `json:"transactional_message_id,omitempty"`
	Identifiers             map[string]string `json:"identifiers"`
	DisableMessageRetention *bool             `json:"disable_message_retention,omitempty"`
	QueueDraft              *bool             `json:"queue_draft,omitempty"`
	SendAt                  *int64            `json:"send_at,omitempty"`
	Language                *string           `json:"language,omitempty"`
}

type SendInAppResponse struct {
	TransactionalResponse
}

// SendInApp sends a single transactional in-app message using the Customer.io transactional API
func (c *APIClient) SendInApp(ctx context.Context, req *SendInAppRequest) (*SendInAppResponse, error) {
	resp, err := c.sendTransactional(ctx, TransactionalTypeInApp, req)
	if err != nil {
		return nil, err
	}

	return &SendInAppResponse{
		*resp,
	}, nil
}
