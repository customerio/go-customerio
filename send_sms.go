package customerio

import (
	"context"
)

type SendSMSRequest struct {
	MessageData             map[string]interface{} `json:"message_data,omitempty"`
	TransactionalMessageID  string                 `json:"transactional_message_id,omitempty"`
	Identifiers             map[string]string      `json:"identifiers"`
	DisableMessageRetention *bool                  `json:"disable_message_retention,omitempty"`
	SendToUnsubscribed      *bool                  `json:"send_to_unsubscribed,omitempty"`
	QueueDraft              *bool                  `json:"queue_draft,omitempty"`
	SendAt                  *int64                 `json:"send_at,omitempty"`
	Language                *string                `json:"language,omitempty"`

	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

type SendSMSResponse struct {
	TransactionalResponse
}

// SendPush sends a single transactional push using the Customer.io transactional API
func (c *APIClient) SendSMS(ctx context.Context, req *SendSMSRequest) (*SendSMSResponse, error) {
	resp, err := c.sendTransactional(ctx, TransactionalTypeSMS, req)
	if err != nil {
		return nil, err
	}

	return &SendSMSResponse{
		*resp,
	}, nil
}
