package customerio

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type APIClient struct {
	Key    string
	URL    string
	Client *http.Client
}

func NewAPIClient(key string) *APIClient {
	return &APIClient{
		Key:    key,
		Client: http.DefaultClient,
		URL:    "https://api.customer.io",
	}
}

func (c *APIClient) doRequest(ctx context.Context, verb, requestPath string, body interface{}) (*http.Response, []byte, error) {

	b, err := json.Marshal(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", c.URL+requestPath, bytes.NewBuffer(b))
	if err != nil {
		return nil, nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Authorization", "Bearer "+c.Key)
	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	var responseBody []byte

	errc := make(chan error, 1)
	go func() {
		resp, err = c.Client.Do(req)
		if err != nil {
			errc <- err
			return
		}
		defer resp.Body.Close()

		responseBody, err = ioutil.ReadAll(resp.Body)
		errc <- err
	}()

	select {
	case <-ctx.Done():
		<-errc // Wait for f to return.
		return resp, responseBody, ctx.Err()
	case err := <-errc:
		return resp, responseBody, err
	}
}
