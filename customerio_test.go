package customerio_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/customerio/go-customerio/v3"
)

type httpClientFunc func(*http.Request) (*http.Response, error)

func (f httpClientFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type testCase struct {
	id     string
	method string
	path   string
	body   any
}

type requestRecord struct {
	method string
	path   string
	body   map[string]any
}

// trackServer creates a per-test HTTP server and CustomerIO client.
// The server records request details into the returned requestRecord
// so tests can assert against them after making a client call.
func trackServer(t *testing.T) (*customerio.CustomerIO, *requestRecord) {
	t.Helper()
	rec := &requestRecord{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() { _ = req.Body.Close() }()

		// Validate basic auth
		s := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 || s[0] != "Basic" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		decoded, err := base64.URLEncoding.DecodeString(s[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		pair := strings.SplitN(string(decoded), ":", 2)
		if len(pair) != 2 || pair[0] != "siteid" || pair[1] != "apikey" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if req.Method != "DELETE" && req.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "expected Content-Type application/json", http.StatusBadRequest)
			return
		}

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

		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	client := customerio.NewTrackClient("siteid", "apikey")
	client.URL = srv.URL
	return client, rec
}

// assertRequest verifies that the recorded request matches the expected
// method, path, and body. The body parameter may be nil (for no-body
// requests like DELETE), a map/struct (compared via JSON marshaling),
// or a raw JSON string.
func assertRequest(t *testing.T, rec *requestRecord, method, path string, body any) {
	t.Helper()
	if rec.method != method {
		t.Errorf("expected method %s got %s", method, rec.method)
	}
	if rec.path != path {
		t.Errorf("expected path %s got %s", path, rec.path)
	}
	if body == nil && rec.body == nil {
		return
	}

	// If body is a raw JSON string, decode it so we can compare normalized JSON.
	expected := body
	if s, ok := body.(string); ok {
		var parsed map[string]any
		if err := json.Unmarshal([]byte(s), &parsed); err != nil {
			t.Fatalf("failed to parse expected body string as JSON: %v", err)
		}
		expected = parsed
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}
	gotJSON, err := json.Marshal(rec.body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expectedJSON, gotJSON) {
		t.Errorf("body mismatch\nexpected: %s\ngot:      %s", expectedJSON, gotJSON)
	}
}

func runCases(t *testing.T, rec *requestRecord, cases []testCase, do func(c testCase) error) {
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if err := do(c); err != nil {
				t.Fatal(err)
			}
			assertRequest(t, rec, c.method, c.path, c.body)
		})
	}
}

func checkParamError(t *testing.T, err error, param string) {
	if err == nil {
		t.Error("expected error")
		return
	}
	pe, ok := err.(customerio.ParamError)
	if !ok {
		t.Error("expected ParamError")
	}
	if pe.Param != param {
		t.Errorf("expected %s got %s", param, pe.Param)
	}
}

func TestIdentify(t *testing.T) {
	client, rec := trackServer(t)

	attributes := map[string]any{
		"a": "1",
	}
	err := client.Identify("", attributes)
	checkParamError(t, err, "customerID")

	runCases(t, rec,
		[]testCase{
			{"1", "PUT", "/api/v1/customers/1", attributes},
			{"1 ", "PUT", "/api/v1/customers/1%20", attributes},
			{"1/", "PUT", "/api/v1/customers/1%2F", attributes},
		},
		func(c testCase) error {
			return client.Identify(c.id, attributes)
		})
}

func TestBasicAuthUsesURLSafeBase64(t *testing.T) {
	siteID := "~~~"
	apiKey := "~~~"
	expectedAuth := "Basic " + base64.URLEncoding.EncodeToString([]byte(siteID+":"+apiKey))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if got := req.Header.Get("Authorization"); got != expectedAuth {
			t.Errorf("expected Authorization %q got %q", expectedAuth, got)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := customerio.NewTrackClient(siteID, apiKey)
	client.URL = srv.URL

	if err := client.Identify("1", map[string]any{"a": "1"}); err != nil {
		t.Fatal(err)
	}
}

func TestTrack(t *testing.T) {
	client, rec := trackServer(t)

	data := map[string]any{
		"a": "1",
	}

	body := map[string]any{
		"name": "test",
		"data": map[string]any{
			"a": "1",
		},
	}
	err := client.Track("", "test", data)
	checkParamError(t, err, "customerID")
	err = client.Track("1", "", data)
	checkParamError(t, err, "eventName")

	runCases(t, rec,
		[]testCase{
			{"1", "POST", "/api/v1/customers/1/events", body},
			{"1 ", "POST", "/api/v1/customers/1%20/events", body},
			{"1/", "POST", "/api/v1/customers/1%2F/events", body},
		},
		func(c testCase) error {
			return client.Track(c.id, "test", data)
		})
}

func TestTrackWithOptions(t *testing.T) {
	client, rec := trackServer(t)

	data := map[string]any{
		"a": "1",
	}
	timestamp := time.Unix(1640995200, 0)

	body := map[string]any{
		"name":      "test",
		"id":        "evt_123",
		"timestamp": timestamp.Unix(),
		"type":      customerio.TrackTypePage,
		"data": map[string]any{
			"a": "1",
		},
	}

	if err := client.Track(
		"1",
		"test",
		data,
		customerio.WithEventID("evt_123"),
		customerio.WithEventTimestamp(timestamp),
		customerio.WithEventType(customerio.TrackTypePage),
	); err != nil {
		t.Fatal(err)
	}
	assertRequest(t, rec, "POST", "/api/v1/customers/1/events", body)
}

func TestTrackAnonymous(t *testing.T) {
	client, rec := trackServer(t)

	data := map[string]any{
		"a": "1",
	}

	body := map[string]any{
		"name":         "test",
		"anonymous_id": "anon123",
		"data": map[string]any{
			"a": "1",
		},
	}

	if err := client.TrackAnonymous("anon123", "test", data); err != nil {
		t.Fatal(err)
	}
	assertRequest(t, rec, "POST", "/api/v1/events", body)
}

func TestTrackAnonymousAllowsEmptyAnonymousID(t *testing.T) {
	client, rec := trackServer(t)

	data := map[string]any{
		"a": "1",
	}

	body := map[string]any{
		"name": "test",
		"data": map[string]any{
			"a": "1",
		},
	}

	if err := client.TrackAnonymous("", "test", data); err != nil {
		t.Fatal(err)
	}
	assertRequest(t, rec, "POST", "/api/v1/events", body)
}

func TestTrackAnonymousWithOptions(t *testing.T) {
	client, rec := trackServer(t)

	data := map[string]any{
		"a": "1",
	}
	timestamp := time.Unix(1640995200, 0)

	body := map[string]any{
		"name":         "test",
		"anonymous_id": "anon123",
		"id":           "evt_123",
		"timestamp":    timestamp.Unix(),
		"type":         customerio.TrackTypeScreen,
		"data": map[string]any{
			"a": "1",
		},
	}

	if err := client.TrackAnonymous(
		"anon123",
		"test",
		data,
		customerio.WithEventID("evt_123"),
		customerio.WithEventTimestamp(timestamp),
		customerio.WithEventType(customerio.TrackTypeScreen),
	); err != nil {
		t.Fatal(err)
	}
	assertRequest(t, rec, "POST", "/api/v1/events", body)
}

func TestDelete(t *testing.T) {
	client, rec := trackServer(t)

	err := client.Delete("")
	checkParamError(t, err, "customerID")
	runCases(t, rec,
		[]testCase{
			{"1", "DELETE", "/api/v1/customers/1", nil},
			{"1 ", "DELETE", "/api/v1/customers/1%20", nil},
			{"1/", "DELETE", "/api/v1/customers/1%2F", nil},
		},
		func(c testCase) error {
			return client.Delete(c.id)
		})
}

func TestDeleteCtxUsesRequestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := customerio.NewTrackClient("siteid", "apikey", customerio.WithHTTPClient(httpClientFunc(func(req *http.Request) (*http.Response, error) {
		if err := req.Context().Err(); err != context.Canceled {
			t.Errorf("expected canceled request context, got %v", err)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})))

	if err := client.DeleteCtx(ctx, "1"); err != nil {
		t.Fatal(err)
	}
}

func TestTrackCtxUsesRequestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := customerio.NewTrackClient("siteid", "apikey", customerio.WithHTTPClient(httpClientFunc(func(req *http.Request) (*http.Response, error) {
		if err := req.Context().Err(); err != context.Canceled {
			t.Errorf("expected canceled request context, got %v", err)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})))

	if err := client.TrackCtx(ctx, "1", "purchase", nil); err != nil {
		t.Fatal(err)
	}
}

func TestCustomerIOErrorAccessors(t *testing.T) {
	const responseBody = "rate limited"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		if _, err := w.Write([]byte(responseBody)); err != nil {
			t.Error(err)
		}
	}))
	defer srv.Close()

	client := customerio.NewTrackClient("siteid", "apikey")
	client.URL = srv.URL

	err := client.Track("1", "purchase", nil)
	var apiErr *customerio.CustomerIOError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected CustomerIOError, got %T", err)
	}
	if apiErr.StatusCode() != http.StatusTooManyRequests {
		t.Errorf("expected status %d got %d", http.StatusTooManyRequests, apiErr.StatusCode())
	}
	if apiErr.URL() != srv.URL+"/api/v1/customers/1/events" {
		t.Errorf("unexpected url: %s", apiErr.URL())
	}
	if apiErr.Error() != fmt.Sprintf("%d: %s %s", http.StatusTooManyRequests, srv.URL+"/api/v1/customers/1/events", responseBody) {
		t.Errorf("unexpected error string: %s", apiErr.Error())
	}
	body := apiErr.Body()
	if string(body) != responseBody {
		t.Errorf("expected body %q got %q", responseBody, string(body))
	}
	body[0] = 'R'
	if string(apiErr.Body()) != responseBody {
		t.Error("Body should return a copy")
	}
}

func TestNewDeviceMarshalsToken(t *testing.T) {
	device, err := customerio.NewDevice("device-id", "ios", map[string]any{"attr1": "value1"})
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(device)
	if err != nil {
		t.Fatal(err)
	}

	var payload map[string]any
	if err := json.Unmarshal(b, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["token"] != "device-id" {
		t.Errorf("expected token to be device-id, got %v", payload["token"])
	}
	if _, ok := payload["id"]; ok {
		t.Errorf("device payload should use token, got id in %s", b)
	}
}

func TestAddDevice(t *testing.T) {
	client, rec := trackServer(t)

	err := client.AddDevice("", "d1", "ios", nil)
	checkParamError(t, err, "customerID")
	err = client.AddDevice("1", "", "ios", nil)
	checkParamError(t, err, "deviceID")
	err = client.AddDevice("1", "d1", "", nil)
	checkParamError(t, err, "platform")

	body := map[string]map[string]any{
		"device": {
			"id":         "d1",
			"platform":   "ios",
			"last_used":  "1606511962",
			"attributes": map[string]any{},
		},
	}
	runCases(t, rec,
		[]testCase{
			{"1", "PUT", "/api/v1/customers/1/devices", body},
			{"1 ", "PUT", "/api/v1/customers/1%20/devices", body},
			{"1/", "PUT", "/api/v1/customers/1%2F/devices", body},
		},
		func(c testCase) error {
			return client.AddDevice(c.id, "d1", "ios", map[string]any{
				"last_used": 1606511962,
			})
		})
}

func TestDeleteDevice(t *testing.T) {
	client, rec := trackServer(t)

	err := client.DeleteDevice("", "d1")
	checkParamError(t, err, "customerID")

	err = client.DeleteDevice("1", "")
	checkParamError(t, err, "deviceID")

	runCases(t, rec,
		[]testCase{
			{"1", "DELETE", "/api/v1/customers/1/devices/d1", nil},
			{"1 ", "DELETE", "/api/v1/customers/1%20/devices/d1", nil},
			{"1/", "DELETE", "/api/v1/customers/1%2F/devices/d1", nil},
			{"2", "DELETE", "/api/v1/customers/d1/devices/2", nil},
			{"2 ", "DELETE", "/api/v1/customers/d1/devices/2%20", nil},
			{"2/", "DELETE", "/api/v1/customers/d1/devices/2%2F", nil},
		},
		func(c testCase) error {
			if c.id[0] == '2' {
				return client.DeleteDevice("d1", c.id)
			} else {
				return client.DeleteDevice(c.id, "d1")
			}
		})
}

func TestRejectsMismatchedBasicAuth(t *testing.T) {
	for _, tc := range []struct {
		name   string
		siteID string
		apiKey string
	}{
		{
			name:   "wrong site id",
			siteID: "wrong",
			apiKey: "apikey",
		},
		{
			name:   "wrong api key",
			siteID: "siteid",
			apiKey: "wrong",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Create a server with auth checking (same as trackServer's handler).
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				s := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
				if len(s) != 2 || s[0] != "Basic" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				decoded, err := base64.URLEncoding.DecodeString(s[1])
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				pair := strings.SplitN(string(decoded), ":", 2)
				if len(pair) != 2 || pair[0] != "siteid" || pair[1] != "apikey" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusOK)
			}))
			t.Cleanup(srv.Close)

			client := customerio.NewTrackClient(tc.siteID, tc.apiKey)
			client.URL = srv.URL

			err := client.Identify("1", map[string]any{"a": "1"})

			var apiErr *customerio.CustomerIOError
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected CustomerIOError, got %T (%v)", err, err)
			}
			if apiErr.StatusCode() != http.StatusUnauthorized {
				t.Errorf("expected status %d got %d", http.StatusUnauthorized, apiErr.StatusCode())
			}
		})
	}
}

func TestMergeCustomers(t *testing.T) {
	client, rec := trackServer(t)

	checkMergeError := func(t *testing.T, err error, prefix, contains string) {
		t.Helper()
		if err == nil {
			t.Fatal("expected error")
		}
		if got := err.Error(); !strings.Contains(got, prefix+": ") || !strings.Contains(got, contains) {
			t.Errorf("expected error containing %q and %q, got %q", prefix, contains, got)
		}
	}

	err1 := client.MergeCustomers(customerio.Identifier{
		Type:  "",
		Value: "id1",
	}, customerio.Identifier{
		Type:  "id",
		Value: "id2",
	})
	checkMergeError(t, err1, "primary", "invalid id type")

	err2 := client.MergeCustomers(customerio.Identifier{
		Type:  "id",
		Value: "",
	}, customerio.Identifier{
		Type:  "id",
		Value: "id2",
	})
	checkMergeError(t, err2, "primary", "invalid id")

	err3 := client.MergeCustomers(customerio.Identifier{
		Type:  "email",
		Value: "id1",
	}, customerio.Identifier{
		Type:  "",
		Value: "id2",
	})
	checkMergeError(t, err3, "secondary", "invalid id type")

	err4 := client.MergeCustomers(customerio.Identifier{
		Type:  "cio_id",
		Value: "id1",
	}, customerio.Identifier{
		Type:  "email",
		Value: "",
	})
	checkMergeError(t, err4, "secondary", "invalid id")

	runCases(t, rec,
		[]testCase{
			{"1", "POST", "/api/v1/merge_customers", `{"primary":{"email":"cool.person@company.com"},"secondary":{"email":"cperson@gmail.com"}}`},
			{"2", "POST", "/api/v1/merge_customers", `{"primary":{"id":"cool.person@company.com"},"secondary":{"cio_id":"person2"}}`},
			{"3", "POST", "/api/v1/merge_customers", `{"primary":{"cio_id":"CIO123"},"secondary":{"id":"person1"}}`},
		},
		func(c testCase) error {
			switch c.id {
			case "1":
				return client.MergeCustomers(customerio.Identifier{
					Type:  "email",
					Value: "cool.person@company.com",
				}, customerio.Identifier{
					Type:  "email",
					Value: "cperson@gmail.com",
				})
			case "2":
				return client.MergeCustomers(customerio.Identifier{
					Type:  "id",
					Value: "cool.person@company.com",
				}, customerio.Identifier{
					Type:  "cio_id",
					Value: "person2",
				})
			default:
				return client.MergeCustomers(customerio.Identifier{
					Type:  customerio.IdentifierTypeCioID,
					Value: "CIO123",
				}, customerio.Identifier{
					Type:  customerio.IdentifierTypeID,
					Value: "person1",
				})
			}
		})
}
