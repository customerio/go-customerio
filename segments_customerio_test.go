package customerio_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestAddPeopleToSegment(t *testing.T) {
	customerIDs := []string{"1", "2", "3"}
	var verify = func(request []byte) {}

	api, srv := segmentsTrackServer(t, verify)
	defer srv.Close()

	err := api.AddPeopleToSegment(context.Background(), testSegmentID, customerIDs)
	if err != nil {
		t.Error(err)
	}
}

func TestAddPeopleToSegmentSegmentParamError(t *testing.T) {
	var customerIDs []string
	var verify = func(request []byte) {}

	api, srv := segmentsTrackServer(t, verify)
	defer srv.Close()

	err := api.AddPeopleToSegment(context.Background(), 0, customerIDs)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}

	if e, ok := err.(customerio.ParamError); !ok {
		t.Errorf("Expected ParamError, got: %#v", e)
	}
}

func TestAddPeopleToSegmentIDsParamError(t *testing.T) {
	var customerIDs []string
	var verify = func(request []byte) {}

	api, srv := segmentsTrackServer(t, verify)
	defer srv.Close()

	err := api.AddPeopleToSegment(context.Background(), testSegmentID, customerIDs)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}

	if e, ok := err.(customerio.ParamError); !ok {
		t.Errorf("Expected ParamError, got: %#v", e)
	}
}

func TestAddPeopleToSegmentError(t *testing.T) {
	customerIDs := []string{"1", "2", "3"}
	var verify = func(request []byte) {}

	api, srv := segmentsTrackServer(t, verify)
	defer srv.Close()

	err := api.AddPeopleToSegment(context.Background(), notFoundID, customerIDs)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}

	if e, ok := err.(*customerio.CustomerIOError); !ok {
		t.Errorf("Expected CustomerIOError, got: %#v", e)
	}
}

func segmentsTrackServer(t *testing.T, verify func(request []byte)) (*customerio.CustomerIO, *httptest.Server) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
		}
		defer req.Body.Close()

		verify(b)

		switch true {
		case req.Method == "POST" && req.URL.Path == "/api/v1/segments/1/add_customers":
			w.WriteHeader(http.StatusOK)
		case req.Method == "POST" && req.URL.Path == "/api/v1/segments/2/add_customers":
			w.WriteHeader(http.StatusNotFound)
		case req.Method == "POST" && req.URL.Path == "/api/v1/segments/1/remove_customers":
			w.WriteHeader(http.StatusOK)
		case req.Method == "POST" && req.URL.Path == "/api/v1/segments/2/remove_customers":
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Errorf("Unexpected request: %s %s", req.Method, req.URL.Path)
		}
	}))

	api := customerio.NewCustomerIO("test", "myKey")
	api.URL = srv.URL

	return api, srv
}

func TestRemovePeopleFromSegment(t *testing.T) {
	customerIDs := []string{"1", "2", "3"}
	var verify = func(request []byte) {}

	api, srv := segmentsTrackServer(t, verify)
	defer srv.Close()

	err := api.RemovePeopleFromSegment(context.Background(), testSegmentID, customerIDs)
	if err != nil {
		t.Error(err)
	}
}

func TestRemovePeopleFromSegmentSegmentParamError(t *testing.T) {
	var customerIDs []string
	var verify = func(request []byte) {}

	api, srv := segmentsTrackServer(t, verify)
	defer srv.Close()

	err := api.RemovePeopleFromSegment(context.Background(), 0, customerIDs)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}

	if e, ok := err.(customerio.ParamError); !ok {
		t.Errorf("Expected ParamError, got: %#v", e)
	}
}

func TestRemovePeopleFromSegmentIDsParamError(t *testing.T) {
	var customerIDs []string
	var verify = func(request []byte) {}

	api, srv := segmentsTrackServer(t, verify)
	defer srv.Close()

	err := api.RemovePeopleFromSegment(context.Background(), testSegmentID, customerIDs)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}

	if e, ok := err.(customerio.ParamError); !ok {
		t.Errorf("Expected ParamError, got: %#v", e)
	}
}

func TestRemovePeopleFromSegmentError(t *testing.T) {
	customerIDs := []string{"1", "2", "3"}
	var verify = func(request []byte) {}

	api, srv := segmentsTrackServer(t, verify)
	defer srv.Close()

	err := api.RemovePeopleFromSegment(context.Background(), notFoundID, customerIDs)
	if err == nil {
		t.Errorf("Expected error, got: %#v", err)
	}

	if e, ok := err.(*customerio.CustomerIOError); !ok {
		t.Errorf("Expected CustomerIOError, got: %#v", e)
	}
}
