package customerio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const errUnexpectedStatusCode = "unexpected status code %d"

// SegmentState represents the possible states of a segment.
// Enum values:
//   - events: currently handling event conditions for this segment
//   - build: currently handling profile attribute conditions for this segment
//   - events_queued: waiting to process event conditions for this segment
//   - build_queued: waiting to process profile attribute conditions for this segment
//   - finished: the segment has completed building
type SegmentState string

const (
	SegmentStateEvents       SegmentState = "events"
	SegmentStateBuild        SegmentState = "build"
	SegmentStateEventsQueued SegmentState = "events_queued"
	SegmentStateBuildQueued  SegmentState = "build_queued"
	SegmentStateFinished     SegmentState = "finished"
)

// SegmentType represents the type of a segment.
// Enum values:
//   - dynamic: segment is dynamic, meaning it changes over time
//   - manual: segment is manually maintained
type SegmentType string

const (
	SegmentTypeDynamic SegmentType = "dynamic"
	SegmentTypeManual  SegmentType = "manual"
)

// Segment represents a segment object returned by the API.
type Segment struct {
	ID            int          `json:"id"`                 // Unique identifier for the segment.
	DeduplicateID string       `json:"deduplicate_id"`     // A string in the format id:timestamp.
	Name          string       `json:"name"`               // Name of the segment.
	Description   string       `json:"description"`        // Description of the segment's purpose.
	State         SegmentState `json:"state"`              // Current state of the segment.
	Progress      *int         `json:"progress,omitempty"` // Progress percentage if the segment is building.
	Type          SegmentType  `json:"type"`               // Type of segment: dynamic or manual.
	Tags          []string     `json:"tags,omitempty"`     // Optional tags associated with the segment.
}

// CreateSegmentRequest represents the payload to create a new segment.
type CreateSegmentRequest struct {
	Segment Segment `json:"segment"` // Segment data to create.
}

// CreateSegmentResponse represents the response body for a segment creation request.
type CreateSegmentResponse struct {
	Segment Segment `json:"segment"` // The created segment.
}

// CreateSegment sends a request to create a new segment and returns the created segment data.
func (c *APIClient) CreateSegment(ctx context.Context, req *CreateSegmentRequest) (*CreateSegmentResponse, error) {
	body, statusCode, err := c.doRequest(ctx, "POST", "/v1/segments", req)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf(errUnexpectedStatusCode, statusCode)
	}

	var response CreateSegmentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ListSegmentsResponse represents the response containing multiple segments.
type ListSegmentsResponse struct {
	Segments []Segment `json:"segments"` // List of segments.
}

// ListSegments retrieves all segments from the API.
func (c *APIClient) ListSegments(ctx context.Context) (*ListSegmentsResponse, error) {
	respBody, statusCode, err := c.doRequest(ctx, "GET", "/v1/segments", nil)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf(errUnexpectedStatusCode, statusCode)
	}

	var response ListSegmentsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetSegmentResponse represents the response for retrieving a single segment.
type GetSegmentResponse struct {
	Segment Segment `json:"segment"` // The requested segment.
}

// GetSegment retrieves a specific segment by its ID.
func (c *APIClient) GetSegment(ctx context.Context, segmentID int) (*GetSegmentResponse, error) {
	respBody, statusCode, err := c.doRequest(ctx, "GET", fmt.Sprintf("/v1/segments/%d", segmentID), nil)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf(errUnexpectedStatusCode, statusCode)
	}

	var response GetSegmentResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// DeleteSegment removes a segment by its ID.
func (c *APIClient) DeleteSegment(ctx context.Context, segmentID int) error {
	_, statusCode, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/v1/segments/%d", segmentID), nil)
	if err != nil {
		return err
	}
	if statusCode != http.StatusNoContent {
		return fmt.Errorf(errUnexpectedStatusCode, statusCode)
	}
	return nil
}

// GetSegmentDependenciesResponse represents the response containing segment dependencies.
type GetSegmentDependenciesResponse struct {
	UsedBy struct {
		Campaigns        []int `json:"campaigns"`         // List of campaigns using this segment.
		SentNewsletters  []int `json:"sent_newsletters"`  // List of sent newsletters using this segment.
		DraftNewsletters []int `json:"draft_newsletters"` // List of draft newsletters using this segment.
	} `json:"used_by"` // Dependencies of the segment.
}

// GetSegmentDependencies returns the dependencies of a specific segment.
func (c *APIClient) GetSegmentDependencies(ctx context.Context, segmentID int) (*GetSegmentDependenciesResponse, error) {
	respBody, statusCode, err := c.doRequest(ctx, "GET", fmt.Sprintf("/v1/segments/%d/used_by", segmentID), nil)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf(errUnexpectedStatusCode, statusCode)
	}

	var response GetSegmentDependenciesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetSegmentCustomerCountResponse represents the response for retrieving customer count in a segment.
type GetSegmentCustomerCountResponse struct {
	Count int `json:"count"` // Number of customers in the segment.
}

// GetSegmentCustomerCount returns the total number of customers in a specific segment.
func (c *APIClient) GetSegmentCustomerCount(ctx context.Context, segmentID int) (*GetSegmentCustomerCountResponse, error) {
	respBody, statusCode, err := c.doRequest(ctx, "GET", fmt.Sprintf("/v1/segments/%d/customer_count", segmentID), nil)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf(errUnexpectedStatusCode, statusCode)
	}

	var response GetSegmentCustomerCountResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ListCustomersInSegmentResponse represents the response for listing customers in a segment.
type ListCustomersInSegmentResponse struct {
	IDs         []string             `json:"ids"`         // List of customer IDs.
	Identifiers []CustomerIdentifier `json:"identifiers"` // List of customer identifiers.
	Next        string               `json:"next"`        // Optional pagination cursor.
}

// CustomerIdentifier represents the customer identifiers in a segment.
type CustomerIdentifier struct {
	Email string `json:"email"`  // Email of the customer.
	ID    int    `json:"id"`     // Internal ID of the customer.
	CioID string `json:"cio_id"` // Customer.io ID of the customer.
}

// ListCustomersInSegment retrieves a list of customers in a specific segment.
func (c *APIClient) ListCustomersInSegment(ctx context.Context, segmentID int) (*ListCustomersInSegmentResponse, error) {
	respBody, statusCode, err := c.doRequest(ctx, "GET", fmt.Sprintf("/v1/segments/%d/membership", segmentID), nil)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf(errUnexpectedStatusCode, statusCode)
	}

	var response ListCustomersInSegmentResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
