package customerio

import "context"

// MessagesOptions filters GetMessages, the cross-resource message list —
// unlike GetCustomerMessages/GetTransactionalMessageDeliveries, it isn't
// scoped to a single customer or transactional template; the *_id filters
// narrow it to a specific campaign/action/newsletter/transactional/etc.
type MessagesOptions struct {
	PaginationOptions
	Drafts              *bool
	Metric              string
	Type                string
	CampaignID          string
	ActionID            string
	NewsletterID        string
	TransactionalID     string
	TriggerID           string
	TemplateID          string
	ContentID           string
	StartTS             int64
	EndTS               int64
	Associations        *bool
	GetTrackedResponses *bool
}

// GetMessages returns messages across the whole workspace, most recent
// first, optionally filtered to one campaign/action/newsletter/
// transactional template/trigger/template/content id.
// See https://docs.customer.io/api/app/#operation/getMessages
func (c *APIClient) GetMessages(ctx context.Context, opts MessagesOptions) (*MessagesResponse, error) {
	q := opts.PaginationOptions.apply(newQuery()).
		setBool("drafts", opts.Drafts).
		setString("metric", opts.Metric).
		setString("type", opts.Type).
		setString("campaign_id", opts.CampaignID).
		setString("action_id", opts.ActionID).
		setString("newsletter_id", opts.NewsletterID).
		setString("transactional_id", opts.TransactionalID).
		setString("trigger_id", opts.TriggerID).
		setString("template_id", opts.TemplateID).
		setString("content_id", opts.ContentID).
		setInt64("start_ts", opts.StartTS).
		setInt64("end_ts", opts.EndTS).
		setBool("associations", opts.Associations).
		setBool("get_tracked_responses", opts.GetTrackedResponses)

	requestPath := "/v1/messages" + q.String()

	var resp MessagesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MessageOptions controls GetMessage's response detail.
type MessageOptions struct {
	ArchivedMessage     *bool
	Associations        *bool
	GetTrackedResponses *bool
}

// MessageResponse wraps a single Message.
type MessageResponse struct {
	Message Message `json:"message"`
}

// GetMessage returns one message by id.
// See https://docs.customer.io/api/app/#operation/getMessage
func (c *APIClient) GetMessage(ctx context.Context, messageID string, opts MessageOptions) (*MessageResponse, error) {
	if messageID == "" {
		return nil, ParamError{Param: "messageID"}
	}

	q := newQuery().
		setBool("archived_message", opts.ArchivedMessage).
		setBool("associations", opts.Associations).
		setBool("get_tracked_responses", opts.GetTrackedResponses)
	requestPath := formatPath("/v1/messages/%s", messageID) + q.String()

	var resp MessageResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ArchivedMessage is a message's full rendered content, as delivered.
type ArchivedMessage struct {
	ID            string            `json:"id,omitempty"`
	Body          string            `json:"body,omitempty"`
	From          string            `json:"from,omitempty"`
	ReplyTo       string            `json:"reply_to,omitempty"`
	Recipient     string            `json:"recipient,omitempty"`
	Subject       string            `json:"subject,omitempty"`
	CC            string            `json:"cc,omitempty"`
	BCC           string            `json:"bcc,omitempty"`
	FakeBCC       *bool             `json:"fake_bcc,omitempty"`
	PreheaderText string            `json:"preheader_text,omitempty"`
	URL           string            `json:"url,omitempty"`
	RequestMethod string            `json:"request_method,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	Forgotten     bool              `json:"forgotten,omitempty"`
}

// ArchivedMessageResponse is the decoded shape of
// GET /v1/messages/{id}/archived_message.
type ArchivedMessageResponse struct {
	ArchivedMessage ArchivedMessage `json:"archived_message"`
}

// GetArchivedMessage returns a message's full rendered content. This
// endpoint is rate-limited by the API (10 requests/second) — the returned
// error's status will be 429 if that limit is exceeded.
// See https://docs.customer.io/api/app/#operation/getArchivedMessage
func (c *APIClient) GetArchivedMessage(ctx context.Context, messageID string) (*ArchivedMessageResponse, error) {
	if messageID == "" {
		return nil, ParamError{Param: "messageID"}
	}

	var resp ArchivedMessageResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/messages/%s/archived_message", messageID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
