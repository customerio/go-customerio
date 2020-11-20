package customerio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TransactionalResponse struct {
	Recipient  string `json:"recipient"`
	DeliveryID string `json:"delivery_id"`
	QueuedAt   int    `json:"queued_at"`
}

type TransactionalError struct {
	Err    string
	Status int
}

func (e *TransactionalError) Error() string {
	return e.Err
}

func (c *APIClient) SendEmail(e Email) (*TransactionalResponse, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/send/email", c.URL), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Key))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var meta struct {
			Meta struct {
				Err string `json:"error"`
			} `json:"meta"`
		}
		if err := json.Unmarshal(body, &meta); err != nil {
			return nil, &TransactionalError{
				Status: resp.StatusCode,
				Err:    string(body),
			}
		}
		return nil, &TransactionalError{
			Status: resp.StatusCode,
			Err:    meta.Meta.Err,
		}
	}

	var result TransactionalResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
