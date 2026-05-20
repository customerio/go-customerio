package customerio

import (
	"context"
	"net/url"
	"strconv"
)

// IDType is the type of ids you want to use.
// All the values in the ids array must be of this type.
// If you don't provide this parameter, we assume that the ids array contains id values.
// Enum values:
//   - "id"
//   - "email"
//   - "cio_id"
type IDType string

const (
	IDTypeID      IDType = "id"
	IDTypeEmail   IDType = "email"
	IDTypeCioID   IDType = "cio_id"
	DefaultIDType        = IDTypeID
)

// AddPeopleToSegment adds people to a segment.
func (c *CustomerIO) AddPeopleToSegment(ctx context.Context, segmentID int, ids []string) error {
	if segmentID == 0 {
		return ParamError{Param: "segmentID"}
	}
	if len(ids) == 0 {
		return ParamError{Param: "ids"}
	}
	u, err := buildURL(c.URL, c.idTypeQuery(), "api", "v1", "segments", strconv.Itoa(segmentID), "add_customers")
	if err != nil {
		return err
	}
	return c.request(ctx, "POST", u, map[string]interface{}{
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
	u, err := buildURL(c.URL, c.idTypeQuery(), "api", "v1", "segments", strconv.Itoa(segmentID), "remove_customers")
	if err != nil {
		return err
	}
	return c.request(ctx, "POST", u, map[string]interface{}{
		"ids": ids,
	})
}

// idTypeQuery returns the id_type query parameter for segment requests, or nil
// when no IDType is configured.
func (c *CustomerIO) idTypeQuery() url.Values {
	if c.IDType == "" {
		return nil
	}

	// Check if the IDType is valid and construct the query parameter accordingly.
	v := url.Values{}
	switch IDType(c.IDType) {
	case IDTypeEmail, IDTypeCioID:
		v.Set("id_type", c.IDType)
	default:
		v.Set("id_type", string(DefaultIDType))
	}
	return v
}
