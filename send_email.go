package customerio

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
)

type SendEmailRequest struct {
	MessageData             map[string]interface{} `json:"message_data,omitempty"`
	TransactionalMessageID  string                 `json:"transactional_message_id,omitempty"`
	Identifiers             map[string]string      `json:"identifiers"`
	Headers                 map[string]string      `json:"headers,omitempty"`
	From                    string                 `json:"from,omitempty"`
	To                      string                 `json:"to,omitempty"`
	ReplyTo                 string                 `json:"reply_to,omitempty"`
	BCC                     string                 `json:"bcc,omitempty"`
	Subject                 string                 `json:"subject,omitempty"`
	Preheader               string                 `json:"preheader,omitempty"`
	Body                    string                 `json:"body,omitempty"`
	PlaintextBody           string                 `json:"body_plain,omitempty"`
	AMPBody                 string                 `json:"body_amp,omitempty"`
	FakeBCC                 *bool                  `json:"fake_bcc,omitempty"`
	Attachments             map[string]string      `json:"attachments,omitempty"`
	DisableMessageRetention *bool                  `json:"disable_message_retention,omitempty"`
	SendToUnsubscribed      *bool                  `json:"send_to_unsubscribed,omitempty"`
	EnableTracking          *bool                  `json:"tracked,omitempty"`
	QueueDraft              *bool                  `json:"queue_draft,omitempty"`
	DisableCSSPreprocessing *bool                  `json:"disable_css_preprocessing,omitempty"`
	SendAt                  *int64                 `json:"send_at,omitempty"`
	Language                *string                `json:"language,omitempty"`
}

var ErrAttachmentExists = errors.New("attachment with this name already exists")

func (e *SendEmailRequest) Attach(name string, value io.Reader) error {
	if e.Attachments == nil {
		e.Attachments = map[string]string{}
	}
	if _, ok := e.Attachments[name]; ok {
		return ErrAttachmentExists
	}

	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	if _, err := io.Copy(enc, value); err != nil {
		return err
	}
	if err := enc.Close(); err != nil {
		return err
	}

	e.Attachments[name] = buf.String()
	return nil
}

type SendEmailResponse struct {
	TransactionalResponse
}

// SendEmail sends a single transactional email using the Customer.io transactional API
func (c *APIClient) SendEmail(ctx context.Context, req *SendEmailRequest) (*SendEmailResponse, error) {
	resp, err := c.sendTransactional(ctx, TransactionalTypeEmail, req)
	if err != nil {
		return nil, err
	}

	return &SendEmailResponse{
		*resp,
	}, nil
}
