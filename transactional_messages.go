package customerio

import "context"

// TransactionalMessage is a transactional message template's metadata (not
// its content — see TransactionalMessageContent).
type TransactionalMessage struct {
	ID                 int    `json:"id"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	SendToUnsubscribed bool   `json:"send_to_unsubscribed"`
	LinkTracking       bool   `json:"link_tracking"`
	OpenTracking       bool   `json:"open_tracking"`
	HideMessageBody    bool   `json:"hide_message_body"`
	QueueDrafts        bool   `json:"queue_drafts"`
	CreatedAt          int64  `json:"created_at,omitempty"`
	UpdatedAt          int64  `json:"updated_at,omitempty"`
}

// ListTransactionalMessagesResponse is the decoded shape of GET /v1/transactional.
type ListTransactionalMessagesResponse struct {
	Messages []TransactionalMessage `json:"messages"`
}

// ListTransactionalMessages returns every transactional message template
// defined in the workspace.
// See https://docs.customer.io/api/app/#operation/listTransactionalMessages
func (c *APIClient) ListTransactionalMessages(ctx context.Context) (*ListTransactionalMessagesResponse, error) {
	var resp ListTransactionalMessagesResponse
	if err := c.doJSON(ctx, "GET", "/v1/transactional", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TransactionalMessageResponse wraps a single TransactionalMessage.
type TransactionalMessageResponse struct {
	Message TransactionalMessage `json:"message"`
}

// GetTransactionalMessage returns one transactional message template's metadata.
// See https://docs.customer.io/api/app/#operation/getTransactionalMessage
func (c *APIClient) GetTransactionalMessage(ctx context.Context, transactionalID string) (*TransactionalMessageResponse, error) {
	if transactionalID == "" {
		return nil, ParamError{Param: "transactionalID"}
	}

	var resp TransactionalMessageResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/transactional/%s", transactionalID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TransactionalMessageContent is one language variant's rendered content
// for a transactional message template.
type TransactionalMessageContent struct {
	ID            int               `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	Created       int64             `json:"created,omitempty"`
	Updated       int64             `json:"updated,omitempty"`
	Body          string            `json:"body,omitempty"`
	BodyAMP       string            `json:"body_amp,omitempty"`
	Language      string            `json:"language,omitempty"`
	Type          string            `json:"type,omitempty"`
	From          string            `json:"from,omitempty"`
	FromID        *int              `json:"from_id,omitempty"`
	ReplyTo       string            `json:"reply_to,omitempty"`
	ReplyToID     *int              `json:"reply_to_id,omitempty"`
	Preprocessor  string            `json:"preprocessor,omitempty"`
	Recipient     string            `json:"recipient,omitempty"`
	Subject       string            `json:"subject,omitempty"`
	PreheaderText string            `json:"preheader_text,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	CC            string            `json:"cc,omitempty"`
	BCC           string            `json:"bcc,omitempty"`
	FakeBCC       *bool             `json:"fake_bcc,omitempty"`
}

// TransactionalMessageContentsResponse is the decoded shape of
// GET /v1/transactional/{id}/contents.
type TransactionalMessageContentsResponse struct {
	Contents []TransactionalMessageContent `json:"contents"`
}

// GetTransactionalMessageContents returns every language variant's content
// for a transactional message template.
// See https://docs.customer.io/api/app/#operation/getTransactionalMessageContents
func (c *APIClient) GetTransactionalMessageContents(ctx context.Context, transactionalID string) (*TransactionalMessageContentsResponse, error) {
	if transactionalID == "" {
		return nil, ParamError{Param: "transactionalID"}
	}

	var resp TransactionalMessageContentsResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/transactional/%s/contents", transactionalID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TransactionalMessageContentResponse wraps the content variant(s) returned
// by the language and content-update endpoints.
type TransactionalMessageContentResponse struct {
	Content []TransactionalMessageContent `json:"content"`
}

// GetTransactionalMessageLanguage returns one language variant's content.
// See https://docs.customer.io/api/app/#operation/getTransactionalMessageLanguage
func (c *APIClient) GetTransactionalMessageLanguage(ctx context.Context, transactionalID, language string) (*TransactionalMessageContentResponse, error) {
	if transactionalID == "" {
		return nil, ParamError{Param: "transactionalID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	var resp TransactionalMessageContentResponse
	requestPath := formatPath("/v1/transactional/%s/language/%s", transactionalID, language)
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateTransactionalMessageLanguage creates or updates a language variant's
// content. data is free-form (subject/body/etc — whichever fields you want
// to set), matching the App API's own untyped request body for this endpoint.
// See https://docs.customer.io/api/app/#operation/updateTransactionalMessageLanguage
func (c *APIClient) UpdateTransactionalMessageLanguage(ctx context.Context, transactionalID, language string, data map[string]any) (*TransactionalMessageContentResponse, error) {
	if transactionalID == "" {
		return nil, ParamError{Param: "transactionalID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	var resp TransactionalMessageContentResponse
	requestPath := formatPath("/v1/transactional/%s/language/%s", transactionalID, language)
	if err := c.doJSON(ctx, "PUT", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TransactionalDeliveriesOptions filters GetTransactionalMessageDeliveries.
type TransactionalDeliveriesOptions struct {
	PaginationOptions
	Metric              string
	StartTS             int64
	EndTS               int64
	GetTrackedResponses *bool
}

// TransactionalDeliveriesResponse is the decoded shape of
// GET /v1/transactional/{id}/messages.
type TransactionalDeliveriesResponse struct {
	Messages []Message `json:"messages"`
	Next     string    `json:"next,omitempty"`
}

// GetTransactionalMessageDeliveries returns the messages sent from a
// transactional message template.
// See https://docs.customer.io/api/app/#operation/getTransactionalMessageDeliveries
func (c *APIClient) GetTransactionalMessageDeliveries(ctx context.Context, transactionalID string, opts TransactionalDeliveriesOptions) (*TransactionalDeliveriesResponse, error) {
	if transactionalID == "" {
		return nil, ParamError{Param: "transactionalID"}
	}

	q := opts.PaginationOptions.apply(newQuery()).
		setString("metric", opts.Metric).
		setInt64("start_ts", opts.StartTS).
		setInt64("end_ts", opts.EndTS).
		setBool("get_tracked_responses", opts.GetTrackedResponses)

	requestPath := formatPath("/v1/transactional/%s/messages", transactionalID) + q.String()

	var resp TransactionalDeliveriesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MetricsPeriod buckets a metrics time series into hours/days/weeks/months.
type MetricsPeriod string

const (
	MetricsPeriodHours  MetricsPeriod = "hours"
	MetricsPeriodDays   MetricsPeriod = "days"
	MetricsPeriodWeeks  MetricsPeriod = "weeks"
	MetricsPeriodMonths MetricsPeriod = "months"
)

// TransactionalMetricsOptions controls the time series window for
// GetTransactionalMessageMetrics.
type TransactionalMetricsOptions struct {
	Period MetricsPeriod
	Steps  int
}

// MetricSeries is a set of per-bucket counts, one entry per event kind that
// applies to the resource being measured. This is the full superset of
// fields the App API reports across every metered resource (transactional
// messages, campaigns, campaign/broadcast actions, broadcasts,
// newsletters/newsletter content); which fields are actually populated
// depends on the resource and channel. The StatusCode* fields only apply
// to webhook actions.
type MetricSeries struct {
	Attempted         []int `json:"attempted,omitempty"`
	Bounced           []int `json:"bounced,omitempty"`
	Clicked           []int `json:"clicked,omitempty"`
	HumanClicked      []int `json:"human_clicked,omitempty"`
	PrefetchClicked   []int `json:"prefetch_clicked,omitempty"`
	Converted         []int `json:"converted,omitempty"`
	Created           []int `json:"created,omitempty"`
	Deferred          []int `json:"deferred,omitempty"`
	Delivered         []int `json:"delivered,omitempty"`
	Drafted           []int `json:"drafted,omitempty"`
	Failed            []int `json:"failed,omitempty"`
	Opened            []int `json:"opened,omitempty"`
	HumanOpened       []int `json:"human_opened,omitempty"`
	PrefetchOpened    []int `json:"prefetch_opened,omitempty"`
	Sent              []int `json:"sent,omitempty"`
	Spammed           []int `json:"spammed,omitempty"`
	Suppressed        []int `json:"suppressed,omitempty"`
	Undeliverable     []int `json:"undeliverable,omitempty"`
	TopicUnsubscribed []int `json:"topic_unsubscribed,omitempty"`
	Unsubscribed      []int `json:"unsubscribed,omitempty"`
	StatusCode2xx     []int `json:"2xx,omitempty"`
	StatusCode3xx     []int `json:"3xx,omitempty"`
	StatusCode4xx     []int `json:"4xx,omitempty"`
	StatusCode5xx     []int `json:"5xx,omitempty"`
}

// TransactionalMetricsResponse wraps a MetricSeries report. Despite the
// name (kept for API stability with where it was first introduced), the
// App API returns this identical {metric:{series:{...}}} shape for every
// metered resource, so it's reused as-is by campaign/broadcast/newsletter
// metrics methods rather than duplicated per resource.
type TransactionalMetricsResponse struct {
	Metric struct {
		Series MetricSeries `json:"series"`
	} `json:"metric"`
}

// GetTransactionalMessageMetrics returns send/delivery/engagement counts
// over time for a transactional message template.
// See https://docs.customer.io/api/app/#operation/getTransactionalMessageMetrics
func (c *APIClient) GetTransactionalMessageMetrics(ctx context.Context, transactionalID string, opts TransactionalMetricsOptions) (*TransactionalMetricsResponse, error) {
	if transactionalID == "" {
		return nil, ParamError{Param: "transactionalID"}
	}

	q := newQuery().setString("period", string(opts.Period)).setInt("steps", opts.Steps)
	requestPath := formatPath("/v1/transactional/%s/metrics", transactionalID) + q.String()

	var resp TransactionalMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TransactionalLinkMetricsOptions controls GetTransactionalMessageLinkMetrics.
type TransactionalLinkMetricsOptions struct {
	TransactionalMetricsOptions
	Unique *bool
}

// LinkMetricSeries is per-bucket click counts for one link.
type LinkMetricSeries struct {
	Clicked        []int `json:"clicked,omitempty"`
	HumanClicked   []int `json:"human_clicked,omitempty"`
	MachineClicked []int `json:"machine_clicked,omitempty"`
}

// LinkMetric is one tracked link's click metrics.
type LinkMetric struct {
	Link struct {
		ID   string `json:"id,omitempty"`
		Href string `json:"href,omitempty"`
	} `json:"link"`
	Metric struct {
		Series LinkMetricSeries `json:"series"`
	} `json:"metric"`
}

// TransactionalLinkMetricsResponse wraps per-link click metrics. Like
// TransactionalMetricsResponse, this shape is shared across every metered
// resource's .../metrics/links endpoint, not just transactional messages.
type TransactionalLinkMetricsResponse struct {
	Links []LinkMetric `json:"links"`
}

// GetTransactionalMessageLinkMetrics returns per-link click metrics for a
// transactional message template.
// See https://docs.customer.io/api/app/#operation/getTransactionalMessageLinkMetrics
func (c *APIClient) GetTransactionalMessageLinkMetrics(ctx context.Context, transactionalID string, opts TransactionalLinkMetricsOptions) (*TransactionalLinkMetricsResponse, error) {
	if transactionalID == "" {
		return nil, ParamError{Param: "transactionalID"}
	}

	q := newQuery().
		setString("period", string(opts.Period)).
		setInt("steps", opts.Steps).
		setBool("unique", opts.Unique)
	requestPath := formatPath("/v1/transactional/%s/metrics/links", transactionalID) + q.String()

	var resp TransactionalLinkMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateTransactionalMessageContent updates one content variant by id. data
// is free-form, matching the App API's own untyped request body.
// See https://docs.customer.io/api/app/#operation/updateTransactionalMessageContent
func (c *APIClient) UpdateTransactionalMessageContent(ctx context.Context, transactionalID, contentID string, data map[string]any) (*TransactionalMessageContentResponse, error) {
	if transactionalID == "" {
		return nil, ParamError{Param: "transactionalID"}
	}
	if contentID == "" {
		return nil, ParamError{Param: "contentID"}
	}

	var resp TransactionalMessageContentResponse
	requestPath := formatPath("/v1/transactional/%s/content/%s", transactionalID, contentID)
	if err := c.doJSON(ctx, "PUT", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
