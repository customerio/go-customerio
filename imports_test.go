package customerio_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

// importServer answers every request with the given status and body, after
// handing the request and its body to verify.
func importServer(t *testing.T, status int, response string, verify func(req *http.Request, body []byte)) (*customerio.APIClient, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
		}
		defer func() { _ = req.Body.Close() }()
		if verify != nil {
			verify(req, b)
		}
		w.WriteHeader(status)
		// errcheck (golangci-lint): explicitly discard return values
		_, _ = w.Write([]byte(response))
	}))

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	return api, srv
}

// importServerNeverCalled fails the test if the client sends any request.
func importServerNeverCalled(t *testing.T) (*customerio.APIClient, *httptest.Server) {
	t.Helper()
	return importServer(t, http.StatusOK, `{"import":{"id":"1"}}`, func(req *http.Request, _ []byte) {
		t.Errorf("server must not be called, got %s %s", req.Method, req.URL.Path)
	})
}

func TestCreateImportDevice(t *testing.T) {
	api, srv := importServer(t, http.StatusOK, `{"import":{"id":"23","type":"device","state":"preprocessing"}}`, func(req *http.Request, body []byte) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		if req.URL.Path != "/v1/imports" {
			t.Errorf("expected /v1/imports, got %s", req.URL.Path)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer myKey" {
			t.Errorf("expected Authorization %q, got %q", "Bearer myKey", got)
		}

		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}

		// Fields left empty (description, people_to_process, data_to_process,
		// object_type_id) must be absent, so the whole body is compared.
		want := map[string]any{
			"import": map[string]any{
				"name":          "Device tokens",
				"data_file_url": "https://example.com/devices.csv",
				"type":          "device",
				"identifier":    "id",
			},
		}
		if !reflect.DeepEqual(payload, want) {
			t.Errorf("body mismatch: want %#v got %#v", want, payload)
		}
	})
	defer srv.Close()

	imp, err := api.CreateImport(context.Background(), &customerio.CreateImportRequest{
		Name:        "Device tokens",
		DataFileURL: "https://example.com/devices.csv",
		Type:        customerio.ImportTypeDevice,
		Identifier:  "id",
	})
	if err != nil {
		t.Fatal(err)
	}
	if imp.ID != 23 {
		t.Errorf("expected ID 23, got %d", imp.ID)
	}
	if imp.Type != customerio.ImportTypeDevice {
		t.Errorf("expected type device, got %q", imp.Type)
	}
	if imp.State != "preprocessing" {
		t.Errorf("expected state preprocessing, got %q", imp.State)
	}
}

func TestCreateImportOptionalFields(t *testing.T) {
	api, srv := importServer(t, http.StatusOK, `{"import":{"id":"24","type":"object","state":"preprocessing","object_type_id":"2"}}`, func(_ *http.Request, body []byte) {
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}

		want := map[string]any{
			"import": map[string]any{
				"name":            "Companies",
				"data_file_url":   "https://example.com/companies.csv",
				"type":            "object",
				"description":     "Monthly sync",
				"data_to_process": "only_new",
				"object_type_id":  "2",
			},
		}
		if !reflect.DeepEqual(payload, want) {
			t.Errorf("body mismatch: want %#v got %#v", want, payload)
		}
	})
	defer srv.Close()

	imp, err := api.CreateImport(context.Background(), &customerio.CreateImportRequest{
		Name:          "Companies",
		DataFileURL:   "https://example.com/companies.csv",
		Type:          customerio.ImportTypeObject,
		Description:   "Monthly sync",
		DataToProcess: "only_new",
		ObjectTypeID:  "2",
	})
	if err != nil {
		t.Fatal(err)
	}
	if imp.ObjectTypeID != "2" {
		t.Errorf("expected object type id 2, got %q", imp.ObjectTypeID)
	}
}

func TestCreateImportNilRequest(t *testing.T) {
	api, srv := importServerNeverCalled(t)
	defer srv.Close()

	_, err := api.CreateImport(context.Background(), nil)
	checkParamError(t, err, "req")
}

func TestCreateImportMissingDataFileURL(t *testing.T) {
	api, srv := importServerNeverCalled(t)
	defer srv.Close()

	_, err := api.CreateImport(context.Background(), &customerio.CreateImportRequest{
		Name: "Device tokens",
		Type: customerio.ImportTypeDevice,
	})
	checkParamError(t, err, "DataFileURL")
}

func TestCreateImportMissingType(t *testing.T) {
	api, srv := importServerNeverCalled(t)
	defer srv.Close()

	_, err := api.CreateImport(context.Background(), &customerio.CreateImportRequest{
		Name:        "Device tokens",
		DataFileURL: "https://example.com/devices.csv",
	})
	checkParamError(t, err, "Type")
}

func TestCreateImportError(t *testing.T) {
	api, srv := importServer(t, http.StatusUnprocessableEntity, `{"errors":[{"detail":"identifier is required","source":{"pointer":"/data/attributes/identifier"},"status":"422"}]}`, nil)
	defer srv.Close()

	_, err := api.CreateImport(context.Background(), &customerio.CreateImportRequest{
		Name:        "Device tokens",
		DataFileURL: "https://example.com/devices.csv",
		Type:        customerio.ImportTypeDevice,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	cioErr, ok := err.(*customerio.CustomerIOError)
	if !ok {
		t.Fatalf("expected *CustomerIOError, got %T", err)
	}
	if cioErr.StatusCode() != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, cioErr.StatusCode())
	}
	if !strings.Contains(string(cioErr.Body()), "identifier is required") {
		t.Errorf("expected detail in Body, got %q", string(cioErr.Body()))
	}
}

func TestGetImport(t *testing.T) {
	response := `{"import":{"id":"23","type":"device","state":"imported","name":"Device tokens","description":"","error":"","created_at":1764200000,"updated_at":1764200060,"identifier":"id","rows_to_import":10,"rows_imported":8,"error_download_url":"https://api.customer.io/v1/exports/7/download"}}`
	api, srv := importServer(t, http.StatusOK, response, func(req *http.Request, body []byte) {
		if req.Method != "GET" {
			t.Errorf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/imports/23" {
			t.Errorf("expected /v1/imports/23, got %s", req.URL.Path)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer myKey" {
			t.Errorf("expected Authorization %q, got %q", "Bearer myKey", got)
		}
		if len(body) != 0 {
			t.Errorf("expected empty body, got %q", string(body))
		}
	})
	defer srv.Close()

	imp, err := api.GetImport(context.Background(), 23)
	if err != nil {
		t.Fatal(err)
	}
	if imp.ID != 23 {
		t.Errorf("expected ID 23, got %d", imp.ID)
	}
	if imp.State != "imported" {
		t.Errorf("expected state imported, got %q", imp.State)
	}
	if imp.CreatedAt != 1764200000 || imp.UpdatedAt != 1764200060 {
		t.Errorf("unexpected timestamps: %d %d", imp.CreatedAt, imp.UpdatedAt)
	}
	if imp.ErrorDownloadURL != "https://api.customer.io/v1/exports/7/download" {
		t.Errorf("unexpected error_download_url: %q", imp.ErrorDownloadURL)
	}
	if imp.RowsToImport == nil || *imp.RowsToImport != 10 {
		t.Errorf("expected rows_to_import 10, got %v", imp.RowsToImport)
	}
	if imp.RowsImported == nil || *imp.RowsImported != 8 {
		t.Errorf("expected rows_imported 8, got %v", imp.RowsImported)
	}
}

func TestGetImportIDZero(t *testing.T) {
	api, srv := importServerNeverCalled(t)
	defer srv.Close()

	_, err := api.GetImport(context.Background(), 0)
	checkParamError(t, err, "importID")
}

func TestGetImportIDNegative(t *testing.T) {
	api, srv := importServerNeverCalled(t)
	defer srv.Close()

	_, err := api.GetImport(context.Background(), -1)
	checkParamError(t, err, "importID")
}

func TestGetImportError(t *testing.T) {
	api, srv := importServer(t, http.StatusNotFound, `{"errors":[{"detail":"import not found","status":"404"}]}`, nil)
	defer srv.Close()

	_, err := api.GetImport(context.Background(), 23)
	cioErr, ok := err.(*customerio.CustomerIOError)
	if !ok {
		t.Fatalf("expected *CustomerIOError, got %T", err)
	}
	if cioErr.StatusCode() != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, cioErr.StatusCode())
	}
}
