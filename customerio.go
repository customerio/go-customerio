package customerio

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const DefaultUserAgent = "Customer.io Go Client/" + Version

// CustomerIO wraps the customer.io track API, see: https://customer.io/docs/api/#apitrackintroduction
type CustomerIO struct {
	siteID    string
	apiKey    string
	URL       string
	UserAgent string
	Client    HTTPClient
}

// CustomerIOError is returned by any method that fails at the API level
type CustomerIOError struct {
	status int
	url    string
	body   []byte
}

func (e *CustomerIOError) Error() string {
	return fmt.Sprintf("%v: %v %v", e.status, e.url, string(e.body))
}

// StatusCode returns the HTTP status code from a failed API response.
func (e *CustomerIOError) StatusCode() int {
	return e.status
}

// URL returns the request URL from a failed API response.
func (e *CustomerIOError) URL() string {
	return e.url
}

// Body returns a copy of the failed API response body.
func (e *CustomerIOError) Body() []byte {
	body := make([]byte, len(e.body))
	copy(body, e.body)
	return body
}

// ParamError is an error returned if a parameter to the track API is invalid.
type ParamError struct {
	Param string // Param is the name of the parameter.
}

func (e ParamError) Error() string { return e.Param + ": missing" }

// NewTrackClient prepares a client for use with the Customer.io track API, see: https://customer.io/docs/api/#apitrackintroduction
// using a Tracking Site ID and API Key pair from https://fly.customer.io/settings/api_credentials
func NewTrackClient(siteID, apiKey string, opts ...Option) *CustomerIO {
	c := &CustomerIO{
		siteID:    siteID,
		apiKey:    apiKey,
		URL:       "https://track.customer.io",
		UserAgent: DefaultUserAgent,
		Client:    newDefaultHTTPClient(),
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyTrack(c)
	}

	return c
}

// NewCustomerIO prepares a client for use with the Customer.io track API, see: https://customer.io/docs/api/#apitrackintroduction
//
// Deprecated: use NewTrackClient.
func NewCustomerIO(siteID, apiKey string) *CustomerIO {
	return NewTrackClient(siteID, apiKey)
}

// IdentifyCtx identifies a customer and sets their attributes
func (c *CustomerIO) IdentifyCtx(ctx context.Context, customerID string, attributes map[string]any) error {
	if customerID == "" {
		return ParamError{Param: "customerID"}
	}
	return c.request(ctx, "PUT",
		fmt.Sprintf("%s/api/v1/customers/%s", c.URL, url.PathEscape(customerID)),
		attributes)
}

// Identify identifies a customer and sets their attributes
func (c *CustomerIO) Identify(customerID string, attributes map[string]any) error {
	return c.IdentifyCtx(context.Background(), customerID, attributes)
}

// TrackCtx sends a single event to Customer.io for the supplied user
func (c *CustomerIO) TrackCtx(ctx context.Context, customerID string, eventName string, data map[string]any, opts ...TrackOption) error {
	if customerID == "" {
		return ParamError{Param: "customerID"}
	}
	if eventName == "" {
		return ParamError{Param: "eventName"}
	}
	return c.request(ctx, "POST",
		fmt.Sprintf("%s/api/v1/customers/%s/events", c.URL, url.PathEscape(customerID)),
		trackPayload(eventName, data, opts...))
}

// Track sends a single event to Customer.io for the supplied user
func (c *CustomerIO) Track(customerID string, eventName string, data map[string]any, opts ...TrackOption) error {
	return c.TrackCtx(context.Background(), customerID, eventName, data, opts...)
}

// TrackAnonymousCtx sends a single event to Customer.io for the anonymous user
func (c *CustomerIO) TrackAnonymousCtx(ctx context.Context, anonymousID, eventName string, data map[string]any, opts ...TrackOption) error {
	if eventName == "" {
		return ParamError{Param: "eventName"}
	}

	payload := trackPayload(eventName, data, opts...)

	if anonymousID != "" {
		payload["anonymous_id"] = anonymousID
	}

	return c.request(ctx, "POST", fmt.Sprintf("%s/api/v1/events", c.URL), payload)
}

// TrackAnonymous sends a single event to Customer.io for the anonymous user
func (c *CustomerIO) TrackAnonymous(anonymousID, eventName string, data map[string]any, opts ...TrackOption) error {
	return c.TrackAnonymousCtx(context.Background(), anonymousID, eventName, data, opts...)
}

// DeleteCtx deletes a customer
func (c *CustomerIO) DeleteCtx(ctx context.Context, customerID string) error {
	if customerID == "" {
		return ParamError{Param: "customerID"}
	}
	return c.request(ctx, "DELETE",
		fmt.Sprintf("%s/api/v1/customers/%s", c.URL, url.PathEscape(customerID)),
		nil)
}

func (c *CustomerIO) auth() string {
	return base64.URLEncoding.EncodeToString(fmt.Appendf(nil, "%v:%v", c.siteID, c.apiKey))
}

func (c *CustomerIO) request(ctx context.Context, method, url string, body any) error {
	respBody, statusCode, err := doHTTP(ctx, c.Client, method, url, c.UserAgent, body, func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %v", c.auth()))
	})
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return &CustomerIOError{
			status: statusCode,
			url:    url,
			body:   respBody,
		}
	}

	return nil
}

type IdentifierType string

const (
	IdentifierTypeID    IdentifierType = "id"
	IdentifierTypeEmail IdentifierType = "email"
	IdentifierTypeCioID IdentifierType = "cio_id"
)

type Identifier struct {
	Type  IdentifierType
	Value string
}

func (id Identifier) kv() map[string]string {
	return map[string]string{
		string(id.Type): id.Value,
	}
}

func (id Identifier) validate() error {
	if id.Type != IdentifierTypeID &&
		id.Type != IdentifierTypeEmail &&
		id.Type != IdentifierTypeCioID {
		return errors.New("invalid id type")
	}

	if strings.TrimSpace(id.Value) == "" {
		return errors.New("invalid id")
	}

	return nil
}

// MergeCustomersCtx sends a request to Customer.io to merge two customer profiles together.
func (c *CustomerIO) MergeCustomersCtx(ctx context.Context, primary Identifier, secondary Identifier) error {
	if err := primary.validate(); err != nil {
		return fmt.Errorf("primary: %w", err)
	}
	if err := secondary.validate(); err != nil {
		return fmt.Errorf("secondary: %w", err)
	}

	return c.request(ctx, "POST",
		fmt.Sprintf("%s/api/v1/merge_customers", c.URL),
		map[string]any{
			"primary":   primary.kv(),
			"secondary": secondary.kv(),
		})
}

// MergeCustomers sends a request to Customer.io to merge two customer profiles together.
func (c *CustomerIO) MergeCustomers(primary Identifier, secondary Identifier) error {
	return c.MergeCustomersCtx(context.Background(), primary, secondary)
}
