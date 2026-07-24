package customerio

import (
	"context"
	"encoding/json"
)

// CustomerIdentifiers identifies a customer by all three identifier kinds
// Customer.io tracks. It appears in search/list results so callers can look
// the customer up by whichever identifier they have on hand.
type CustomerIdentifiers struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	CioID string `json:"cio_id,omitempty"`
}

// searchCustomersRequest is the POST /v1/customers body.
type searchCustomersRequest struct {
	Filter AudienceFilter `json:"filter"`
}

// SearchCustomersResponse is the decoded shape of POST /v1/customers.
type SearchCustomersResponse struct {
	Identifiers []CustomerIdentifiers `json:"identifiers"`
	IDs         []string              `json:"ids"`
	Next        string                `json:"next,omitempty"`
}

// SearchCustomers finds customers matching filter (built with FilterBySegment,
// FilterByAttribute, and the FilterAnd/FilterOr/FilterNot combinators).
// See https://docs.customer.io/api/app/#operation/searchCustomers
func (c *APIClient) SearchCustomers(ctx context.Context, filter AudienceFilter, opts PaginationOptions) (*SearchCustomersResponse, error) {
	requestPath := "/v1/customers" + opts.apply(newQuery()).String()

	var resp SearchCustomersResponse
	if err := c.doJSON(ctx, "POST", requestPath, searchCustomersRequest{Filter: filter}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCustomersByEmailResponse is the decoded shape of GET /v1/customers?email=.
type GetCustomersByEmailResponse struct {
	Results []CustomerIdentifiers `json:"results"`
}

// GetCustomersByEmail looks up every customer profile with the given email —
// Customer.io allows duplicate emails across profiles, so this can return
// more than one result. See https://docs.customer.io/api/app/#operation/getCustomersByEmail
func (c *APIClient) GetCustomersByEmail(ctx context.Context, email string) (*GetCustomersByEmailResponse, error) {
	if email == "" {
		return nil, ParamError{Param: "email"}
	}

	requestPath := "/v1/customers" + newQuery().setString("email", email).String()

	var resp GetCustomersByEmailResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CustomerDevice is a push notification device registered to a customer.
type CustomerDevice struct {
	ID         string         `json:"id"`
	LastUsed   int64          `json:"last_used,omitempty"`
	Platform   string         `json:"platform"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

// Customer is a customer profile as returned by the attributes endpoints.
// Attributes values are always strings on the wire.
type Customer struct {
	ID           string               `json:"id"`
	Identifiers  *CustomerIdentifiers `json:"identifiers,omitempty"`
	Attributes   map[string]string    `json:"attributes"`
	Unsubscribed bool                 `json:"unsubscribed"`
	Devices      []CustomerDevice     `json:"devices,omitempty"`
}

// CustomerAttributesResponse is the decoded shape of GET /v1/customers/{id}/attributes.
type CustomerAttributesResponse struct {
	Customer Customer `json:"customer"`
}

// GetAttributes returns a customer's profile attributes. idType selects
// which kind of identifier id is; the zero value defaults to IdentifierTypeID.
// See https://docs.customer.io/api/app/#operation/getPersonAttributes
func (c *APIClient) GetAttributes(ctx context.Context, id string, idType IdentifierType) (*CustomerAttributesResponse, error) {
	if id == "" {
		return nil, ParamError{Param: "id"}
	}
	if idType == "" {
		idType = IdentifierTypeID
	}

	requestPath := formatPath("/v1/customers/%s/attributes", id) + newQuery().setString("id_type", string(idType)).String()

	var resp CustomerAttributesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CustomerActivitiesOptions filters GetCustomerActivities.
type CustomerActivitiesOptions struct {
	PaginationOptions
	IDType IdentifierType
	Type   string
	Name   string
}

// GetCustomerActivities returns activity records (sends, opens, clicks,
// events, attribute changes, ...) for one customer, most recent first.
// See https://docs.customer.io/api/app/#operation/getPersonActivities
func (c *APIClient) GetCustomerActivities(ctx context.Context, customerID string, opts CustomerActivitiesOptions) (*ActivitiesResponse, error) {
	if customerID == "" {
		return nil, ParamError{Param: "customerID"}
	}

	q := opts.PaginationOptions.apply(newQuery()).
		setString("id_type", string(opts.IDType)).
		setString("type", opts.Type).
		setString("name", opts.Name)

	requestPath := formatPath("/v1/customers/%s/activities", customerID) + q.String()

	var resp ActivitiesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Message is a sent email/push/SMS/etc. and its delivery metadata — the
// shared shape returned by GET /v1/customers/{id}/messages,
// GET /v1/transactional/{id}/messages, and GET /v1/messages. Open/click/
// bounce state is not part of this object; join against
// GetCustomerActivities (or GetMessage's associations) for delivery
// lifecycle events.
type Message struct {
	ID                  string               `json:"id"`
	DeduplicateID       string               `json:"deduplicate_id,omitempty"`
	CustomerID          string               `json:"customer_id,omitempty"`
	CustomerIdentifiers *CustomerIdentifiers `json:"customer_identifiers,omitempty"`
	CampaignID          int                  `json:"campaign_id,omitempty"`
	ActionID            int                  `json:"action_id,omitempty"`
	ParentActionID      *int                 `json:"parent_action_id,omitempty"`
	TriggerEventID      string               `json:"trigger_event_id,omitempty"`
	Recipient           string               `json:"recipient,omitempty"`
	Subject             string               `json:"subject,omitempty"`
	Metrics             map[string]int64     `json:"metrics,omitempty"`
	Created             int64                `json:"created,omitempty"`
	FailureMessage      string               `json:"failure_message,omitempty"`
	NewsletterID        *int                 `json:"newsletter_id,omitempty"`
	ContentID           *int                 `json:"content_id,omitempty"`
	BroadcastID         *int                 `json:"broadcast_id,omitempty"`
	Type                string               `json:"type,omitempty"`
	Forgotten           bool                 `json:"forgotten,omitempty"`
	// TrackedResponses' shape isn't documented; kept opaque rather than guessed.
	TrackedResponses any `json:"tracked_responses,omitempty"`
}

// MessagesResponse is the decoded shape of GET /v1/customers/{id}/messages.
type MessagesResponse struct {
	Messages []Message `json:"messages"`
	Next     string    `json:"next,omitempty"`
}

// CustomerMessagesOptions filters GetCustomerMessages. IDType lets callers
// look a customer up by email or cio_id directly (id_type=email skips a
// separate SearchCustomers/GetCustomersByEmail round trip).
type CustomerMessagesOptions struct {
	PaginationOptions
	IDType  IdentifierType
	StartTS int64
	EndTS   int64
}

// GetCustomerMessages returns messages sent to one customer, most recent
// first. customerID is interpreted according to opts.IDType (default
// IdentifierTypeID) — pass an email address with IDType: IdentifierTypeEmail
// to look a customer up by email directly.
// See https://docs.customer.io/api/app/#operation/getPersonMessages
func (c *APIClient) GetCustomerMessages(ctx context.Context, customerID string, opts CustomerMessagesOptions) (*MessagesResponse, error) {
	if customerID == "" {
		return nil, ParamError{Param: "customerID"}
	}

	q := opts.PaginationOptions.apply(newQuery()).
		setString("id_type", string(opts.IDType)).
		setInt64("start_ts", opts.StartTS).
		setInt64("end_ts", opts.EndTS)

	requestPath := formatPath("/v1/customers/%s/messages", customerID) + q.String()

	var resp MessagesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ObjectIdentifiers identifies an object by the two identifier kinds
// Customer.io tracks for objects (as opposed to CustomerIdentifiers, which
// is for customers).
type ObjectIdentifiers struct {
	ObjectID    string `json:"object_id,omitempty"`
	CioObjectID string `json:"cio_object_id,omitempty"`
}

// CustomerRelationship is a customer's relationship to an object.
// ObjectTypeID is decoded leniently (json.Number): the schema types it as
// an integer, but the analogous GetObjectRelationships example shows it as
// a numeric string.
type CustomerRelationship struct {
	ObjectTypeID json.Number       `json:"object_type_id,omitempty"`
	Identifiers  ObjectIdentifiers `json:"identifiers"`
	Attributes   map[string]any    `json:"attributes,omitempty"`
	Timestamps   map[string]int64  `json:"timestamps,omitempty"`
}

// RelationshipsResponse is the decoded shape of GET /v1/customers/{id}/relationships.
type RelationshipsResponse struct {
	Relationships []CustomerRelationship `json:"cio_relationships"`
	Next          string                 `json:"next,omitempty"`
}

// GetCustomerRelationships returns the objects a customer is related to.
// See https://docs.customer.io/api/app/#operation/getPersonRelationships
func (c *APIClient) GetCustomerRelationships(ctx context.Context, customerID string, opts PaginationOptions) (*RelationshipsResponse, error) {
	if customerID == "" {
		return nil, ParamError{Param: "customerID"}
	}

	requestPath := formatPath("/v1/customers/%s/relationships", customerID) + opts.apply(newQuery()).String()

	var resp RelationshipsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CustomerSegment is a manual or data-driven segment a customer belongs to.
type CustomerSegment struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CustomerSegmentsResponse is the decoded shape of GET /v1/customers/{id}/segments.
type CustomerSegmentsResponse struct {
	Segments []CustomerSegment `json:"segments"`
}

// GetCustomerSegments returns the segments a customer belongs to. idType
// selects which kind of identifier customerID is; the zero value defaults
// to IdentifierTypeID. See https://docs.customer.io/api/app/#operation/getPersonSegments
func (c *APIClient) GetCustomerSegments(ctx context.Context, customerID string, idType IdentifierType) (*CustomerSegmentsResponse, error) {
	if customerID == "" {
		return nil, ParamError{Param: "customerID"}
	}
	if idType == "" {
		idType = IdentifierTypeID
	}

	requestPath := formatPath("/v1/customers/%s/segments", customerID) + newQuery().setString("id_type", string(idType)).String()

	var resp CustomerSegmentsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SubscriptionTopicPreference is one topic's subscription state for a
// customer. ID is decoded leniently: the documented schema types it as a
// string, but example responses show a JSON integer.
type SubscriptionTopicPreference struct {
	ID          json.Number `json:"id"`
	Subscribed  bool        `json:"subscribed"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
}

// SubscriptionPreferencesHeader is optional branding shown on the hosted
// subscription preferences page.
type SubscriptionPreferencesHeader struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}

// SubscriptionPreferencesCustomer is the customer object returned by
// GetCustomerSubscriptionPreferences.
type SubscriptionPreferencesCustomer struct {
	ID           string                         `json:"id"`
	Identifiers  *CustomerIdentifiers           `json:"identifiers,omitempty"`
	Topics       []SubscriptionTopicPreference  `json:"topics,omitempty"`
	Unsubscribed bool                           `json:"unsubscribed"`
	Header       *SubscriptionPreferencesHeader `json:"header,omitempty"`
}

// SubscriptionPreferencesResponse is the decoded shape of
// GET /v1/customers/{id}/subscription_preferences.
type SubscriptionPreferencesResponse struct {
	Customer SubscriptionPreferencesCustomer `json:"customer"`
}

// CustomerSubscriptionPreferencesOptions configures GetCustomerSubscriptionPreferences.
type CustomerSubscriptionPreferencesOptions struct {
	IDType   IdentifierType
	Language string
}

// GetCustomerSubscriptionPreferences returns a customer's topic subscription
// state. See https://docs.customer.io/api/app/#operation/getPersonSubscriptionPreferences
func (c *APIClient) GetCustomerSubscriptionPreferences(ctx context.Context, customerID string, opts CustomerSubscriptionPreferencesOptions) (*SubscriptionPreferencesResponse, error) {
	if customerID == "" {
		return nil, ParamError{Param: "customerID"}
	}

	q := newQuery().setString("id_type", string(opts.IDType)).setString("language", opts.Language)
	requestPath := formatPath("/v1/customers/%s/subscription_preferences", customerID) + q.String()

	var resp SubscriptionPreferencesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// getCustomersAttributesRequest is the POST /v1/customers/attributes body.
type getCustomersAttributesRequest struct {
	IDs []string `json:"ids"`
}

// CustomerAttributesEntry wraps one customer in GetCustomersAttributesResponse.
type CustomerAttributesEntry struct {
	Customer Customer `json:"customer"`
}

// GetCustomersAttributesResponse is the decoded shape of POST /v1/customers/attributes.
type GetCustomersAttributesResponse struct {
	Customers []CustomerAttributesEntry `json:"customers"`
}

// GetCustomersAttributes returns profile attributes for up to 1000
// customers by id in one call. See https://docs.customer.io/api/app/#operation/getCustomersAttributes
func (c *APIClient) GetCustomersAttributes(ctx context.Context, ids []string) (*GetCustomersAttributesResponse, error) {
	if len(ids) == 0 {
		return nil, ParamError{Param: "ids"}
	}

	var resp GetCustomersAttributesResponse
	if err := c.doJSON(ctx, "POST", "/v1/customers/attributes", getCustomersAttributesRequest{IDs: ids}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SubscriptionCenterTokenResponse is the decoded shape of
// GET /v1/subscription_center/{customerId}/token.
type SubscriptionCenterTokenResponse struct {
	Token string `json:"token"`
	URL   string `json:"url,omitempty"`
}

// GetSubscriptionCenterToken returns a signed token (and hosted URL) a
// customer can use to manage their own subscription preferences.
// See https://docs.customer.io/api/app/#operation/getSubscriptionCenterAccessToken
func (c *APIClient) GetSubscriptionCenterToken(ctx context.Context, customerID string) (*SubscriptionCenterTokenResponse, error) {
	if customerID == "" {
		return nil, ParamError{Param: "customerID"}
	}

	requestPath := formatPath("/v1/subscription_center/%s/token", customerID)

	var resp SubscriptionCenterTokenResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
