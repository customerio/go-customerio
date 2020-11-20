package customerio

import "net/http"

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
