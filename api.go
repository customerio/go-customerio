package customerio

import "net/http"

type APIClient struct {
	Key    string
	Host   string
	Client *http.Client
}

func NewAPIClient(key string) *APIClient {
	return &APIClient{
		Key:    key,
		Client: http.DefaultClient,
		Host:   "api.customer.io",
	}
}
