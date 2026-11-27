package customerio

import (
	"context"
	"encoding/json"
	"net/http"
)

// ImportType is the kind of data a CSV import holds. The API validates the
// value, so a type added to the API later can be passed as ImportType("name").
type ImportType string

const (
	ImportTypePeople       ImportType = "people"
	ImportTypeEvent        ImportType = "event"
	ImportTypeObject       ImportType = "object"
	ImportTypeRelationship ImportType = "relationship"
	ImportTypeDevice       ImportType = "device"
)

// CreateImportRequest describes a CSV import to start with CreateImport.
type CreateImportRequest struct {
	// Name is a required label for the import.
	Name string `json:"name"`
	// DataFileURL is the required public URL of the CSV file.
	DataFileURL string `json:"data_file_url"`
	// Type is the required kind of data in the file.
	Type ImportType `json:"type"`
	// Identifier is the column that identifies a person: "id" or "email"
	// ("cio_id" also works for relationship imports). Required for every type
	// except object.
	Identifier string `json:"identifier,omitempty"`
	// Description is an optional note about the import.
	Description string `json:"description,omitempty"`
	// PeopleToProcess limits the rows to import: "all", "only_existing" or
	// "only_new". Set PeopleToProcess or DataToProcess, never both. Device
	// imports accept neither.
	PeopleToProcess string `json:"people_to_process,omitempty"`
	// DataToProcess limits the rows to import: "all", "only_existing" or
	// "only_new". Set PeopleToProcess or DataToProcess, never both. Device
	// imports accept neither.
	DataToProcess string `json:"data_to_process,omitempty"`
	// ObjectTypeID is the object type to import into. Object imports only.
	ObjectTypeID string `json:"object_type_id,omitempty"`
}

// Import is the state of a CSV import, returned by CreateImport and GetImport.
type Import struct {
	ID          int        `json:"id,string"`
	Type        ImportType `json:"type"`
	State       string     `json:"state"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Error       string     `json:"error"`
	// CreatedAt and UpdatedAt are unix timestamps in seconds.
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
	Identifier      string `json:"identifier,omitempty"`
	PeopleToProcess string `json:"people_to_process,omitempty"`
	DataToProcess   string `json:"data_to_process,omitempty"`
	ObjectTypeID    string `json:"object_type_id,omitempty"`
	RowsToImport    *int   `json:"rows_to_import,omitempty"`
	RowsImported    *int   `json:"rows_imported,omitempty"`
	// ErrorDownloadURL points to a file listing the rows that were skipped and
	// why. It can be set even when the import succeeds.
	ErrorDownloadURL   string `json:"error_download_url,omitempty"`
	WarningDownloadURL string `json:"warning_download_url,omitempty"`
}

// importRequestEnvelope is the wire shape for POST /v1/imports.
type importRequestEnvelope struct {
	Import *CreateImportRequest `json:"import"`
}

// importResponseEnvelope is the wire shape both import calls answer with.
type importResponseEnvelope struct {
	Import Import `json:"import"`
}

// CreateImport starts a CSV import by POSTing to /v1/imports. The import runs
// in the background; use GetImport with the returned ID to follow it.
func (c *APIClient) CreateImport(ctx context.Context, req *CreateImportRequest) (*Import, error) {
	if req == nil {
		return nil, ParamError{Param: "req"}
	}
	if req.DataFileURL == "" {
		return nil, ParamError{Param: "DataFileURL"}
	}
	if req.Type == "" {
		return nil, ParamError{Param: "Type"}
	}

	return c.doImportRequest(ctx, "POST", "/v1/imports", importRequestEnvelope{Import: req})
}

// GetImport reads the state of an import by sending a GET to /v1/imports/{id}.
func (c *APIClient) GetImport(ctx context.Context, importID int) (*Import, error) {
	if importID <= 0 {
		return nil, ParamError{Param: "importID"}
	}

	return c.doImportRequest(ctx, "GET", formatPath("/v1/imports/%d", importID), nil)
}

func (c *APIClient) doImportRequest(ctx context.Context, verb, requestPath string, payload any) (*Import, error) {
	body, statusCode, err := c.doRequest(ctx, verb, requestPath, payload)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, &CustomerIOError{
			status: statusCode,
			url:    c.URL + requestPath,
			body:   body,
		}
	}

	var resp importResponseEnvelope
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp.Import, nil
}
