package customerio

import (
	"context"
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
