package customerio_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/customerio/go-customerio"
)

var cio *customerio.CustomerIO

func TestMain(m *testing.M) {
	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()
	u, err := url.Parse(srv.URL)
	if err != nil {
		panic(err)
	}

	cio = customerio.NewCustomerIO("siteid", "apikey")
	cio.Host = u.Host
	cio.SSL = false

	os.Exit(m.Run())
}

type testCase struct {
	id     string
	method string
	path   string
	body   interface{}
}

func runCases(t *testing.T, cases []testCase, do func(c testCase) error) {
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			expect(c.method, c.path, c.body)
			if err := do(c); err != nil {
				t.Error(err.Error())
			}
		})
	}
}
func checkParamError(t *testing.T, err error, param string) {
	pe, ok := err.(customerio.ParamError)
	if !ok {
		t.Error("expected ParamError")
	}
	if pe.Param != param {
		t.Errorf("expected %s got %s", param, pe.Param)
	}
}

func TestIdentify(t *testing.T) {
	attributes := map[string]interface{}{
		"a": "1",
	}
	if err := cio.Identify("", attributes); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}

	runCases(t,
		[]testCase{
			{"1", "PUT", "/api/v1/customers/1", attributes},
			{"1 ", "PUT", "/api/v1/customers/1%20", attributes},
			{"1/", "PUT", "/api/v1/customers/1%2F", attributes},
		},
		func(c testCase) error {
			return cio.Identify(c.id, attributes)
		})
}

func TestTrack(t *testing.T) {
	data := map[string]interface{}{
		"a": "1",
	}

	body := map[string]interface{}{
		"name": "test",
		"data": map[string]interface{}{
			"a": "1",
		},
	}
	if err := cio.Track("", "test", data); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}
	if err := cio.Track("1", "", data); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "eventName")
	}

	runCases(t,
		[]testCase{
			{"1", "POST", "/api/v1/customers/1/events", body},
			{"1 ", "POST", "/api/v1/customers/1%20/events", body},
			{"1/", "POST", "/api/v1/customers/1%2F/events", body},
		},
		func(c testCase) error {
			return cio.Track(c.id, "test", data)
		})
}

func TestTrackAnonymous(t *testing.T) {
	data := map[string]interface{}{
		"a": "1",
	}

	body := map[string]interface{}{
		"name": "test",
		"data": map[string]interface{}{
			"a": "1",
		},
	}

	expect("POST", "/api/v1/events", body)
	if err := cio.TrackAnonymous("test", data); err != nil {
		t.Error(err.Error())
	}
}

func TestDelete(t *testing.T) {
	if err := cio.Delete(""); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}
	runCases(t,
		[]testCase{
			{"1", "DELETE", "/api/v1/customers/1", nil},
			{"1 ", "DELETE", "/api/v1/customers/1%20", nil},
			{"1/", "DELETE", "/api/v1/customers/1%2F", nil},
		},
		func(c testCase) error {
			return cio.Delete(c.id)
		})
}

func TestAddDevice(t *testing.T) {
	if err := cio.AddDevice("", "d1", "ios", nil); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}
	if err := cio.AddDevice("1", "", "ios", nil); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "deviceID")
	}
	if err := cio.AddDevice("1", "d1", "", nil); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "platform")
	}

	body := map[string]map[string]interface{}{
		"device": {
			"id":        "d1",
			"platform":  "ios",
			"last_used": 1606511962,
		},
	}
	runCases(t,
		[]testCase{
			{"1", "PUT", "/api/v1/customers/1/devices", body},
			{"1 ", "PUT", "/api/v1/customers/1%20/devices", body},
			{"1/", "PUT", "/api/v1/customers/1%2F/devices", body},
		},
		func(c testCase) error {
			return cio.AddDevice(c.id, "d1", "ios", map[string]interface{}{
				"last_used": 1606511962,
			})
		})
}

func TestDeleteDevice(t *testing.T) {
	if err := cio.DeleteDevice("", "d1"); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}
	if err := cio.DeleteDevice("1", ""); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "deviceID")
	}
	runCases(t,
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
				return cio.DeleteDevice("d1", c.id)
			} else {
				return cio.DeleteDevice(c.id, "d1")
			}
		})
}

func TestAddCustomersToSegment(t *testing.T) {
	ids := []string{"1", "2", "3"}
	if err := cio.AddCustomersToSegment("", ids); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "segmentID")
	}
	if err := cio.AddCustomersToSegment("1", nil); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "customerIDs")
	}
	body := map[string]interface{}{
		"ids": ids,
	}
	runCases(t,
		[]testCase{
			{"1", "POST", "/api/v1/segments/1/add_customers", body},
			{"1 ", "POST", "/api/v1/segments/1%20/add_customers", nil},
			{"1/", "POST", "/api/v1/segments/1%2F/add_customers", nil},
		},
		func(c testCase) error {
			return cio.AddCustomersToSegment(c.id, ids)
		})
}

func TestRemoveCustomersFromSegment(t *testing.T) {
	ids := []string{"1", "2", "3"}
	if err := cio.RemoveCustomersFromSegment("", ids); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "segmentID")
	}
	if err := cio.RemoveCustomersFromSegment("1", nil); err == nil {
		t.Error("expected error")
	} else {
		checkParamError(t, err, "customerIDs")
	}
	body := map[string]interface{}{
		"ids": ids,
	}
	runCases(t,
		[]testCase{
			{"1", "POST", "/api/v1/segments/1/remove_customers", body},
			{"1 ", "POST", "/api/v1/segments/1%20/remove_customers", nil},
			{"1/", "POST", "/api/v1/segments/1%2F/remove_customers", nil},
		},
		func(c testCase) error {
			return cio.RemoveCustomersFromSegment(c.id, ids)
		})
}

var (
	expectedMethod string
	expectedPath   string
	expectedBody   interface{}
)

func handler(w http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	s := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	pair := strings.SplitN(string(decoded), ":", 2)
	if len(pair) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if pair[0] != "siteid" && pair[1] != "apikey" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if req.Method != "DELETE" && req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "expected Content-Type application/json", http.StatusBadRequest)
	}

	var data map[string]interface{}
	if len(b) > 0 {
		dec := json.NewDecoder(bytes.NewReader(b))
		dec.UseNumber()
		if err := dec.Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	validate := func(method, path string, body interface{}) error {
		if method != expectedMethod {
			return fmt.Errorf("expected %s got %s", expectedMethod, method)
		}
		if path != expectedPath {
			return fmt.Errorf("expected %s got %s", expectedPath, path)
		}
		expected, err := json.Marshal(body)
		if err != nil {
			return err
		}
		got, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if bytes.Compare(expected, got) != 0 {
			return fmt.Errorf("expected %v got %v", expected, got)
		}
		return nil
	}
	if err := validate(req.Method, req.RequestURI, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func expect(method, path string, body interface{}) {
	expectedMethod = method
	expectedPath = path
	expectedBody = body
}
