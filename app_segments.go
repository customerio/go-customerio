package customerio

import (
	"context"
	"net/http"
)

// Segment is a manual or data-driven audience segment.
type Segment struct {
	ID            int      `json:"id"`
	DeduplicateID string   `json:"deduplicate_id,omitempty"`
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	State         string   `json:"state,omitempty"`
	Progress      *int     `json:"progress,omitempty"`
	Type          string   `json:"type,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	CreatedAt     int64    `json:"created_at,omitempty"`
	UpdatedAt     int64    `json:"updated_at,omitempty"`
}

// SegmentResponse wraps a single Segment, as returned by CreateSegment and GetSegment.
type SegmentResponse struct {
	Segment Segment `json:"segment"`
}

// ListSegmentsResponse is the decoded shape of GET /v1/segments.
type ListSegmentsResponse struct {
	Segments []Segment `json:"segments"`
}

// ListSegments returns every segment defined in the workspace.
// See https://docs.customer.io/api/app/#operation/listSegments
func (c *APIClient) ListSegments(ctx context.Context) (*ListSegmentsResponse, error) {
	var resp ListSegmentsResponse
	if err := c.doJSON(ctx, "GET", "/v1/segments", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SegmentInput describes a new manual segment.
type SegmentInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type createSegmentRequest struct {
	Segment SegmentInput `json:"segment"`
}

// CreateSegment creates a new manual segment.
// See https://docs.customer.io/api/app/#operation/createSegment
func (c *APIClient) CreateSegment(ctx context.Context, input SegmentInput) (*SegmentResponse, error) {
	if input.Name == "" {
		return nil, ParamError{Param: "name"}
	}

	var resp SegmentResponse
	if err := c.doJSON(ctx, "POST", "/v1/segments", createSegmentRequest{Segment: input}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSegment returns one segment by id.
// See https://docs.customer.io/api/app/#operation/getSegment
func (c *APIClient) GetSegment(ctx context.Context, segmentID int) (*SegmentResponse, error) {
	if segmentID <= 0 {
		return nil, ParamError{Param: "segmentID"}
	}

	var resp SegmentResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/segments/%d", segmentID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteSegment deletes a manual segment by id.
// See https://docs.customer.io/api/app/#operation/deleteSegment
func (c *APIClient) DeleteSegment(ctx context.Context, segmentID int) error {
	if segmentID <= 0 {
		return ParamError{Param: "segmentID"}
	}

	return c.doJSON(ctx, "DELETE", formatPath("/v1/segments/%d", segmentID), nil, nil, http.StatusOK, http.StatusNoContent)
}

// SegmentCustomerCountResponse is the decoded shape of
// GET /v1/segments/{id}/customer_count.
type SegmentCustomerCountResponse struct {
	SegmentID int `json:"segment_id"`
	Count     int `json:"count"`
}

// GetSegmentCustomerCount returns how many customers currently belong to a segment.
// See https://docs.customer.io/api/app/#operation/getSegmentCustomerCount
func (c *APIClient) GetSegmentCustomerCount(ctx context.Context, segmentID int) (*SegmentCustomerCountResponse, error) {
	if segmentID <= 0 {
		return nil, ParamError{Param: "segmentID"}
	}

	var resp SegmentCustomerCountResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/segments/%d/customer_count", segmentID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SegmentMembershipResponse is the decoded shape of
// GET /v1/segments/{id}/membership. Which of IDs/Identifiers is populated
// (and which CustomerIdentifiers fields are set) depends on the workspace's
// identifier configuration.
type SegmentMembershipResponse struct {
	SegmentID   int                   `json:"segment_id"`
	IDs         []string              `json:"ids,omitempty"`
	Identifiers []CustomerIdentifiers `json:"identifiers,omitempty"`
	Next        string                `json:"next,omitempty"`
}

// GetSegmentMembership returns the customers belonging to a segment.
// See https://docs.customer.io/api/app/#operation/getSegmentMembership
func (c *APIClient) GetSegmentMembership(ctx context.Context, segmentID int, opts PaginationOptions) (*SegmentMembershipResponse, error) {
	if segmentID <= 0 {
		return nil, ParamError{Param: "segmentID"}
	}

	requestPath := formatPath("/v1/segments/%d/membership", segmentID) + opts.apply(newQuery()).String()

	var resp SegmentMembershipResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SegmentUsedBy lists the campaigns/newsletters referencing a segment.
type SegmentUsedBy struct {
	Campaigns        []int `json:"campaigns,omitempty"`
	SentNewsletters  []int `json:"sent_newsletters,omitempty"`
	DraftNewsletters []int `json:"draft_newsletters,omitempty"`
}

// SegmentUsedByResponse is the decoded shape of GET /v1/segments/{id}/used_by.
type SegmentUsedByResponse struct {
	SegmentID int           `json:"segment_id"`
	UsedBy    SegmentUsedBy `json:"used_by"`
}

// GetSegmentUsedBy returns what campaigns/newsletters reference a segment —
// useful to check before deleting it.
// See https://docs.customer.io/api/app/#operation/getSegmentUsedBy
func (c *APIClient) GetSegmentUsedBy(ctx context.Context, segmentID int) (*SegmentUsedByResponse, error) {
	if segmentID <= 0 {
		return nil, ParamError{Param: "segmentID"}
	}

	var resp SegmentUsedByResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/segments/%d/used_by", segmentID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
