package customerio

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type APIClient struct {
	Key       string
	URL       string
	UserAgent string
	Client    HTTPClient
}

// NewAPIClient prepares a client for use with the Customer.io API, see: https://customer.io/docs/api/#apicoreintroduction
// using an App API Key from https://fly.customer.io/settings/api_credentials?keyType=app
func NewAPIClient(key string, opts ...Option) *APIClient {
	client := &APIClient{
		Key:       key,
		Client:    http.DefaultClient,
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
	b, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, verb, c.URL+requestPath, bytes.NewBuffer(b))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", c.UserAgent)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}
