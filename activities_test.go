package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestListActivities(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"activities": []map[string]any{
		{"id": "act1", "type": "event", "timestamp": 1700000000, "name": "purchased"},
	}, "next": "cursor2"})

	deleted := false
	got, err := api.ListActivities(context.Background(), customerio.ListActivitiesOptions{
		Type:              "event",
		Name:              "purchased",
		Deleted:           &deleted,
		CustomerID:        "1",
		IDType:            customerio.IdentifierTypeID,
		PaginationOptions: customerio.PaginationOptions{Limit: 25},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/activities?customer_id=1&deleted=false&id_type=id&limit=25&name=purchased&type=event")
	if len(got.Activities) != 1 || got.Activities[0].Name != "purchased" || got.Next != "cursor2" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestListActivitiesNoOptions(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"activities": []any{}})

	if _, err := api.ListActivities(context.Background(), customerio.ListActivitiesOptions{}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/activities")
}
