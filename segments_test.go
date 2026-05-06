package customerio_test

import (
	"testing"

	"github.com/customerio/go-customerio/v3"
)

var (
	idsByID    = []string{"1", "2", "3"}
	idsByEmail = []string{"alice@example.com", "bob@example.com"}
	idsByCioID = []string{"cio_abc123", "cio_def456"}
)

func TestAddPeopleToSegment(t *testing.T) {
	checkParamError(t, cio.AddPeopleToSegment(0, idsByID), "segmentID")
	checkParamError(t, cio.AddPeopleToSegment(-1, idsByID), "segmentID")
	checkParamError(t, cio.AddPeopleToSegment(7, nil), "ids")
	checkParamError(t, cio.AddPeopleToSegment(7, []string{}), "ids")

	runCases(t,
		[]testCase{
			{"default", "POST", "/api/v1/segments/7/add_customers", map[string]interface{}{"ids": idsByID}},
			{"email", "POST", "/api/v1/segments/7/add_customers?id_type=email", map[string]interface{}{"ids": idsByEmail}},
			{"cio_id", "POST", "/api/v1/segments/7/add_customers?id_type=cio_id", map[string]interface{}{"ids": idsByCioID}},
		},
		func(c testCase) error {
			switch c.id {
			case "email":
				return cio.AddPeopleToSegment(7, idsByEmail, customerio.WithSegmentIDType(customerio.IdentifierTypeEmail))
			case "cio_id":
				return cio.AddPeopleToSegment(7, idsByCioID, customerio.WithSegmentIDType(customerio.IdentifierTypeCioID))
			default:
				return cio.AddPeopleToSegment(7, idsByID)
			}
		})
}

func TestRemovePeopleFromSegment(t *testing.T) {
	checkParamError(t, cio.RemovePeopleFromSegment(0, idsByID), "segmentID")
	checkParamError(t, cio.RemovePeopleFromSegment(-1, idsByID), "segmentID")
	checkParamError(t, cio.RemovePeopleFromSegment(7, nil), "ids")
	checkParamError(t, cio.RemovePeopleFromSegment(7, []string{}), "ids")

	runCases(t,
		[]testCase{
			{"default", "POST", "/api/v1/segments/7/remove_customers", map[string]interface{}{"ids": idsByID}},
			{"email", "POST", "/api/v1/segments/7/remove_customers?id_type=email", map[string]interface{}{"ids": idsByEmail}},
			{"cio_id", "POST", "/api/v1/segments/7/remove_customers?id_type=cio_id", map[string]interface{}{"ids": idsByCioID}},
		},
		func(c testCase) error {
			switch c.id {
			case "email":
				return cio.RemovePeopleFromSegment(7, idsByEmail, customerio.WithSegmentIDType(customerio.IdentifierTypeEmail))
			case "cio_id":
				return cio.RemovePeopleFromSegment(7, idsByCioID, customerio.WithSegmentIDType(customerio.IdentifierTypeCioID))
			default:
				return cio.RemovePeopleFromSegment(7, idsByID)
			}
		})
}
