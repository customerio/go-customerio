package customerio

import (
	"context"
	"fmt"
	"net/url"
)

// AddPeopleToSegmentCtx adds customers to a manual segment by segment ID.
// See https://docs.customer.io/api/track/#operation/add_customers
func (c *CustomerIO) AddPeopleToSegmentCtx(ctx context.Context, segmentID int, ids []string, opts ...SegmentOption) error {
	return c.segmentMembership(ctx, "add_customers", segmentID, ids, opts...)
}

// AddPeopleToSegment adds customers to a manual segment by segment ID.
func (c *CustomerIO) AddPeopleToSegment(segmentID int, ids []string, opts ...SegmentOption) error {
	return c.AddPeopleToSegmentCtx(context.Background(), segmentID, ids, opts...)
}

// RemovePeopleFromSegmentCtx removes customers from a manual segment by segment ID.
// See https://docs.customer.io/api/track/#operation/remove_customers
func (c *CustomerIO) RemovePeopleFromSegmentCtx(ctx context.Context, segmentID int, ids []string, opts ...SegmentOption) error {
	return c.segmentMembership(ctx, "remove_customers", segmentID, ids, opts...)
}

// RemovePeopleFromSegment removes customers from a manual segment by segment ID.
func (c *CustomerIO) RemovePeopleFromSegment(segmentID int, ids []string, opts ...SegmentOption) error {
	return c.RemovePeopleFromSegmentCtx(context.Background(), segmentID, ids, opts...)
}

func (c *CustomerIO) segmentMembership(ctx context.Context, action string, segmentID int, ids []string, opts ...SegmentOption) error {
	if segmentID <= 0 {
		return ParamError{Param: "segmentID"}
	}
	if len(ids) == 0 {
		return ParamError{Param: "ids"}
	}

	u := fmt.Sprintf("%s/api/v1/segments/%d/%s", c.URL, segmentID, action)

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
