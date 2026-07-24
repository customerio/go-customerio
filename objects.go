package customerio

import (
	"context"
	"encoding/json"
)

// ObjectAttributes is the object payload returned by GetObjectAttributes.
// ObjectTypeID is decoded leniently (json.Number) since the documented
// schema types it as a string here but as an integer on other object
// endpoints.
type ObjectAttributes struct {
	ObjectTypeID json.Number       `json:"object_type_id,omitempty"`
	Identifiers  ObjectIdentifiers `json:"identifiers"`
	Attributes   map[string]any    `json:"attributes,omitempty"`
	Timestamps   map[string]int64  `json:"timestamps,omitempty"`
}

// ObjectAttributesResponse is the decoded shape of
// GET /v1/objects/{objectTypeId}/{objectId}/attributes.
type ObjectAttributesResponse struct {
	Object ObjectAttributes `json:"object"`
}

// GetObjectAttributes returns an object's attributes. idType selects which
// kind of identifier objectID is; pass "" to let the API apply its default.
// See https://docs.customer.io/api/app/#operation/getObjectAttributes
func (c *APIClient) GetObjectAttributes(ctx context.Context, objectTypeID, objectID string, idType ObjectIdentifierType) (*ObjectAttributesResponse, error) {
	if objectTypeID == "" {
		return nil, ParamError{Param: "objectTypeID"}
	}
	if objectID == "" {
		return nil, ParamError{Param: "objectID"}
	}

	requestPath := formatPath("/v1/objects/%s/%s/attributes", objectTypeID, objectID) +
		newQuery().setString("id_type", string(idType)).String()

	var resp ObjectAttributesResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ObjectRelationship is an object's relationship to a customer.
// ObjectTypeID is decoded leniently (json.Number): see ObjectAttributes.
type ObjectRelationship struct {
	ObjectTypeID       json.Number         `json:"object_type_id,omitempty"`
	ObjectTypeDisabled bool                `json:"object_type_disabled,omitempty"`
	Identifiers        CustomerIdentifiers `json:"identifiers"`
}

// ObjectRelationshipsResponse is the decoded shape of
// GET /v1/objects/{objectTypeId}/{objectId}/relationships.
type ObjectRelationshipsResponse struct {
	Relationships []ObjectRelationship `json:"cio_relationships"`
	Next          string               `json:"next,omitempty"`
}

// ObjectRelationshipsOptions filters GetObjectRelationships.
type ObjectRelationshipsOptions struct {
	PaginationOptions
	IDType ObjectIdentifierType
}

// GetObjectRelationships returns the customers related to an object.
// See https://docs.customer.io/api/app/#operation/getObjectRelationships
func (c *APIClient) GetObjectRelationships(ctx context.Context, objectTypeID, objectID string, opts ObjectRelationshipsOptions) (*ObjectRelationshipsResponse, error) {
	if objectTypeID == "" {
		return nil, ParamError{Param: "objectTypeID"}
	}
	if objectID == "" {
		return nil, ParamError{Param: "objectID"}
	}

	q := opts.PaginationOptions.apply(newQuery()).setString("id_type", string(opts.IDType))
	requestPath := formatPath("/v1/objects/%s/%s/relationships", objectTypeID, objectID) + q.String()

	var resp ObjectRelationshipsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// findObjectsRequest is the POST /v1/objects body.
type findObjectsRequest struct {
	ObjectTypeID string       `json:"object_type_id"`
	Filter       ObjectFilter `json:"filter"`
}

// FindObjectsResponse is the decoded shape of POST /v1/objects.
type FindObjectsResponse struct {
	Identifiers []ObjectIdentifiers `json:"identifiers"`
	IDs         []string            `json:"ids"`
	Next        string              `json:"next,omitempty"`
}

// FindObjects finds objects of type objectTypeID matching filter (built
// with ObjectFilterByAttribute and the ObjectFilterAnd/Or/Not combinators).
// See https://docs.customer.io/api/app/#operation/findObjects
func (c *APIClient) FindObjects(ctx context.Context, objectTypeID string, filter ObjectFilter, opts PaginationOptions) (*FindObjectsResponse, error) {
	if objectTypeID == "" {
		return nil, ParamError{Param: "objectTypeID"}
	}

	requestPath := "/v1/objects" + opts.apply(newQuery()).String()

	var resp FindObjectsResponse
	req := findObjectsRequest{ObjectTypeID: objectTypeID, Filter: filter}
	if err := c.doJSON(ctx, "POST", requestPath, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ObjectType describes an object type configured in a workspace (e.g. a
// "Company" or "Concert" schema).
type ObjectType struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	SingularName string `json:"singular_name,omitempty"`
	Slug         string `json:"slug,omitempty"`
	SingularSlug string `json:"singular_slug,omitempty"`
	Enabled      bool   `json:"enabled"`
	Icon         string `json:"icon,omitempty"`
}

// ListObjectTypesResponse is the decoded shape of GET /v1/object_types.
type ListObjectTypesResponse struct {
	Types []ObjectType `json:"types"`
}

// ListObjectTypes returns every object type configured in the workspace.
// See https://docs.customer.io/api/app/#operation/listObjectTypes
func (c *APIClient) ListObjectTypes(ctx context.Context) (*ListObjectTypesResponse, error) {
	var resp ListObjectTypesResponse
	if err := c.doJSON(ctx, "GET", "/v1/object_types", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
