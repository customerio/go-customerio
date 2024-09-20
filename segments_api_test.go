package customerio_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

var (
	testSegmentID     = 1
	notFoundID        = 2
	testDeduplicateID = "deduplicate_id"
)

func TestCreateSegment(t *testing.T) {
	createSegmentRequest := &customerio.CreateSegmentRequest{
		Segment: customerio.Segment{
			Name:        "name",
			Description: "description",
			State:       "state",
			Progress:    intPtr(1),
			Type:        "type",
			Tags:        []string{"tags"},
		},
	}

	var verify = func(request []byte) {
		var body customerio.CreateSegmentRequest
		if err := json.Unmarshal(request, &body); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(&body, createSegmentRequest) {
			t.Errorf("Request differed, want: %#v, got: %#v", request, body)
		}
	}

	api, srv := segmentsAppServer(t, verify)
	defer srv.Close()

	resp, err := api.CreateSegment(context.Background(), createSegmentRequest)
	if err != nil {
		t.Error(err)
	}

	expect := &customerio.CreateSegmentResponse{
		Segment: customerio.Segment{
			ID:            testSegmentID,
			DeduplicateID: testDeduplicateID,
		},
	}

	if !reflect.DeepEqual(resp, expect) {
		t.Errorf("Expect: %#v, Got: %#v", expect, resp)
	}
}

func TestCreateSegmentDoRequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.CreateSegment(context.Background(), &customerio.CreateSegmentRequest{
		Segment: customerio.Segment{
			Name:        "name",
			Description: "description",
			State:       "state",
			Progress:    intPtr(1),
			Type:        "type",
			Tags:        []string{"tags"},
		},
	})
	if err == nil {
		t.Errorf("Expected error due to request failure, got: nil")
	}
}

func TestCreateSegmentNotOKError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(502)
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.CreateSegment(context.Background(), &customerio.CreateSegmentRequest{
		Segment: customerio.Segment{
			Name:        "name",
			Description: "description",
			State:       "state",
			Progress:    intPtr(1),
			Type:        "type",
			Tags:        []string{"tags"},
		},
	})
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}
}

func TestCreateSegmentUnmarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`invalid-json`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.CreateSegment(context.Background(), &customerio.CreateSegmentRequest{
		Segment: customerio.Segment{
			Name:        "name",
			Description: "description",
			State:       "state",
			Progress:    intPtr(1),
			Type:        "type",
			Tags:        []string{"tags"},
		},
	})

	if err == nil {
		t.Errorf("Expected error due to invalid JSON, got: nil")
	}
}

func TestListSegments(t *testing.T) {
	var verify = func(request []byte) {}

	api, srv := segmentsAppServer(t, verify)
	defer srv.Close()

	resp, err := api.ListSegments(context.Background())
	if err != nil {
		t.Error(err)
	}

	expect := &customerio.ListSegmentsResponse{
		Segments: []customerio.Segment{
			{
				ID:            testSegmentID,
				DeduplicateID: testDeduplicateID,
			},
		},
	}

	if !reflect.DeepEqual(resp, expect) {
		t.Errorf("Expect: %#v, Got: %#v", expect, resp)
	}
}

func TestListSegmentsDoRequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.ListSegments(context.Background())
	if err == nil {
		t.Errorf("Expected error due to request failure, got: nil")
	}
}

func TestListSegmentsNotOKError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(502)
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.ListSegments(context.Background())
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}
}

func TestListSegmentsUnmarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`invalid-json`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.ListSegments(context.Background())
	if err == nil {
		t.Errorf("Expected error due to invalid JSON, got: nil")
	}
}

func TestGetSegment(t *testing.T) {
	var verify = func(request []byte) {}

	api, srv := segmentsAppServer(t, verify)
	defer srv.Close()

	resp, err := api.GetSegment(context.Background(), testSegmentID)
	if err != nil {
		t.Error(err)
	}

	expect := &customerio.GetSegmentResponse{
		Segment: customerio.Segment{
			ID:            testSegmentID,
			DeduplicateID: testDeduplicateID,
		},
	}

	if !reflect.DeepEqual(resp, expect) {
		t.Errorf("Expect: %#v, Got: %#v", expect, resp)
	}
}

func TestGetSegmentDoRequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.GetSegment(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error due to request failure, got: nil")
	}
}

func TestGetSegmentNotOKError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(502)
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.GetSegment(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}
}

func TestGetSegmentUnmarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`invalid-json`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.GetSegment(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error due to invalid JSON, got: nil")
	}
}

func TestDeleteSegment(t *testing.T) {
	var verify = func(request []byte) {}

	api, srv := segmentsAppServer(t, verify)
	defer srv.Close()

	err := api.DeleteSegment(context.Background(), testSegmentID)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteSegmentDoRequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	err := api.DeleteSegment(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error due to request failure, got: nil")
	}
}

func TestDeleteSegmentNotOKError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(502)
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	err := api.DeleteSegment(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}
}

func TestDeleteSegmentUnmarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`invalid-json`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	err := api.DeleteSegment(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error due to invalid JSON, got: nil")
	}
}

func TestGetSegmentDependencies(t *testing.T) {
	var verify = func(request []byte) {}

	api, srv := segmentsAppServer(t, verify)
	defer srv.Close()

	resp, err := api.GetSegmentDependencies(context.Background(), testSegmentID)
	if err != nil {
		t.Error(err)
	}

	expect := &customerio.GetSegmentDependenciesResponse{
		UsedBy: struct {
			Campaigns        []int `json:"campaigns"`
			SentNewsletters  []int `json:"sent_newsletters"`
			DraftNewsletters []int `json:"draft_newsletters"`
		}{
			Campaigns:        []int{1},
			SentNewsletters:  []int{2},
			DraftNewsletters: []int{3},
		},
	}

	if !reflect.DeepEqual(resp, expect) {
		t.Errorf("Expect: %#v, Got: %#v", expect, resp)
	}
}

func TestGetSegmentDependenciesDoRequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.GetSegmentDependencies(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error due to request failure, got: nil")
	}
}

func TestGetSegmentDependenciesNotOKError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(502)
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.GetSegmentDependencies(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}
}

func TestGetSegmentDependenciesUnmarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`invalid-json`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.GetSegmentDependencies(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error due to invalid JSON, got: nil")
	}
}

func TestGetSegmentCustomerCount(t *testing.T) {
	var verify = func(request []byte) {}

	api, srv := segmentsAppServer(t, verify)
	defer srv.Close()

	resp, err := api.GetSegmentCustomerCount(context.Background(), testSegmentID)
	if err != nil {
		t.Error(err)
	}

	expect := &customerio.GetSegmentCustomerCountResponse{
		Count: 1,
	}

	if !reflect.DeepEqual(resp, expect) {
		t.Errorf("Expect: %#v, Got: %#v", expect, resp)
	}
}

func TestGetSegmentCustomerCountDoRequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.GetSegmentCustomerCount(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error due to request failure, got: nil")
	}
}

func TestGetSegmentCustomerCountNotOKError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(502)
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.GetSegmentCustomerCount(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}
}

func TestGetSegmentCustomerCountUnmarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`invalid-json`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.GetSegmentCustomerCount(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error due to invalid JSON, got: nil")
	}
}

func TestListCustomersInSegment(t *testing.T) {
	var verify = func(request []byte) {}

	api, srv := segmentsAppServer(t, verify)
	defer srv.Close()

	resp, err := api.ListCustomersInSegment(context.Background(), testSegmentID)
	if err != nil {
		t.Error(err)
	}

	expect := &customerio.ListCustomersInSegmentResponse{
		IDs: []string{"string"},
		Identifiers: []customerio.CustomerIdentifier{
			{
				Email: "test@example.com",
				ID:    2,
				CioID: "a3000001",
			},
		},
		Next: "string",
	}

	if !reflect.DeepEqual(resp, expect) {
		t.Errorf("Expect: %#v, Got: %#v", expect, resp)
	}
}

func TestListCustomersInSegmentDoRequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.ListCustomersInSegment(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error due to request failure, got: nil")
	}
}

func TestListCustomersInSegmentNotOKError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(502)
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.ListCustomersInSegment(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}
}

func TestListCustomersInSegmentUnmarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`invalid-json`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	_, err := api.ListCustomersInSegment(context.Background(), testSegmentID)
	if err == nil {
		t.Errorf("Expected error due to invalid JSON, got: nil")
	}
}

func intPtr(s int) *int {
	return &s
}

func segmentsAppServer(t *testing.T, verify func(request []byte)) (*customerio.APIClient, *httptest.Server) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
		}
		defer req.Body.Close()

		verify(b)

		switch true {
		case req.Method == "POST" && req.URL.Path == "/v1/segments":
			w.Write([]byte(`{
				"segment": {
					"id": ` + fmt.Sprintf("%d", testSegmentID) + `,
					"deduplicate_id": "` + testDeduplicateID + `"
			}
		  }`))
		case req.Method == "GET" && req.URL.Path == "/v1/segments":
			w.Write([]byte(`{
				"segments": [
					{
						"id": ` + fmt.Sprintf("%d", testSegmentID) + `,
						"deduplicate_id": "` + testDeduplicateID + `"
			}
		  ]}`))
		case req.Method == "GET" && req.URL.Path == "/v1/segments/1":
			w.Write([]byte(`{
				"segment": {
					"id": ` + fmt.Sprintf("%d", testSegmentID) + `,
					"deduplicate_id": "` + testDeduplicateID + `"
			}
		  }`))
		case req.Method == "DELETE" && req.URL.Path == "/v1/segments/1":
			w.WriteHeader(http.StatusNoContent)
		case req.Method == "GET" && req.URL.Path == "/v1/segments/1/customer_count":
			w.Write([]byte(`{
				"count": 1
			}`))
		case req.Method == "GET" && req.URL.Path == "/v1/segments/1/membership":
			w.Write([]byte(`{
				"ids": ["string"],
				"identifiers": [
					{
						"email": "test@example.com",
						"id": 2,
						"cio_id": "a3000001"
					}
				],
				"next": "string"
			}`))
		case req.Method == "GET" && req.URL.Path == "/v1/segments/1/used_by":
			w.Write([]byte(`{
				"used_by": {
					"campaigns": [1],
					"sent_newsletters": [2],
					"draft_newsletters": [3]
				}
			}`))
		default:
			t.Errorf("Unexpected request: %s %s", req.Method, req.URL.Path)
		}
	}))

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	return api, srv
}
