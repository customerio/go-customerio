package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

var (
	idsByID    = []string{"1", "2", "3"}
	idsByEmail = []string{"alice@example.com", "bob@example.com"}
	idsByCioID = []string{"cio_abc123", "cio_def456"}
)

func TestAddPeopleToSegment(t *testing.T) {
	ctx := context.Background()

	checkParamError(t, cio.AddPeopleToSegment(ctx, 0, idsByID), "segmentID")
	checkParamError(t, cio.AddPeopleToSegment(ctx, -1, idsByID), "segmentID")
	checkParamError(t, cio.AddPeopleToSegment(ctx, 7, nil), "ids")
	checkParamError(t, cio.AddPeopleToSegment(ctx, 7, []string{}), "ids")

	runCases(t,
		[]testCase{
			{"default", "POST", "/api/v1/segments/7/add_customers", map[string]interface{}{"ids": idsByID}},
			{"email", "POST", "/api/v1/segments/7/add_customers?id_type=email", map[string]interface{}{"ids": idsByEmail}},
			{"cio_id", "POST", "/api/v1/segments/7/add_customers?id_type=cio_id", map[string]interface{}{"ids": idsByCioID}},
		},
		func(c testCase) error {
			switch c.id {
			case "email":
				return cio.AddPeopleToSegment(ctx, 7, idsByEmail, customerio.WithSegmentIDType(customerio.IdentifierTypeEmail))
			case "cio_id":
				return cio.AddPeopleToSegment(ctx, 7, idsByCioID, customerio.WithSegmentIDType(customerio.IdentifierTypeCioID))
			default:
				return cio.AddPeopleToSegment(ctx, 7, idsByID)
			}
		})
}

func TestRemovePeopleFromSegment(t *testing.T) {
	ctx := context.Background()

	checkParamError(t, cio.RemovePeopleFromSegment(ctx, 0, idsByID), "segmentID")
	checkParamError(t, cio.RemovePeopleFromSegment(ctx, -1, idsByID), "segmentID")
	checkParamError(t, cio.RemovePeopleFromSegment(ctx, 7, nil), "ids")
	checkParamError(t, cio.RemovePeopleFromSegment(ctx, 7, []string{}), "ids")

	runCases(t,
		[]testCase{
			{"default", "POST", "/api/v1/segments/7/remove_customers", map[string]interface{}{"ids": idsByID}},
			{"email", "POST", "/api/v1/segments/7/remove_customers?id_type=email", map[string]interface{}{"ids": idsByEmail}},
			{"cio_id", "POST", "/api/v1/segments/7/remove_customers?id_type=cio_id", map[string]interface{}{"ids": idsByCioID}},
		},
		func(c testCase) error {
			switch c.id {
			case "email":
				return cio.RemovePeopleFromSegment(ctx, 7, idsByEmail, customerio.WithSegmentIDType(customerio.IdentifierTypeEmail))
			case "cio_id":
				return cio.RemovePeopleFromSegment(ctx, 7, idsByCioID, customerio.WithSegmentIDType(customerio.IdentifierTypeCioID))
			default:
				return cio.RemovePeopleFromSegment(ctx, 7, idsByID)
			}
		})
}
