package customerio

import (
	"context"
	"fmt"
)

// AddPeopleToSegment adds people to a segment.
func (c *CustomerIO) AddPeopleToSegment(ctx context.Context, segmentID int, ids []string) error {
	if segmentID == 0 {
		return ParamError{Param: "segmentID"}
	}
	if len(ids) == 0 {
		return ParamError{Param: "ids"}
	}
	return c.request(ctx, "POST",
		fmt.Sprintf("%s/api/v1/segments/%d/add_customers", c.URL, segmentID),
		map[string]interface{}{
			"ids": ids,
		})
}

// RemovePeopleFromSegment removes people from a segment
func (c *CustomerIO) RemovePeopleFromSegment(ctx context.Context, segmentID int, ids []string) error {
	if segmentID == 0 {
		return ParamError{Param: "segmentID"}
	}
	if len(ids) == 0 {
		return ParamError{Param: "ids"}
	}
	return c.request(ctx, "POST",
		fmt.Sprintf("%s/api/v1/segments/%d/remove_customers", c.URL, segmentID),
		map[string]interface{}{
			"ids": ids,
		})
}
