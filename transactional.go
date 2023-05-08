package customerio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type TransactionalType int

const (
	TransactionalTypeEmail = 0
	TransactionalTypePush  = 1
)

var typeToApi = map[TransactionalType]string{
	TransactionalTypeEmail: "email",
	TransactionalTypePush:  "push",
}

var ErrInvalidTransactionalMessageType = errors.New("unknown transactional message type")

func (c *APIClient) sendTransactional(ctx context.Context, typ TransactionalType, req interface{}) (*TransactionalResponse, error) {
	api, ok := typeToApi[typ]
	if !ok {
		return nil, ErrInvalidTransactionalMessageType
	}

	body, statusCode, err := c.doRequest(ctx, "POST", fmt.Sprintf("/v1/send/%s", api), req)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		var meta struct {
			Meta struct {
				Err string `json:"error"`
			} `json:"meta"`
		}
		if err := json.Unmarshal(body, &meta); err != nil {
			return nil, &TransactionalError{
				StatusCode: statusCode,
				Err:        string(body),
			}
		}
		return nil, &TransactionalError{
			StatusCode: statusCode,
			Err:        meta.Meta.Err,
		}
	}

	var resp TransactionalResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// TransactionalResponse  is a response to the send of a transactional message.
type TransactionalResponse struct {
	// DeliveryID is a unique id for the given message.
	DeliveryID string `json:"delivery_id"`
	// QueuedAt is when the message was queued.
	QueuedAt time.Time `json:"queued_at"`
}

func (t *TransactionalResponse) UnmarshalJSON(b []byte) error {
	var r struct {
		DeliveryID string `json:"delivery_id"`
		QueuedAt   int64  `json:"queued_at"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	t.DeliveryID = r.DeliveryID
	t.QueuedAt = time.Unix(r.QueuedAt, 0)
	return nil
}

// TransactionalError is returned if a transactional message fails to send.
type TransactionalError struct {
	// Err is a more specific error message.
	Err string
	// StatusCode is the http status code for the error.
	StatusCode int
}

func (e *TransactionalError) Error() string {
	return e.Err
}
