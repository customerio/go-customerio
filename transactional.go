package customerio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TransactionalAPIClient struct {
	Key    string
	Host   string
	Client *http.Client
}

func NewTransactionalClient(key string) *TransactionalAPIClient {
	return &TransactionalAPIClient{
		Key:    key,
		Client: http.DefaultClient,
		Host:   "api.customer.io",
	}
}

type TransactionalResponse struct {
	Recipient  string `json:"recipient"`
	DeliveryID string `json:"delivery_id"`
	QueuedAt   int    `json:"queued_at"`
}

type TransactionalError struct {
	Meta struct {
		Err string `json:"error"`
	} `json:"meta"`
	Status int `json:"-"`
}

func (e *TransactionalError) Error() string {
	return e.Meta.Err
}

func (c *TransactionalAPIClient) SendEmail(e Email) (*TransactionalResponse, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/v1/send/email", c.Host), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Key))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 {
		var e TransactionalError
		if err := d.Decode(&e); err != nil {
			return nil, err
		}
		return nil, &e
	}

	var result TransactionalResponse
	if err := d.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
