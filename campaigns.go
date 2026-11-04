package customerio

import "context"

// ActionSummary is a campaign or broadcast's abbreviated reference to one
// of its actions, as embedded in the Campaign/Broadcast object's Actions list.
type ActionSummary struct {
	Type string `json:"type,omitempty"`
	ID   int    `json:"id,omitempty"`
}

// Campaign is an email/push/SMS/etc. journey triggered by a segment or event.
type Campaign struct {
	ID                int             `json:"id"`
	DeduplicateID     string          `json:"deduplicate_id,omitempty"`
	Name              string          `json:"name,omitempty"`
	Description       string          `json:"description,omitempty"`
	Type              string          `json:"type,omitempty"`
	State             string          `json:"state,omitempty"`
	Active            bool            `json:"active,omitempty"`
	Created           int64           `json:"created,omitempty"`
	Updated           int64           `json:"updated,omitempty"`
	FirstStarted      int64           `json:"first_started,omitempty"`
	Tags              []string        `json:"tags,omitempty"`
	TriggerSegmentIDs []int           `json:"trigger_segment_ids,omitempty"`
	FilterSegmentIDs  []int           `json:"filter_segment_ids,omitempty"`
	Actions           []ActionSummary `json:"actions,omitempty"`
}

// ListCampaignsResponse is the decoded shape of GET /v1/campaigns.
type ListCampaignsResponse struct {
	Campaigns []Campaign `json:"campaigns"`
}

// ListCampaigns returns every campaign defined in the workspace.
// See https://docs.customer.io/api/app/#operation/listCampaigns
func (c *APIClient) ListCampaigns(ctx context.Context) (*ListCampaignsResponse, error) {
	var resp ListCampaignsResponse
	if err := c.doJSON(ctx, "GET", "/v1/campaigns", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CampaignResponse wraps a single Campaign.
type CampaignResponse struct {
	Campaign Campaign `json:"campaign"`
}

// GetCampaign returns one campaign by id.
// See https://docs.customer.io/api/app/#operation/getCampaign
func (c *APIClient) GetCampaign(ctx context.Context, campaignID string) (*CampaignResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}

	var resp CampaignResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/campaigns/%s", campaignID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Action is one step (email/push/SMS/webhook/...) in a campaign or
// broadcast. BroadcastID names the parent campaign or broadcast id despite
// the field name — that's the wire field name for both resource types.
type Action struct {
	ID            int    `json:"id"`
	BroadcastID   int    `json:"broadcast_id,omitempty"`
	DeduplicateID string `json:"deduplicate_id,omitempty"`
	Name          string `json:"name,omitempty"`
	Layout        string `json:"layout,omitempty"`
	Created       int64  `json:"created,omitempty"`
	Updated       int64  `json:"updated,omitempty"`
	Body          string `json:"body,omitempty"`
	Type          string `json:"type,omitempty"`
	SendingState  string `json:"sending_state,omitempty"`
	Language      string `json:"language,omitempty"`
	From          string `json:"from,omitempty"`
	FromID        *int   `json:"from_id,omitempty"`
	ReplyTo       string `json:"reply_to,omitempty"`
	ReplyToID     *int   `json:"reply_to_id,omitempty"`
	Subject       string `json:"subject,omitempty"`
	Recipient     string `json:"recipient,omitempty"`
	CC            string `json:"cc,omitempty"`
	BCC           string `json:"bcc,omitempty"`
}

// ActionResponse wraps a single Action.
type ActionResponse struct {
	Action Action `json:"action"`
}

// ActionsResponse is the decoded shape of an actions-list endpoint.
type ActionsResponse struct {
	Actions []Action `json:"actions"`
	Next    string   `json:"next,omitempty"`
}

// GetCampaignActions returns a campaign's actions. start is an opaque
// pagination cursor from a previous call, or "" for the first page.
// See https://docs.customer.io/api/app/#operation/getCampaignActions
func (c *APIClient) GetCampaignActions(ctx context.Context, campaignID, start string) (*ActionsResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}

	requestPath := formatPath("/v1/campaigns/%s/actions", campaignID) + newQuery().setString("start", start).String()

	var resp ActionsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCampaignAction returns one campaign action by id.
// See https://docs.customer.io/api/app/#operation/getCampaignAction
func (c *APIClient) GetCampaignAction(ctx context.Context, campaignID, actionID string) (*ActionResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}

	requestPath := formatPath("/v1/campaigns/%s/actions/%s", campaignID, actionID)

	var resp ActionResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCampaignAction updates a campaign action. data is free-form,
// matching the App API's own untyped request body for this endpoint.
// See https://docs.customer.io/api/app/#operation/updateCampaignAction
func (c *APIClient) UpdateCampaignAction(ctx context.Context, campaignID, actionID string, data map[string]any) (*ActionResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}

	requestPath := formatPath("/v1/campaigns/%s/actions/%s", campaignID, actionID)

	var resp ActionResponse
	if err := c.doJSON(ctx, "PUT", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCampaignActionLanguage returns one language variant of a campaign action.
// See https://docs.customer.io/api/app/#operation/getCampaignActionLanguage
func (c *APIClient) GetCampaignActionLanguage(ctx context.Context, campaignID, actionID, language string) (*ActionResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/campaigns/%s/actions/%s/language/%s", campaignID, actionID, language)

	var resp ActionResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCampaignActionLanguage creates or updates a language variant of a
// campaign action. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateCampaignActionLanguage
func (c *APIClient) UpdateCampaignActionLanguage(ctx context.Context, campaignID, actionID, language string, data map[string]any) (*ActionResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/campaigns/%s/actions/%s/language/%s", campaignID, actionID, language)

	var resp ActionResponse
	if err := c.doJSON(ctx, "PUT", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CampaignActionMetricsOptions controls GetCampaignActionMetrics.
type CampaignActionMetricsOptions struct {
	TransactionalMetricsOptions
	Version CampaignMetricsVersion
	Res     MetricResolution
	TZ      string
	Start   int64
	End     int64
}

// GetCampaignActionMetrics returns send/delivery/engagement counts over
// time for one campaign action.
// See https://docs.customer.io/api/app/#operation/getCampaignActionMetrics
func (c *APIClient) GetCampaignActionMetrics(ctx context.Context, campaignID, actionID string, opts CampaignActionMetricsOptions) (*TransactionalMetricsResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}

	q := newQuery().
		setString("version", string(opts.Version)).
		setString("res", string(opts.Res)).
		setString("tz", opts.TZ).
		setInt64("start", opts.Start).
		setInt64("end", opts.End).
		setString("period", string(opts.Period)).
		setInt("steps", opts.Steps)
	requestPath := formatPath("/v1/campaigns/%s/actions/%s/metrics", campaignID, actionID) + q.String()

	var resp TransactionalMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCampaignActionMetricsLinks returns per-link click metrics for one
// campaign action.
// See https://docs.customer.io/api/app/#operation/getCampaignActionMetricsLinks
func (c *APIClient) GetCampaignActionMetricsLinks(ctx context.Context, campaignID, actionID string, opts TransactionalLinkMetricsOptions) (*TransactionalLinkMetricsResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}

	q := newQuery().
		setString("period", string(opts.Period)).
		setInt("steps", opts.Steps).
		setBool("unique", opts.Unique)
	requestPath := formatPath("/v1/campaigns/%s/actions/%s/metrics/links", campaignID, actionID) + q.String()

	var resp TransactionalLinkMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CampaignMetricsOptions controls GetCampaignMetrics.
type CampaignMetricsOptions struct {
	TransactionalMetricsOptions
	Version CampaignMetricsVersion
	Type    MetricType
	Res     MetricResolution
	TZ      string
	Start   int64
	End     int64
}

// GetCampaignMetrics returns send/delivery/engagement counts over time for
// a whole campaign.
// See https://docs.customer.io/api/app/#operation/getCampaignMetrics
func (c *APIClient) GetCampaignMetrics(ctx context.Context, campaignID string, opts CampaignMetricsOptions) (*TransactionalMetricsResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}

	q := newQuery().
		setString("version", string(opts.Version)).
		setString("type", string(opts.Type)).
		setString("res", string(opts.Res)).
		setString("tz", opts.TZ).
		setInt64("start", opts.Start).
		setInt64("end", opts.End).
		setString("period", string(opts.Period)).
		setInt("steps", opts.Steps)
	requestPath := formatPath("/v1/campaigns/%s/metrics", campaignID) + q.String()

	var resp TransactionalMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCampaignMetricsLinks returns per-link click metrics for a whole campaign.
// See https://docs.customer.io/api/app/#operation/getCampaignMetricsLinks
func (c *APIClient) GetCampaignMetricsLinks(ctx context.Context, campaignID string, opts TransactionalLinkMetricsOptions) (*TransactionalLinkMetricsResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}

	q := newQuery().
		setString("period", string(opts.Period)).
		setInt("steps", opts.Steps).
		setBool("unique", opts.Unique)
	requestPath := formatPath("/v1/campaigns/%s/metrics/links", campaignID) + q.String()

	var resp TransactionalLinkMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// JourneyMetricsOptions controls GetCampaignJourneyMetrics. Unlike other
// metrics endpoints, Start, End, and Resolution are all required.
type JourneyMetricsOptions struct {
	Start      int64
	End        int64
	Resolution MetricResolution
}

// JourneyMetricSeries is a funnel report of how customers moved through a
// campaign, one entry per stage.
type JourneyMetricSeries struct {
	Started        []int `json:"started,omitempty"`
	Activated      []int `json:"activated,omitempty"`
	ExitedEarly    []int `json:"exited_early,omitempty"`
	Finished       []int `json:"finished,omitempty"`
	Converted      []int `json:"converted,omitempty"`
	NeverActivated []int `json:"never_activated,omitempty"`
	Messaged       []int `json:"messaged,omitempty"`
}

// JourneyMetricsResponse is the decoded shape of GET /v1/campaigns/{id}/journey_metrics.
type JourneyMetricsResponse struct {
	JourneyMetric JourneyMetricSeries `json:"journey_metric"`
}

// GetCampaignJourneyMetrics returns a funnel report of how customers moved
// through a campaign (started/activated/exited_early/finished/converted/
// never_activated/messaged), bucketed by opts.Resolution.
// See https://docs.customer.io/api/app/#operation/getCampaignJourneyMetrics
func (c *APIClient) GetCampaignJourneyMetrics(ctx context.Context, campaignID string, opts JourneyMetricsOptions) (*JourneyMetricsResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}
	if opts.Start == 0 {
		return nil, ParamError{Param: "start"}
	}
	if opts.End == 0 {
		return nil, ParamError{Param: "end"}
	}
	if opts.Resolution == "" {
		return nil, ParamError{Param: "resolution"}
	}

	// Wire query param is "resolution"; the option field is named Resolution
	// to match (the Node SDK's option field is the abbreviated "res", but the
	// query param it sends is spelled out — kept explicit here instead).
	q := newQuery().setInt64("start", opts.Start).setInt64("end", opts.End).setString("resolution", string(opts.Resolution))
	requestPath := formatPath("/v1/campaigns/%s/journey_metrics", campaignID) + q.String()

	var resp JourneyMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CampaignMessagesOptions filters GetCampaignMessages.
type CampaignMessagesOptions struct {
	PaginationOptions
	Type                MetricType
	Metric              string
	Drafts              *bool
	StartTS             int64
	EndTS               int64
	GetTrackedResponses *bool
}

// GetCampaignMessages returns messages sent from a campaign.
// See https://docs.customer.io/api/app/#operation/getCampaignMessages
func (c *APIClient) GetCampaignMessages(ctx context.Context, campaignID string, opts CampaignMessagesOptions) (*MessagesResponse, error) {
	if campaignID == "" {
		return nil, ParamError{Param: "campaignID"}
	}

	q := opts.PaginationOptions.apply(newQuery()).
		setString("type", string(opts.Type)).
		setString("metric", opts.Metric).
		setBool("drafts", opts.Drafts).
		setInt64("start_ts", opts.StartTS).
		setInt64("end_ts", opts.EndTS).
		setBool("get_tracked_responses", opts.GetTrackedResponses)
	requestPath := formatPath("/v1/campaigns/%s/messages", campaignID) + q.String()

	var resp MessagesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
