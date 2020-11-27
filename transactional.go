package customerio

import (
	"encoding/json"
	"time"
)

// TransactionalResponse  is a response to the send of a transactional message.
type TransactionalResponse struct {
	// Recipient is the recipient of the message.
	Recipient string `json:"recipient"`
	// DeliveryID is a unique id for the given message.
	DeliveryID string `json:"delivery_id"`
	// QueuedAt is when the message was queued.
	QueuedAt time.Time `json:"queued_at"`
}

func (t *TransactionalResponse) UnmarshalJSON(b []byte) error {
	var r struct {
		Recipient  string `json:"recipient"`
		DeliveryID string `json:"delivery_id"`
		QueuedAt   int64  `json:"queued_at"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	t.Recipient = r.Recipient
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
