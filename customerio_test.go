package customerio_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/customerio/go-customerio/v3"
)

var cio *customerio.CustomerIO

func TestMain(m *testing.M) {
	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	cio = customerio.NewCustomerIO("siteid", "apikey")
	cio.URL = srv.URL

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
	attributes := map[string]interface{}{
		"a": "1",
	}
	err := cio.Identify("", attributes)
	checkParamError(t, err, "customerID")

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
			"a":         "1",
			"timestamp": time.Now().UTC(),
		},
	}
	err := cio.Track("", "test", data)
	checkParamError(t, err, "customerID")
	err = cio.Track("1", "", data)
	checkParamError(t, err, "eventName")

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
		"name":         "test",
		"anonymous_id": "anon123",
		"data": map[string]interface{}{
			"a": "1",
		},
	}

	expect("POST", "/api/v1/events", body)
	if err := cio.TrackAnonymous("anon123", "test", data); err != nil {
		t.Error(err.Error())
	}
}

func TestDelete(t *testing.T) {
	err := cio.Delete("")
	checkParamError(t, err, "customerID")
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
	err := cio.AddDevice("", "d1", "ios", nil)
	checkParamError(t, err, "customerID")
	err = cio.AddDevice("1", "", "ios", nil)
	checkParamError(t, err, "deviceID")
	err = cio.AddDevice("1", "d1", "", nil)
	checkParamError(t, err, "platform")

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
	err := cio.DeleteDevice("", "d1")
	checkParamError(t, err, "customerID")

	err = cio.DeleteDevice("1", "")
	checkParamError(t, err, "deviceID")

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

func TestMergeCustomers(t *testing.T) {
	err1 := cio.MergeCustomers(customerio.Identifier{
		Type:  "",
		Value: "id1",
	}, customerio.Identifier{
		Type:  "id",
		Value: "id2",
	})
	checkParamError(t, err1, "primary")

	err2 := cio.MergeCustomers(customerio.Identifier{
		Type:  "id",
		Value: "",
	}, customerio.Identifier{
		Type:  "id",
		Value: "id2",
	})
	checkParamError(t, err2, "primary")

	err3 := cio.MergeCustomers(customerio.Identifier{
		Type:  "email",
		Value: "id1",
	}, customerio.Identifier{
		Type:  "",
		Value: "id2",
	})
	checkParamError(t, err3, "secondary")

	err4 := cio.MergeCustomers(customerio.Identifier{
		Type:  "cio_id",
		Value: "id1",
	}, customerio.Identifier{
		Type:  "email",
		Value: "",
	})
	checkParamError(t, err4, "secondary")

	runCases(t,
		[]testCase{
			{"1", "POST", "/api/v1/merge_customers", `{"primary":{"email":"cool.person@company.com"},"secondary":{"email":"cperson@gmail.com"}}`},
			{"2", "POST", "/api/v1/merge_customers", `{"primary":{"id":"cool.person@company.com"},"secondary":{"cio_id":"person2"}}`},
			{"3", "POST", "/api/v1/merge_customers", `{"primary":{"cio_id":"CIO123"},"secondary":{"id":"person1"}}`},
		},
		func(c testCase) error {
			if c.id == "1" {
				return cio.MergeCustomers(customerio.Identifier{
					Type:  "email",
					Value: "cool.person@company.com",
				}, customerio.Identifier{
					Type:  "email",
					Value: "cperson@gmail.com",
				})
			} else if c.id == "2" {
				return cio.MergeCustomers(customerio.Identifier{
					Type:  "id",
					Value: "cool.person@company.com",
				}, customerio.Identifier{
					Type:  "cio_id",
					Value: "person2",
				})
			} else {
				return cio.MergeCustomers(customerio.Identifier{
					Type:  customerio.IdentifierTypeCioID,
					Value: "CIO123",
				}, customerio.Identifier{
					Type:  customerio.IdentifierTypeID,
					Value: "person1",
				})
			}
		})
}
