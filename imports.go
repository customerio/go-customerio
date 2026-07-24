package customerio

import "context"

// ImportType is the kind of data an import job loads.
type ImportType string

const (
	ImportTypePeople       ImportType = "people"
	ImportTypeEvent        ImportType = "event"
	ImportTypeObject       ImportType = "object"
	ImportTypeRelationship ImportType = "relationship"
)

// ImportProcessScope controls which rows an import job actually applies.
type ImportProcessScope string

const (
	ImportProcessAll          ImportProcessScope = "all"
	ImportProcessOnlyExisting ImportProcessScope = "only_existing"
	ImportProcessOnlyNew      ImportProcessScope = "only_new"
)

// ImportInput starts a bulk import job from a file at DataFileURL.
type ImportInput struct {
	DataFileURL     string             `json:"data_file_url"`
	Type            ImportType         `json:"type"`
	Identifier      IdentifierType     `json:"identifier,omitempty"`
	ObjectTypeID    string             `json:"object_type_id,omitempty"`
	Name            string             `json:"name,omitempty"`
	Description     string             `json:"description,omitempty"`
	PeopleToProcess ImportProcessScope `json:"people_to_process,omitempty"`
	DataToProcess   ImportProcessScope `json:"data_to_process,omitempty"`
}

// importRequest is the POST /v1/imports body — the App API wraps the input
// under an "import" key.
type importRequest struct {
	Import ImportInput `json:"import"`
}

// Import reports a bulk import job's progress.
type Import struct {
	ID            int    `json:"id"`
	CreatedAt     int64  `json:"created_at,omitempty"`
	UpdatedAt     int64  `json:"updated_at,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	RowsToImport  int    `json:"rows_to_import,omitempty"`
	RowsImported  int    `json:"rows_imported,omitempty"`
	State         string `json:"state,omitempty"`
	Type          string `json:"type,omitempty"`
	Identifier    string `json:"identifier,omitempty"`
	DataToProcess string `json:"data_to_process,omitempty"`
	ObjectTypeID  string `json:"object_type_id,omitempty"`
	Error         string `json:"error,omitempty"`
}

// ImportResponse wraps a single Import.
type ImportResponse struct {
	Import Import `json:"import"`
}

// CreateImport starts a bulk import job.
// See https://docs.customer.io/api/app/#operation/createImport
func (c *APIClient) CreateImport(ctx context.Context, input ImportInput) (*ImportResponse, error) {
	if input.DataFileURL == "" {
		return nil, ParamError{Param: "dataFileURL"}
	}
	if input.Type == "" {
		return nil, ParamError{Param: "type"}
	}

	var resp ImportResponse
	if err := c.doJSON(ctx, "POST", "/v1/imports", importRequest{Import: input}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetImport returns one import job's progress by id.
// See https://docs.customer.io/api/app/#operation/getImport
func (c *APIClient) GetImport(ctx context.Context, importID string) (*ImportResponse, error) {
	if importID == "" {
		return nil, ParamError{Param: "importID"}
	}

	var resp ImportResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/imports/%s", importID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DataIndexAttribute describes a customer/object attribute's indexing
// configuration for BatchUpdateAttributes.
type DataIndexAttribute struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	ObjectTypeID   *int   `json:"object_type_id,omitempty"`
	IsRelationship *bool  `json:"is_relationship,omitempty"`
	EventName      string `json:"event_name,omitempty"`
	PrivacyLevel   *int   `json:"privacy_level,omitempty"`
}

type batchUpdateAttributesRequest struct {
	Attributes []DataIndexAttribute `json:"attributes"`
}

// BatchUpdateAttributes updates indexing configuration for up to 100
// attributes in one call. This endpoint returns no response body on
// success (204). See https://docs.customer.io/api/app/#operation/batchUpdateAttributes
func (c *APIClient) BatchUpdateAttributes(ctx context.Context, attributes []DataIndexAttribute) error {
	if len(attributes) == 0 {
		return ParamError{Param: "attributes"}
	}

	return c.doJSON(ctx, "POST", "/v1/data_index/attributes", batchUpdateAttributesRequest{Attributes: attributes}, nil, 200, 204)
}

// DataIndexEvent describes an event's indexing configuration for BatchUpdateEvents.
type DataIndexEvent struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type batchUpdateEventsRequest struct {
	Events []DataIndexEvent `json:"events"`
}

// BatchUpdateEvents updates indexing configuration for up to 100 events in
// one call. This endpoint returns no response body on success (204).
// See https://docs.customer.io/api/app/#operation/batchUpdateEvents
func (c *APIClient) BatchUpdateEvents(ctx context.Context, events []DataIndexEvent) error {
	if len(events) == 0 {
		return ParamError{Param: "events"}
	}

	return c.doJSON(ctx, "POST", "/v1/data_index/events", batchUpdateEventsRequest{Events: events}, nil, 200, 204)
}
