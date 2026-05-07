package customerio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// BroadcastRecipients defines who receives a broadcast trigger.
// Set Segment for segment-based targeting, or set exactly one of
// Ids, Emails, PerUserData, or DataFileURL for direct targeting.
type BroadcastRecipients struct {
	Segment     map[string]interface{}   `json:"segment,omitempty"`
	Ids         []string                 `json:"ids,omitempty"`
	Emails      []string                 `json:"emails,omitempty"`
	PerUserData []map[string]interface{} `json:"per_user_data,omitempty"`
	DataFileURL string                   `json:"data_file_url,omitempty"`

	IDIgnoreMissing    *bool `json:"id_ignore_missing,omitempty"`
	EmailIgnoreMissing *bool `json:"email_ignore_missing,omitempty"`
	EmailAddDuplicates *bool `json:"email_add_duplicates,omitempty"`
}

// BroadcastResponse is returned when a broadcast is successfully triggered.
type BroadcastResponse struct {
	ID int `json:"id"`
}

// TriggerBroadcast triggers a broadcast by POSTing to /v1/campaigns/{id}/triggers.
// For segment-based targeting, set recipients.Segment. For direct targeting, set exactly one
// of recipients.Ids, recipients.Emails, recipients.PerUserData, or recipients.DataFileURL.
func (c *APIClient) TriggerBroadcast(ctx context.Context, broadcastID int, data map[string]interface{}, recipients BroadcastRecipients) (*BroadcastResponse, error) {
	if broadcastID <= 0 {
		return nil, ParamError{Param: "broadcastID"}
	}

	payload := buildBroadcastPayload(data, recipients)

	requestPath := fmt.Sprintf("/v1/campaigns/%d/triggers", broadcastID)
	body, statusCode, err := c.doRequest(ctx, "POST", requestPath, payload)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, &CustomerIOError{
			StatusCode: statusCode,
			URL:        c.URL + requestPath,
			Body:       body,
		}
	}

	var resp BroadcastResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// If a direct recipient field (ids, emails, per_user_data, data_file_url) is present,
// it and its allowed companion options are flattened to the top level alongside data.
// Otherwise the full recipients struct is nested under a "recipients" key.
func buildBroadcastPayload(data map[string]interface{}, r BroadcastRecipients) map[string]interface{} {
	if len(r.Ids) > 0 {
		p := map[string]interface{}{"data": data, "ids": r.Ids}
		if r.IDIgnoreMissing != nil {
			p["id_ignore_missing"] = *r.IDIgnoreMissing
		}
		return p
	}

	if len(r.Emails) > 0 {
		p := map[string]interface{}{"data": data, "emails": r.Emails}
		if r.EmailIgnoreMissing != nil {
			p["email_ignore_missing"] = *r.EmailIgnoreMissing
		}
		if r.EmailAddDuplicates != nil {
			p["email_add_duplicates"] = *r.EmailAddDuplicates
		}
		return p
	}

	if len(r.PerUserData) > 0 {
		p := map[string]interface{}{"data": data, "per_user_data": r.PerUserData}
		if r.IDIgnoreMissing != nil {
			p["id_ignore_missing"] = *r.IDIgnoreMissing
		}
		if r.EmailIgnoreMissing != nil {
			p["email_ignore_missing"] = *r.EmailIgnoreMissing
		}
		if r.EmailAddDuplicates != nil {
			p["email_add_duplicates"] = *r.EmailAddDuplicates
		}
		return p
	}

	if r.DataFileURL != "" {
		p := map[string]interface{}{"data": data, "data_file_url": r.DataFileURL}
		if r.IDIgnoreMissing != nil {
			p["id_ignore_missing"] = *r.IDIgnoreMissing
		}
		if r.EmailIgnoreMissing != nil {
			p["email_ignore_missing"] = *r.EmailIgnoreMissing
		}
		if r.EmailAddDuplicates != nil {
			p["email_add_duplicates"] = *r.EmailAddDuplicates
		}
		return p
	}

	return map[string]interface{}{"data": data, "recipients": r}
}
