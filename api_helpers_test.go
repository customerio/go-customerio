package customerio_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

// apiRecord captures the most recent request made through an apiServer.
type apiRecord struct {
	method string
	path   string
	body   map[string]any
}

// apiServer creates a per-test HTTP server and APIClient for App API tests.
// Every request gets statusCode and responseBody (JSON-encoded); the
// request actually sent is captured into the returned apiRecord so tests
// can assert against it afterwards.
func apiServer(t *testing.T, statusCode int, responseBody any) (*customerio.APIClient, *apiRecord) {
	t.Helper()
	rec := &apiRecord{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() { _ = req.Body.Close() }()

		rec.method = req.Method
		rec.path = req.RequestURI
		rec.body = nil
		if len(b) > 0 {
			dec := json.NewDecoder(bytes.NewReader(b))
			dec.UseNumber()
			if err := dec.Decode(&rec.body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(statusCode)
		if responseBody != nil {
			// errcheck (golangci-lint): explicitly discard return values
			_ = json.NewEncoder(w).Encode(responseBody)
		}
	}))
	t.Cleanup(srv.Close)

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL
	return api, rec
}

func assertAPIRequest(t *testing.T, rec *apiRecord, method, path string) {
	t.Helper()
	if rec.method != method {
		t.Errorf("expected method %s got %s", method, rec.method)
	}
	if rec.path != path {
		t.Errorf("expected path %s got %s", path, rec.path)
	}
}

// assertJSONEqual compares got and want by marshaling both to JSON, rather
// than reflect.DeepEqual — apiServer decodes request bodies with UseNumber,
// so a captured numeric field is a json.Number, which never DeepEqual's a
// float64 literal in a hand-written expectation even when both serialize
// identically.
func assertJSONEqual(t *testing.T, got, want any) {
	t.Helper()
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotJSON, wantJSON) {
		t.Errorf("body mismatch\nexpected: %s\ngot:      %s", wantJSON, gotJSON)
	}
}
