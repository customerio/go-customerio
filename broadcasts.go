package customerio

import "context"

// Broadcast is a one-off (non-recurring) email/push/SMS/etc. send, similar
// to a Campaign but without a Description or trigger/filter segments.
type Broadcast struct {
	ID             int             `json:"id"`
	DeduplicateID  string          `json:"deduplicate_id,omitempty"`
	Name           string          `json:"name,omitempty"`
	Type           string          `json:"type,omitempty"`
	State          string          `json:"state,omitempty"`
	Active         bool            `json:"active,omitempty"`
	Created        int64           `json:"created,omitempty"`
	Updated        int64           `json:"updated,omitempty"`
	FirstStarted   int64           `json:"first_started,omitempty"`
	Tags           []string        `json:"tags,omitempty"`
	Actions        []ActionSummary `json:"actions,omitempty"`
	MsgTemplateIDs []int           `json:"msg_template_ids,omitempty"`
}

// ListBroadcastsResponse is the decoded shape of GET /v1/broadcasts.
type ListBroadcastsResponse struct {
	Broadcasts []Broadcast `json:"broadcasts"`
}

// ListBroadcasts returns every broadcast defined in the workspace.
// See https://docs.customer.io/api/app/#operation/listBroadcasts
func (c *APIClient) ListBroadcasts(ctx context.Context) (*ListBroadcastsResponse, error) {
	var resp ListBroadcastsResponse
	if err := c.doJSON(ctx, "GET", "/v1/broadcasts", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BroadcastResponseWrapper wraps a single Broadcast. Named distinctly from
// BroadcastResponse (trigger_broadcast.go), which is TriggerBroadcast's
// {id} acknowledgment and predates this file.
type BroadcastResponseWrapper struct {
	Broadcast Broadcast `json:"broadcast"`
}

// GetBroadcast returns one broadcast by id.
// See https://docs.customer.io/api/app/#operation/getBroadcast
func (c *APIClient) GetBroadcast(ctx context.Context, broadcastID string) (*BroadcastResponseWrapper, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}

	var resp BroadcastResponseWrapper
	if err := c.doJSON(ctx, "GET", formatPath("/v1/broadcasts/%s", broadcastID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBroadcastActions returns a broadcast's actions. Unlike
// GetCampaignActions, this endpoint takes no pagination parameter.
// See https://docs.customer.io/api/app/#operation/getBroadcastActions
func (c *APIClient) GetBroadcastActions(ctx context.Context, broadcastID string) (*ActionsResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}

	var resp ActionsResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/broadcasts/%s/actions", broadcastID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBroadcastAction returns one broadcast action by id.
// See https://docs.customer.io/api/app/#operation/getBroadcastAction
func (c *APIClient) GetBroadcastAction(ctx context.Context, broadcastID, actionID string) (*ActionResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}

	requestPath := formatPath("/v1/broadcasts/%s/actions/%s", broadcastID, actionID)

	var resp ActionResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateBroadcastAction updates a broadcast action. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateBroadcastAction
func (c *APIClient) UpdateBroadcastAction(ctx context.Context, broadcastID, actionID string, data map[string]any) (*ActionResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}

	requestPath := formatPath("/v1/broadcasts/%s/actions/%s", broadcastID, actionID)

	var resp ActionResponse
	if err := c.doJSON(ctx, "PUT", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBroadcastActionLanguage returns one language variant of a broadcast action.
// See https://docs.customer.io/api/app/#operation/getBroadcastActionLanguage
func (c *APIClient) GetBroadcastActionLanguage(ctx context.Context, broadcastID, actionID, language string) (*ActionResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/broadcasts/%s/actions/%s/language/%s", broadcastID, actionID, language)

	var resp ActionResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateBroadcastActionLanguage creates or updates a language variant of a
// broadcast action. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateBroadcastActionLanguage
func (c *APIClient) UpdateBroadcastActionLanguage(ctx context.Context, broadcastID, actionID, language string, data map[string]any) (*ActionResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/broadcasts/%s/actions/%s/language/%s", broadcastID, actionID, language)

	var resp ActionResponse
	if err := c.doJSON(ctx, "PUT", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBroadcastActionMetrics returns send/delivery/engagement counts over
// time for one broadcast action. Unlike campaign action metrics, this
// endpoint only accepts Period/Steps — no version/res/tz/start/end.
// See https://docs.customer.io/api/app/#operation/getBroadcastActionMetrics
func (c *APIClient) GetBroadcastActionMetrics(ctx context.Context, broadcastID, actionID string, opts TransactionalMetricsOptions) (*TransactionalMetricsResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}

	q := newQuery().setString("period", string(opts.Period)).setInt("steps", opts.Steps)
	requestPath := formatPath("/v1/broadcasts/%s/actions/%s/metrics", broadcastID, actionID) + q.String()

	var resp TransactionalMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBroadcastActionMetricsLinks returns per-link click metrics for one
// broadcast action.
// See https://docs.customer.io/api/app/#operation/getBroadcastActionMetricsLinks
func (c *APIClient) GetBroadcastActionMetricsLinks(ctx context.Context, broadcastID, actionID string, opts TransactionalLinkMetricsOptions) (*TransactionalLinkMetricsResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}
	if actionID == "" {
		return nil, ParamError{Param: "actionID"}
	}

	q := newQuery().
		setString("period", string(opts.Period)).
		setInt("steps", opts.Steps).
		setBool("unique", opts.Unique)
	requestPath := formatPath("/v1/broadcasts/%s/actions/%s/metrics/links", broadcastID, actionID) + q.String()

	var resp TransactionalLinkMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BroadcastMetricsOptions controls GetBroadcastMetrics.
type BroadcastMetricsOptions struct {
	TransactionalMetricsOptions
	Type MetricType
}

// GetBroadcastMetrics returns send/delivery/engagement counts over time for
// a whole broadcast.
// See https://docs.customer.io/api/app/#operation/getBroadcastMetrics
func (c *APIClient) GetBroadcastMetrics(ctx context.Context, broadcastID string, opts BroadcastMetricsOptions) (*TransactionalMetricsResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}

	q := newQuery().
		setString("period", string(opts.Period)).
		setInt("steps", opts.Steps).
		setString("type", string(opts.Type))
	requestPath := formatPath("/v1/broadcasts/%s/metrics", broadcastID) + q.String()

	var resp TransactionalMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBroadcastMetricsLinks returns per-link click metrics for a whole broadcast.
// See https://docs.customer.io/api/app/#operation/getBroadcastMetricsLinks
func (c *APIClient) GetBroadcastMetricsLinks(ctx context.Context, broadcastID string, opts TransactionalLinkMetricsOptions) (*TransactionalLinkMetricsResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}

	q := newQuery().
		setString("period", string(opts.Period)).
		setInt("steps", opts.Steps).
		setBool("unique", opts.Unique)
	requestPath := formatPath("/v1/broadcasts/%s/metrics/links", broadcastID) + q.String()

	var resp TransactionalLinkMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BroadcastMessagesOptions filters GetBroadcastMessages.
type BroadcastMessagesOptions struct {
	PaginationOptions
	Metric              string
	Type                MetricType
	StartTS             int64
	EndTS               int64
	GetTrackedResponses *bool
}

// GetBroadcastMessages returns messages sent from a broadcast.
// See https://docs.customer.io/api/app/#operation/getBroadcastMessages
func (c *APIClient) GetBroadcastMessages(ctx context.Context, broadcastID string, opts BroadcastMessagesOptions) (*MessagesResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}

	q := opts.PaginationOptions.apply(newQuery()).
		setString("metric", opts.Metric).
		setString("type", string(opts.Type)).
		setInt64("start_ts", opts.StartTS).
		setInt64("end_ts", opts.EndTS).
		setBool("get_tracked_responses", opts.GetTrackedResponses)
	requestPath := formatPath("/v1/broadcasts/%s/messages", broadcastID) + q.String()

	var resp MessagesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BroadcastTriggersResponse is the decoded shape of GET /v1/broadcasts/{id}/triggers.
// The per-entry shape isn't independently documented; it's inferred to
// match GetBroadcastTriggerStatus's object, by convention with every other
// list endpoint in this API returning the same shape as its single-item GET.
type BroadcastTriggersResponse struct {
	Triggers []BroadcastTriggerStatus `json:"triggers"`
}

// GetBroadcastTriggers lists every trigger fired for a broadcast. To check
// one trigger's status or errors, see GetBroadcastTriggerStatus and
// GetBroadcastTriggerErrors (trigger_broadcast.go) — those two address the
// same triggers under /v1/campaigns/{broadcastID}/triggers/{triggerID}, a
// real (if surprising) App API path inconsistency, not a typo here.
// See https://docs.customer.io/api/app/#operation/getBroadcastTriggers
func (c *APIClient) GetBroadcastTriggers(ctx context.Context, broadcastID string) (*BroadcastTriggersResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}

	var resp BroadcastTriggersResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/broadcasts/%s/triggers", broadcastID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
