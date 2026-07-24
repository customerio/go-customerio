package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestListSegmentsApp(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"segments": []map[string]any{
		{"id": 1, "name": "VIP", "type": "manual"},
	}})
	got, err := api.ListSegments(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/segments")
	if len(got.Segments) != 1 || got.Segments[0].Name != "VIP" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateSegment(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateSegment(context.Background(), customerio.SegmentInput{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "name")
	}

	api, rec := apiServer(t, 200, map[string]any{"segment": map[string]any{"id": 9, "name": "New"}})
	got, err := api.CreateSegment(context.Background(), customerio.SegmentInput{Name: "New", Description: "desc"})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/segments")
	wantBody := map[string]any{"segment": map[string]any{"name": "New", "description": "desc"}}
	assertJSONEqual(t, rec.body, wantBody)
	if got.Segment.ID != 9 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetSegmentApp(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetSegment(context.Background(), 0); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "segmentID")
	}

	api, rec := apiServer(t, 200, map[string]any{"segment": map[string]any{"id": 7, "name": "Gold"}})
	got, err := api.GetSegment(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/segments/7")
	if got.Segment.Name != "Gold" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestDeleteSegment(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if err := api.DeleteSegment(context.Background(), -1); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "segmentID")
	}

	api, rec := apiServer(t, 204, nil)
	if err := api.DeleteSegment(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/segments/7")
}

func TestGetSegmentCustomerCount(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"segment_id": 7, "count": 42})
	got, err := api.GetSegmentCustomerCount(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/segments/7/customer_count")
	if got.Count != 42 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetSegmentMembership(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{
		"segment_id":  7,
		"ids":         []string{"1", "2"},
		"identifiers": []map[string]any{{"id": "1"}, {"id": "2"}},
		"next":        "cursor",
	})
	got, err := api.GetSegmentMembership(context.Background(), 7, customerio.PaginationOptions{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/segments/7/membership?limit=2")
	if len(got.IDs) != 2 || got.Next != "cursor" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetSegmentUsedBy(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{
		"segment_id": 7,
		"used_by":    map[string]any{"campaigns": []int{1, 2}, "sent_newsletters": []int{3}},
	})
	got, err := api.GetSegmentUsedBy(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/segments/7/used_by")
	if len(got.UsedBy.Campaigns) != 2 || len(got.UsedBy.SentNewsletters) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestAppSegmentsError(t *testing.T) {
	api, _ := apiServer(t, 404, map[string]any{"errors": []map[string]any{{"detail": "not found"}}})

	_, err := api.GetSegment(context.Background(), 1)
	cioErr, ok := err.(*customerio.CustomerIOError)
	if !ok {
		t.Fatalf("expected *CustomerIOError, got %T", err)
	}
	if cioErr.StatusCode() != 404 {
		t.Errorf("expected status 404, got %d", cioErr.StatusCode())
	}
}
