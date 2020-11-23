package customerio

import (
	"encoding/json"
	"time"
)

type TransactionalResponse struct {
	Recipient  string    `json:"recipient"`
	DeliveryID string    `json:"delivery_id"`
	QueuedAt   time.Time `json:"queued_at"`
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

type TransactionalError struct {
	Err    string
	Status int
}

func (e *TransactionalError) Error() string {
	return e.Err
}
