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
	client, rec := trackServer(t)
	ctx := context.Background()

	checkParamError(t, client.AddPeopleToSegment(ctx, 0, idsByID), "segmentID")
	checkParamError(t, client.AddPeopleToSegment(ctx, -1, idsByID), "segmentID")
	checkParamError(t, client.AddPeopleToSegment(ctx, 7, nil), "ids")
	checkParamError(t, client.AddPeopleToSegment(ctx, 7, []string{}), "ids")

	runCases(t, rec,
		[]testCase{
			{"default", "POST", "/api/v1/segments/7/add_customers", map[string]any{"ids": idsByID}},
			{"email", "POST", "/api/v1/segments/7/add_customers?id_type=email", map[string]any{"ids": idsByEmail}},
			{"cio_id", "POST", "/api/v1/segments/7/add_customers?id_type=cio_id", map[string]any{"ids": idsByCioID}},
		},
		func(c testCase) error {
			switch c.id {
			case "email":
				return client.AddPeopleToSegment(ctx, 7, idsByEmail, customerio.WithSegmentIDType(customerio.IdentifierTypeEmail))
			case "cio_id":
				return client.AddPeopleToSegment(ctx, 7, idsByCioID, customerio.WithSegmentIDType(customerio.IdentifierTypeCioID))
			default:
				return client.AddPeopleToSegment(ctx, 7, idsByID)
			}
		})
}

func TestRemovePeopleFromSegment(t *testing.T) {
	client, rec := trackServer(t)
	ctx := context.Background()

	checkParamError(t, client.RemovePeopleFromSegment(ctx, 0, idsByID), "segmentID")
	checkParamError(t, client.RemovePeopleFromSegment(ctx, -1, idsByID), "segmentID")
	checkParamError(t, client.RemovePeopleFromSegment(ctx, 7, nil), "ids")
	checkParamError(t, client.RemovePeopleFromSegment(ctx, 7, []string{}), "ids")

	runCases(t, rec,
		[]testCase{
			{"default", "POST", "/api/v1/segments/7/remove_customers", map[string]any{"ids": idsByID}},
			{"email", "POST", "/api/v1/segments/7/remove_customers?id_type=email", map[string]any{"ids": idsByEmail}},
			{"cio_id", "POST", "/api/v1/segments/7/remove_customers?id_type=cio_id", map[string]any{"ids": idsByCioID}},
		},
		func(c testCase) error {
			switch c.id {
			case "email":
				return client.RemovePeopleFromSegment(ctx, 7, idsByEmail, customerio.WithSegmentIDType(customerio.IdentifierTypeEmail))
			case "cio_id":
				return client.RemovePeopleFromSegment(ctx, 7, idsByCioID, customerio.WithSegmentIDType(customerio.IdentifierTypeCioID))
			default:
				return client.RemovePeopleFromSegment(ctx, 7, idsByID)
			}
		})
}
