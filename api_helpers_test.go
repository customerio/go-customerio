package customerio_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

// apiRecord captures the most recent request made through an apiServer.
// body is `any` rather than map[string]any since some endpoints
// (e.g. UpdateCollectionContent) send a top-level JSON array.
type apiRecord struct {
	method string
	path   string
	body   any
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

// assertJSONEqual compares got and want by marshaling both to JSON and then
// unmarshaling into a generic any, rather than comparing raw bytes or using
// reflect.DeepEqual directly on the inputs. Two round trips are needed: (1)
// apiServer decodes request bodies with UseNumber, so a captured numeric
// field is a json.Number, which never DeepEqual's a float64 literal in a
// hand-written expectation even when both serialize identically; (2) when
// want is a typed struct, json.Marshal preserves its field declaration
// order, while got (already a decoded map) always marshals with
// alphabetically sorted keys — comparing raw bytes would spuriously fail on
// key order alone. Round-tripping both sides through the same
// marshal-then-unmarshal-to-any path normalizes both away.
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

	var gotGeneric, wantGeneric any
	if err := json.Unmarshal(gotJSON, &gotGeneric); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(wantJSON, &wantGeneric); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(gotGeneric, wantGeneric) {
		t.Errorf("body mismatch\nexpected: %s\ngot:      %s", wantJSON, gotJSON)
	}
}
