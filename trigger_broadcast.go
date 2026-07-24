package customerio

import (
	"context"
	"encoding/json"
	"net/http"
)

// BroadcastRecipients defines who receives a broadcast trigger.
// Set Segment for segment-based targeting, or set exactly one of
// Ids, Emails, PerUserData, or DataFileURL for direct targeting.
type BroadcastRecipients struct {
	Segment     map[string]any   `json:"segment,omitempty"`
	Ids         []string         `json:"ids,omitempty"`
	Emails      []string         `json:"emails,omitempty"`
	PerUserData []map[string]any `json:"per_user_data,omitempty"`
	DataFileURL string           `json:"data_file_url,omitempty"`
}

// BroadcastOptions controls how the API processes direct recipients.
// Each flag only applies to certain recipient types; flags that don't
// apply to the chosen recipient type are dropped before sending, since
// extra fields cause server-side rejection.
type BroadcastOptions struct {
	IDIgnoreMissing    *bool
	EmailIgnoreMissing *bool
	EmailAddDuplicates *bool
}

// BroadcastResponse is returned when a broadcast is successfully triggered.
type BroadcastResponse struct {
	ID int `json:"id"`
}

// broadcastInput bundles the inputs to buildBroadcastPayload.
type broadcastInput struct {
	Data       map[string]any
	Recipients BroadcastRecipients
	Options    BroadcastOptions
}

// broadcastPayload is the wire shape for /v1/campaigns/{id}/triggers.
// Direct recipient fields (Ids/Emails/PerUserData/DataFileURL) sit at the top
// level; for segment-based targeting, the full BroadcastRecipients goes under
// Recipients. omitempty enforces mutual exclusion at marshal time.
type broadcastPayload struct {
	Data               map[string]any       `json:"data"`
	Ids                []string             `json:"ids,omitempty"`
	Emails             []string             `json:"emails,omitempty"`
	PerUserData        []map[string]any     `json:"per_user_data,omitempty"`
	DataFileURL        string               `json:"data_file_url,omitempty"`
	Recipients         *BroadcastRecipients `json:"recipients,omitempty"`
	IDIgnoreMissing    *bool                `json:"id_ignore_missing,omitempty"`
	EmailIgnoreMissing *bool                `json:"email_ignore_missing,omitempty"`
	EmailAddDuplicates *bool                `json:"email_add_duplicates,omitempty"`
}

// TriggerBroadcast triggers a broadcast by POSTing to /v1/campaigns/{id}/triggers.
// For segment-based targeting, set recipients.Segment. For direct targeting, set exactly one
// of recipients.Ids, recipients.Emails, recipients.PerUserData, or recipients.DataFileURL.
// opts.IDIgnoreMissing/EmailIgnoreMissing/EmailAddDuplicates apply only to direct
// targeting and are filtered to the recipient type in use.
func (c *APIClient) TriggerBroadcast(ctx context.Context, broadcastID int, data map[string]any, recipients BroadcastRecipients, opts BroadcastOptions) (*BroadcastResponse, error) {
	if broadcastID <= 0 {
		return nil, ParamError{Param: "broadcastID"}
	}

	payload := buildBroadcastPayload(broadcastInput{Data: data, Recipients: recipients, Options: opts})

	requestPath := formatPath("/v1/campaigns/%d/triggers", broadcastID)
	body, statusCode, err := c.doRequest(ctx, "POST", requestPath, payload)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, &CustomerIOError{
			status: statusCode,
			url:    c.URL + requestPath,
			body:   body,
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
func buildBroadcastPayload(in broadcastInput) broadcastPayload {
	r := in.Recipients
	o := in.Options
	p := broadcastPayload{Data: in.Data}

	switch {
	case len(r.Ids) > 0:
		p.Ids = r.Ids
		p.IDIgnoreMissing = o.IDIgnoreMissing
	case len(r.Emails) > 0:
		p.Emails = r.Emails
		p.EmailIgnoreMissing = o.EmailIgnoreMissing
		p.EmailAddDuplicates = o.EmailAddDuplicates
	case len(r.PerUserData) > 0:
		p.PerUserData = r.PerUserData
		p.IDIgnoreMissing = o.IDIgnoreMissing
		p.EmailIgnoreMissing = o.EmailIgnoreMissing
		p.EmailAddDuplicates = o.EmailAddDuplicates
	case r.DataFileURL != "":
		p.DataFileURL = r.DataFileURL
		p.IDIgnoreMissing = o.IDIgnoreMissing
		p.EmailIgnoreMissing = o.EmailIgnoreMissing
		p.EmailAddDuplicates = o.EmailAddDuplicates
	default:
		p.Recipients = &r
	}

	return p
}

// BroadcastTriggerStatus reports one broadcast trigger's processing state.
type BroadcastTriggerStatus struct {
	ID          int   `json:"id"`
	BroadcastID int   `json:"broadcast_id,omitempty"`
	CreatedAt   int64 `json:"created_at,omitempty"`
	ProcessedAt int64 `json:"processed_at,omitempty"`
}

// GetBroadcastTriggerStatus returns one broadcast trigger's processing
// state. Despite the "broadcast" naming, this addresses
// /v1/campaigns/{broadcastID}/triggers/{triggerID} — the App API files
// broadcast triggers under the campaigns resource; GetBroadcastTriggers
// (broadcasts.go) is the analogous list endpoint under /v1/broadcasts/{id}/triggers.
// See https://docs.customer.io/api/app/#operation/getBroadcastTriggerStatus
func (c *APIClient) GetBroadcastTriggerStatus(ctx context.Context, broadcastID, triggerID string) (*BroadcastTriggerStatus, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}
	if triggerID == "" {
		return nil, ParamError{Param: "triggerID"}
	}

	requestPath := formatPath("/v1/campaigns/%s/triggers/%s", broadcastID, triggerID)

	var resp BroadcastTriggerStatus
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BroadcastTriggerErrorsResponse is the decoded shape of
// GET /v1/campaigns/{broadcastID}/triggers/{triggerID}/errors. Unlike most
// other paginated App API responses, Next here is a numeric offset, not an
// opaque string cursor.
type BroadcastTriggerErrorsResponse struct {
	Errors []string `json:"errors"`
	Next   int      `json:"next,omitempty"`
}

// GetBroadcastTriggerErrors returns the per-recipient errors from one
// broadcast trigger (e.g. "line 4: invalid email").
// See https://docs.customer.io/api/app/#operation/getBroadcastTriggerErrors
func (c *APIClient) GetBroadcastTriggerErrors(ctx context.Context, broadcastID, triggerID string, opts PaginationOptions) (*BroadcastTriggerErrorsResponse, error) {
	if broadcastID == "" {
		return nil, ParamError{Param: "broadcastID"}
	}
	if triggerID == "" {
		return nil, ParamError{Param: "triggerID"}
	}

	requestPath := formatPath("/v1/campaigns/%s/triggers/%s/errors", broadcastID, triggerID) + opts.apply(newQuery()).String()

	var resp BroadcastTriggerErrorsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
