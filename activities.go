package customerio

import "context"

// Activity is an entry from GET /activities or GET /customers/{id}/activities.
// Data's shape varies by Type: message-delivery activities (sent_email,
// delivered_email, opened_email, ...) carry {delivery_id, delivered, opened};
// attribute_change activities carry a map of attribute name to {from, to};
// event/page/screen activities have no documented shape. It is left as
// map[string]any rather than modeled per Type for that reason.
type Activity struct {
	ID                  string               `json:"id"`
	Type                string               `json:"type"`
	Timestamp           int64                `json:"timestamp"`
	CustomerID          string               `json:"customer_id,omitempty"`
	CustomerIdentifiers *CustomerIdentifiers `json:"customer_identifiers,omitempty"`
	DeliveryID          string               `json:"delivery_id,omitempty"`
	DeliveryType        string               `json:"delivery_type,omitempty"`
	Name                string               `json:"name,omitempty"`
	URL                 string               `json:"url,omitempty"`
	Data                map[string]any       `json:"data,omitempty"`
}

// ActivitiesResponse is the decoded shape of GET /activities and
// GET /customers/{id}/activities.
type ActivitiesResponse struct {
	Activities []Activity `json:"activities"`
	Next       string     `json:"next,omitempty"`
}

// ListActivitiesOptions filters ListActivities.
type ListActivitiesOptions struct {
	PaginationOptions
	Type       string
	Name       string
	Deleted    *bool
	CustomerID string
	IDType     IdentifierType
}

// ListActivities returns activity records across all customers, most recent
// first. See https://docs.customer.io/api/app/#operation/listActivities
func (c *APIClient) ListActivities(ctx context.Context, opts ListActivitiesOptions) (*ActivitiesResponse, error) {
	q := opts.PaginationOptions.apply(newQuery()).
		setString("type", opts.Type).
		setString("name", opts.Name).
		setBool("deleted", opts.Deleted).
		setString("customer_id", opts.CustomerID).
		setString("id_type", string(opts.IDType))

	requestPath := "/v1/activities" + q.String()

	var resp ActivitiesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
