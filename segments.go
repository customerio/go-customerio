package customerio

import (
	"context"
	"net/url"
)

// AddPeopleToSegment adds customers to a manual segment by segment ID.
// See https://docs.customer.io/api/track/#operation/add_customers
func (c *CustomerIO) AddPeopleToSegment(ctx context.Context, segmentID int, ids []string, opts ...SegmentOption) error {
	return c.segmentMembership(ctx, "add_customers", segmentID, ids, opts...)
}

// RemovePeopleFromSegment removes customers from a manual segment by segment ID.
// See https://docs.customer.io/api/track/#operation/remove_customers
func (c *CustomerIO) RemovePeopleFromSegment(ctx context.Context, segmentID int, ids []string, opts ...SegmentOption) error {
	return c.segmentMembership(ctx, "remove_customers", segmentID, ids, opts...)
}

func (c *CustomerIO) segmentMembership(ctx context.Context, action string, segmentID int, ids []string, opts ...SegmentOption) error {
	if segmentID <= 0 {
		return ParamError{Param: "segmentID"}
	}
	if len(ids) == 0 {
		return ParamError{Param: "ids"}
	}

	u := c.URL + formatPath("/api/v1/segments/%d/%s", segmentID, action)

	// url.Values handles encoding for any future option that carries an
	// arbitrary string; current options (id_type) only emit URL-safe values.
	q := url.Values{}
	for _, opt := range opts {
		opt(q)
	}
	if encoded := q.Encode(); encoded != "" {
		u += "?" + encoded
	}

	return c.request(ctx, "POST", u, map[string]interface{}{
		"ids": ids,
	})
}
