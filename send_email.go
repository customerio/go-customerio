package customerio

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type SendEmailRequest struct {
	MessageData             map[string]interface{} `json:"message_data,omitempty"`
	TransactionalMessageID  string                 `json:"transactional_message_id,omitempty"`
	Identifiers             map[string]string      `json:"identifiers"`
	Headers                 map[string]string      `json:"headers,omitempty"`
	From                    string                 `json:"from,omitempty"`
	FromID                  string                 `json:"from_id,omitempty"`
	To                      string                 `json:"to,omitempty"`
	ReplyTo                 string                 `json:"reply_to,omitempty"`
	ReplyToID               string                 `json:"reply_to_id,omitempty"`
	BCC                     string                 `json:"bcc,omitempty"`
	Subject                 string                 `json:"subject,omitempty"`
	Preheader               string                 `json:"preheader,omitempty"`
	Body                    string                 `json:"body,omitempty"`
	PlaintextBody           string                 `json:"plaintext_body,omitempty"`
	AMPBody                 string                 `json:"amp_body,omitempty"`
	FakeBCC                 *bool                  `json:"fake_bcc,omitempty"`
	Attachments             map[string]string      `json:"attachments,omitempty"`
	DisableMessageRetention *bool                  `json:"disable_message_retention,omitempty"`
	SendToUnsubscribed      *bool                  `json:"send_to_unsubscribed,omitempty"`
	EnableTracking          *bool                  `json:"tracked,omitempty"`
	QueueDraft              *bool                  `json:"queue_draft,omitempty"`
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

	e.Attachments[name] = string(buf.Bytes())
	return nil
}

type SendEmailResponse struct {
	TransactionalResponse
}

// SendEmail sends a single transactional email using the Customer.io transactional API
func (c *APIClient) SendEmail(ctx context.Context, req *SendEmailRequest) (*SendEmailResponse, error) {
	body, statusCode, err := c.doRequest(ctx, "POST", "/v1/send/email", req)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		var meta struct {
			Meta struct {
				Err string `json:"error"`
			} `json:"meta"`
		}
		if err := json.Unmarshal(body, &meta); err != nil {
			return nil, &TransactionalError{
				StatusCode: statusCode,
				Err:        string(body),
			}
		}
		return nil, &TransactionalError{
			StatusCode: statusCode,
			Err:        meta.Meta.Err,
		}
	}

	var result SendEmailResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
