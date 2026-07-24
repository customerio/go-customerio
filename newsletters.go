package customerio

import (
	"context"
	"encoding/json"
)

// Newsletter is a one-off or recurring bulk send to a segment. Unlike
// Campaign/Broadcast, newsletters have no documented State field.
type Newsletter struct {
	ID                  int      `json:"id"`
	DeduplicateID       string   `json:"deduplicate_id,omitempty"`
	ContentIDs          []int    `json:"content_ids,omitempty"`
	Name                string   `json:"name,omitempty"`
	SentAt              int64    `json:"sent_at,omitempty"`
	Created             int64    `json:"created,omitempty"`
	Updated             int64    `json:"updated,omitempty"`
	Type                string   `json:"type,omitempty"`
	Tags                []string `json:"tags,omitempty"`
	RecipientSegmentIDs []int    `json:"recipient_segment_ids,omitempty"`
	SubscriptionTopicID *int     `json:"subscription_topic_id,omitempty"`
}

// NewsletterResponse wraps a single Newsletter.
type NewsletterResponse struct {
	Newsletter Newsletter `json:"newsletter"`
}

// ListNewslettersOptions filters ListNewsletters.
type ListNewslettersOptions struct {
	PaginationOptions
	Sort SortDirection
}

// ListNewslettersResponse is the decoded shape of GET /v1/newsletters.
type ListNewslettersResponse struct {
	Newsletters []Newsletter `json:"newsletters"`
}

// ListNewsletters returns every newsletter defined in the workspace.
// See https://docs.customer.io/api/app/#operation/listNewsletters
func (c *APIClient) ListNewsletters(ctx context.Context, opts ListNewslettersOptions) (*ListNewslettersResponse, error) {
	q := opts.PaginationOptions.apply(newQuery()).setString("sort", string(opts.Sort))
	requestPath := "/v1/newsletters" + q.String()

	var resp ListNewslettersResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateNewsletter creates a newsletter. data is free-form, matching the
// App API's own untyped request body for this endpoint.
// See https://docs.customer.io/api/app/#operation/createNewsletter
func (c *APIClient) CreateNewsletter(ctx context.Context, data map[string]any) (*NewsletterResponse, error) {
	var resp NewsletterResponse
	if err := c.doJSON(ctx, "POST", "/v1/newsletters", data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNewsletter returns one newsletter by id.
// See https://docs.customer.io/api/app/#operation/getNewsletter
func (c *APIClient) GetNewsletter(ctx context.Context, newsletterID string) (*NewsletterResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}

	var resp NewsletterResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/newsletters/%s", newsletterID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteNewsletter deletes a newsletter by id.
// See https://docs.customer.io/api/app/#operation/deleteNewsletter
func (c *APIClient) DeleteNewsletter(ctx context.Context, newsletterID string) error {
	if newsletterID == "" {
		return ParamError{Param: "newsletterID"}
	}

	return c.doJSON(ctx, "DELETE", formatPath("/v1/newsletters/%s", newsletterID), nil, nil, 200, 204)
}

// NewsletterContent is one language/variant's rendered content for a newsletter.
type NewsletterContent struct {
	ID            int               `json:"id,omitempty"`
	NewsletterID  json.Number       `json:"newsletter_id,omitempty"`
	DeduplicateID string            `json:"deduplicate_id,omitempty"`
	Name          string            `json:"name,omitempty"`
	Layout        string            `json:"layout,omitempty"`
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
	CC            string            `json:"cc,omitempty"`
	BCC           string            `json:"bcc,omitempty"`
	FakeBCC       *bool             `json:"fake_bcc,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
}

// NewsletterContentsResponse is the decoded shape of GET /v1/newsletters/{id}/contents.
type NewsletterContentsResponse struct {
	Contents []NewsletterContent `json:"contents"`
}

// GetNewsletterContents returns every content variant for a newsletter.
// See https://docs.customer.io/api/app/#operation/getNewsletterContents
func (c *APIClient) GetNewsletterContents(ctx context.Context, newsletterID string) (*NewsletterContentsResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}

	var resp NewsletterContentsResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/newsletters/%s/contents", newsletterID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// NewsletterContentResponse wraps a single NewsletterContent — shared by
// the content, language, and test-group-language endpoints below, which
// all return the same object shape.
type NewsletterContentResponse struct {
	Content NewsletterContent `json:"content"`
}

// GetNewsletterContent returns one content variant by id.
// See https://docs.customer.io/api/app/#operation/getNewsletterContent
func (c *APIClient) GetNewsletterContent(ctx context.Context, newsletterID, contentID string) (*NewsletterContentResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}
	if contentID == "" {
		return nil, ParamError{Param: "contentID"}
	}

	requestPath := formatPath("/v1/newsletters/%s/contents/%s", newsletterID, contentID)

	var resp NewsletterContentResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateNewsletterContent updates a content variant. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateNewsletterContent
func (c *APIClient) UpdateNewsletterContent(ctx context.Context, newsletterID, contentID string, data map[string]any) (*NewsletterContentResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}
	if contentID == "" {
		return nil, ParamError{Param: "contentID"}
	}

	requestPath := formatPath("/v1/newsletters/%s/contents/%s", newsletterID, contentID)

	var resp NewsletterContentResponse
	if err := c.doJSON(ctx, "PUT", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNewsletterContentMetrics returns send/delivery/engagement counts over
// time for one newsletter content variant.
// See https://docs.customer.io/api/app/#operation/getNewsletterContentMetrics
func (c *APIClient) GetNewsletterContentMetrics(ctx context.Context, newsletterID, contentID string, opts TransactionalMetricsOptions) (*TransactionalMetricsResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}
	if contentID == "" {
		return nil, ParamError{Param: "contentID"}
	}

	q := newQuery().setString("period", string(opts.Period)).setInt("steps", opts.Steps)
	requestPath := formatPath("/v1/newsletters/%s/contents/%s/metrics", newsletterID, contentID) + q.String()

	var resp TransactionalMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNewsletterContentMetricsLinks returns per-link click metrics for one
// newsletter content variant.
// See https://docs.customer.io/api/app/#operation/getNewsletterContentMetricsLinks
func (c *APIClient) GetNewsletterContentMetricsLinks(ctx context.Context, newsletterID, contentID string, opts TransactionalLinkMetricsOptions) (*TransactionalLinkMetricsResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}
	if contentID == "" {
		return nil, ParamError{Param: "contentID"}
	}

	q := newQuery().
		setString("period", string(opts.Period)).
		setInt("steps", opts.Steps).
		setBool("unique", opts.Unique)
	requestPath := formatPath("/v1/newsletters/%s/contents/%s/metrics/links", newsletterID, contentID) + q.String()

	var resp TransactionalLinkMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNewsletterMetrics returns send/delivery/engagement counts over time
// for a whole newsletter.
// See https://docs.customer.io/api/app/#operation/getNewsletterMetrics
func (c *APIClient) GetNewsletterMetrics(ctx context.Context, newsletterID string, opts TransactionalMetricsOptions) (*TransactionalMetricsResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}

	q := newQuery().setString("period", string(opts.Period)).setInt("steps", opts.Steps)
	requestPath := formatPath("/v1/newsletters/%s/metrics", newsletterID) + q.String()

	var resp TransactionalMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNewsletterMetricsLinks returns per-link click metrics for a whole newsletter.
// See https://docs.customer.io/api/app/#operation/getNewsletterMetricsLinks
func (c *APIClient) GetNewsletterMetricsLinks(ctx context.Context, newsletterID string, opts TransactionalLinkMetricsOptions) (*TransactionalLinkMetricsResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}

	q := newQuery().
		setString("period", string(opts.Period)).
		setInt("steps", opts.Steps).
		setBool("unique", opts.Unique)
	requestPath := formatPath("/v1/newsletters/%s/metrics/links", newsletterID) + q.String()

	var resp TransactionalLinkMetricsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// NewsletterMessagesOptions filters GetNewsletterMessages.
type NewsletterMessagesOptions struct {
	PaginationOptions
	Metric              string
	Type                NewsletterChannelType
	StartTS             int64
	EndTS               int64
	GetTrackedResponses *bool
}

// GetNewsletterMessages returns messages sent from a newsletter.
// See https://docs.customer.io/api/app/#operation/getNewsletterMessages
func (c *APIClient) GetNewsletterMessages(ctx context.Context, newsletterID string, opts NewsletterMessagesOptions) (*MessagesResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}

	q := opts.PaginationOptions.apply(newQuery()).
		setString("metric", opts.Metric).
		setString("type", string(opts.Type)).
		setInt64("start_ts", opts.StartTS).
		setInt64("end_ts", opts.EndTS).
		setBool("get_tracked_responses", opts.GetTrackedResponses)
	requestPath := formatPath("/v1/newsletters/%s/messages", newsletterID) + q.String()

	var resp MessagesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendNewsletter sends a newsletter immediately. data is free-form.
// See https://docs.customer.io/api/app/#operation/sendNewsletter
func (c *APIClient) SendNewsletter(ctx context.Context, newsletterID string, data map[string]any) (*NewsletterResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}

	var resp NewsletterResponse
	if err := c.doJSON(ctx, "POST", formatPath("/v1/newsletters/%s/send", newsletterID), data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ScheduleNewsletter schedules a newsletter to send later. data is free-form.
// See https://docs.customer.io/api/app/#operation/scheduleNewsletter
func (c *APIClient) ScheduleNewsletter(ctx context.Context, newsletterID string, data map[string]any) (*NewsletterResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}

	var resp NewsletterResponse
	if err := c.doJSON(ctx, "POST", formatPath("/v1/newsletters/%s/schedule", newsletterID), data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateNewsletterLanguage adds a language variant to a newsletter. data is
// free-form (must include the language tag — see the App API docs for the
// expected field name).
// See https://docs.customer.io/api/app/#operation/createNewsletterLanguage
func (c *APIClient) CreateNewsletterLanguage(ctx context.Context, newsletterID string, data map[string]any) (*NewsletterContentResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}

	var resp NewsletterContentResponse
	if err := c.doJSON(ctx, "POST", formatPath("/v1/newsletters/%s/language", newsletterID), data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNewsletterLanguage returns one language variant of a newsletter.
// See https://docs.customer.io/api/app/#operation/getNewsletterLanguage
func (c *APIClient) GetNewsletterLanguage(ctx context.Context, newsletterID, language string) (*NewsletterContentResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/newsletters/%s/language/%s", newsletterID, language)

	var resp NewsletterContentResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateNewsletterLanguage updates a language variant of a newsletter. data
// is free-form.
// See https://docs.customer.io/api/app/#operation/updateNewsletterLanguage
func (c *APIClient) UpdateNewsletterLanguage(ctx context.Context, newsletterID, language string, data map[string]any) (*NewsletterContentResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/newsletters/%s/language/%s", newsletterID, language)

	var resp NewsletterContentResponse
	if err := c.doJSON(ctx, "PUT", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteNewsletterLanguage removes a language variant from a newsletter.
// See https://docs.customer.io/api/app/#operation/deleteNewsletterLanguage
func (c *APIClient) DeleteNewsletterLanguage(ctx context.Context, newsletterID, language string) error {
	if newsletterID == "" {
		return ParamError{Param: "newsletterID"}
	}
	if language == "" {
		return ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/newsletters/%s/language/%s", newsletterID, language)
	return c.doJSON(ctx, "DELETE", requestPath, nil, nil, 200, 204)
}

// NewsletterTestGroup is one A/B test variant of a newsletter. ContentIDs
// is decoded leniently (json.Number elements): the schema declares it as
// string[], but example responses show integers.
type NewsletterTestGroup struct {
	ID         int           `json:"id"`
	Name       string        `json:"name,omitempty"`
	Label      string        `json:"label,omitempty"`
	Winner     bool          `json:"winner,omitempty"`
	ContentIDs []json.Number `json:"content_ids,omitempty"`
}

// NewsletterTestGroupsResponse is the decoded shape of GET /v1/newsletters/{id}/test_groups.
type NewsletterTestGroupsResponse struct {
	TestGroups []NewsletterTestGroup `json:"test_groups"`
}

// GetNewsletterTestGroups returns every A/B test group defined for a newsletter.
// See https://docs.customer.io/api/app/#operation/getNewsletterTestGroups
func (c *APIClient) GetNewsletterTestGroups(ctx context.Context, newsletterID string) (*NewsletterTestGroupsResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}

	var resp NewsletterTestGroupsResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/newsletters/%s/test_groups", newsletterID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// NewsletterTestGroupResponse wraps a single NewsletterTestGroup.
type NewsletterTestGroupResponse struct {
	TestGroup NewsletterTestGroup `json:"test_group"`
}

// CreateNewsletterTestGroup adds a new (empty) A/B test group to a
// newsletter — this endpoint takes no request body.
// See https://docs.customer.io/api/app/#operation/createNewsletterTestGroup
func (c *APIClient) CreateNewsletterTestGroup(ctx context.Context, newsletterID string) (*NewsletterTestGroupResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}

	var resp NewsletterTestGroupResponse
	if err := c.doJSON(ctx, "POST", formatPath("/v1/newsletters/%s/test_groups", newsletterID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateNewsletterTestGroupLanguage adds a language variant to a test
// group's content. data is free-form. Note the wire path uses singular
// "test_group" here (and in the three methods below), while
// GetNewsletterTestGroups/CreateNewsletterTestGroup use plural
// "test_groups" — a real App API path inconsistency, not a typo.
// See https://docs.customer.io/api/app/#operation/createNewsletterTestGroupLanguage
func (c *APIClient) CreateNewsletterTestGroupLanguage(ctx context.Context, newsletterID, testGroupID string, data map[string]any) (*NewsletterContentResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}
	if testGroupID == "" {
		return nil, ParamError{Param: "testGroupID"}
	}

	requestPath := formatPath("/v1/newsletters/%s/test_group/%s/language", newsletterID, testGroupID)

	var resp NewsletterContentResponse
	if err := c.doJSON(ctx, "POST", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNewsletterTestGroupLanguage returns one language variant of a test group.
// See https://docs.customer.io/api/app/#operation/getNewsletterTestGroupLanguage
func (c *APIClient) GetNewsletterTestGroupLanguage(ctx context.Context, newsletterID, testGroupID, language string) (*NewsletterContentResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}
	if testGroupID == "" {
		return nil, ParamError{Param: "testGroupID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/newsletters/%s/test_group/%s/language/%s", newsletterID, testGroupID, language)

	var resp NewsletterContentResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateNewsletterTestGroupLanguage updates a language variant of a test
// group. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateNewsletterTestGroupLanguage
func (c *APIClient) UpdateNewsletterTestGroupLanguage(ctx context.Context, newsletterID, testGroupID, language string, data map[string]any) (*NewsletterContentResponse, error) {
	if newsletterID == "" {
		return nil, ParamError{Param: "newsletterID"}
	}
	if testGroupID == "" {
		return nil, ParamError{Param: "testGroupID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/newsletters/%s/test_group/%s/language/%s", newsletterID, testGroupID, language)

	var resp NewsletterContentResponse
	if err := c.doJSON(ctx, "PUT", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteNewsletterTestGroupLanguage removes a language variant from a test group.
// See https://docs.customer.io/api/app/#operation/deleteNewsletterTestGroupLanguage
func (c *APIClient) DeleteNewsletterTestGroupLanguage(ctx context.Context, newsletterID, testGroupID, language string) error {
	if newsletterID == "" {
		return ParamError{Param: "newsletterID"}
	}
	if testGroupID == "" {
		return ParamError{Param: "testGroupID"}
	}
	if language == "" {
		return ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/newsletters/%s/test_group/%s/language/%s", newsletterID, testGroupID, language)
	return c.doJSON(ctx, "DELETE", requestPath, nil, nil, 200, 204)
}
