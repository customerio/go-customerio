package customerio

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type APIClient struct {
	// Deprecated: Use NewAPIClient constructor options instead. Will be unexported in v4.
	Key string
	// Deprecated: Use NewAPIClient with WithURL or WithRegion instead. Will be unexported in v4.
	URL string
	// Deprecated: Use NewAPIClient with WithUserAgent instead. Will be unexported in v4.
	UserAgent string
	// Deprecated: Use NewAPIClient with WithHTTPClient instead. Will be unexported in v4.
	Client HTTPClient
}

// NewAPIClient prepares a client for use with the Customer.io API, see: https://customer.io/docs/api/#apicoreintroduction
// using an App API Key from https://fly.customer.io/settings/api_credentials?keyType=app
func NewAPIClient(key string, opts ...Option) *APIClient {
	client := &APIClient{
		Key:       key,
		Client:    newDefaultHTTPClient(),
		URL:       "https://api.customer.io",
		UserAgent: DefaultUserAgent,
	}

	for _, opt := range opts {
		if opt != nil {
			opt.applyAPI(client)
		}
	}
	return client
}

func (c *APIClient) doRequest(ctx context.Context, verb, requestPath string, body any) ([]byte, int, error) {
	return doHTTP(ctx, c.Client, verb, c.URL+requestPath, c.UserAgent, body, func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+c.Key)
	})
}

// doJSON executes verb against requestPath, marshaling body as the JSON
// request payload (nil for none) and unmarshaling the response into out
// (nil to discard it). A response whose status isn't in okStatuses
// (defaulting to just http.StatusOK) is returned as a *CustomerIOError.
func (c *APIClient) doJSON(ctx context.Context, verb, requestPath string, body, out any, okStatuses ...int) error {
	if len(okStatuses) == 0 {
		okStatuses = []int{http.StatusOK}
	}

	respBody, statusCode, err := c.doRequest(ctx, verb, requestPath, body)
	if err != nil {
		return err
	}

	ok := false
	for _, s := range okStatuses {
		if statusCode == s {
			ok = true
			break
		}
	}
	if !ok {
		return &CustomerIOError{status: statusCode, url: c.URL + requestPath, body: respBody}
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return err
		}
	}
	return nil
}

// doMultipartJSON is doJSON's counterpart for multipart/form-data uploads —
// currently only CreateAsset. body/contentType come from a
// *multipart.Writer (contentType must include its boundary, from
// FormDataContentType()).
func (c *APIClient) doMultipartJSON(ctx context.Context, verb, requestPath, contentType string, body io.Reader, out any) error {
	respBody, statusCode, err := doMultipartHTTP(ctx, c.Client, verb, c.URL+requestPath, c.UserAgent, contentType, body, func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+c.Key)
	})
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return &CustomerIOError{status: statusCode, url: c.URL + requestPath, body: respBody}
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return err
		}
	}
	return nil
}
