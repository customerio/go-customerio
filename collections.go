package customerio

import "context"

// Collection is a workspace-level reference dataset (loaded from inline
// data or a source URL) usable in message liquid via lookups. The read
// object exposes Schema/Rows/Bytes rather than the original Data/URL —
// those are write-only concepts on the create/update request.
type Collection struct {
	ID        int      `json:"id"`
	Name      string   `json:"name,omitempty"`
	Bytes     int64    `json:"bytes,omitempty"`
	Rows      int      `json:"rows,omitempty"`
	Schema    []string `json:"schema,omitempty"`
	CreatedAt int64    `json:"created_at,omitempty"`
	UpdatedAt int64    `json:"updated_at,omitempty"`
}

// ListCollectionsResponse is the decoded shape of GET /v1/collections.
type ListCollectionsResponse struct {
	Collections []Collection `json:"collections"`
}

// ListCollections returns every collection defined in the workspace.
// See https://docs.customer.io/api/app/#operation/listCollections
func (c *APIClient) ListCollections(ctx context.Context) (*ListCollectionsResponse, error) {
	var resp ListCollectionsResponse
	if err := c.doJSON(ctx, "GET", "/v1/collections", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CollectionResponse wraps a single Collection.
type CollectionResponse struct {
	Collection Collection `json:"collection"`
}

// CollectionInput creates a collection from either inline Data or a source
// URL — the App API's create request is documented as accepting exactly
// one of the two (oneOf {name,data} / {name,url}); like the Node SDK, this
// client doesn't enforce that locally, so set exactly one of Data/URL
// yourself.
type CollectionInput struct {
	Name string           `json:"name"`
	Data []map[string]any `json:"data,omitempty"`
	URL  string           `json:"url,omitempty"`
}

// CreateCollection creates a collection.
// See https://docs.customer.io/api/app/#operation/createCollection
func (c *APIClient) CreateCollection(ctx context.Context, input CollectionInput) (*CollectionResponse, error) {
	if input.Name == "" {
		return nil, ParamError{Param: "name"}
	}

	var resp CollectionResponse
	if err := c.doJSON(ctx, "POST", "/v1/collections", input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCollection returns one collection by id.
// See https://docs.customer.io/api/app/#operation/getCollection
func (c *APIClient) GetCollection(ctx context.Context, collectionID string) (*CollectionResponse, error) {
	if collectionID == "" {
		return nil, ParamError{Param: "collectionID"}
	}

	var resp CollectionResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/collections/%s", collectionID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CollectionUpdate updates a collection's name and/or backing data/url — all
// fields optional, set whichever you want to change.
type CollectionUpdate struct {
	Name string           `json:"name,omitempty"`
	Data []map[string]any `json:"data,omitempty"`
	URL  string           `json:"url,omitempty"`
}

// UpdateCollection updates a collection.
// See https://docs.customer.io/api/app/#operation/updateCollection
func (c *APIClient) UpdateCollection(ctx context.Context, collectionID string, updates CollectionUpdate) (*CollectionResponse, error) {
	if collectionID == "" {
		return nil, ParamError{Param: "collectionID"}
	}

	var resp CollectionResponse
	if err := c.doJSON(ctx, "PUT", formatPath("/v1/collections/%s", collectionID), updates, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCollection deletes a collection.
// See https://docs.customer.io/api/app/#operation/deleteCollection
func (c *APIClient) DeleteCollection(ctx context.Context, collectionID string) error {
	if collectionID == "" {
		return ParamError{Param: "collectionID"}
	}

	return c.doJSON(ctx, "DELETE", formatPath("/v1/collections/%s", collectionID), nil, nil, 200, 204)
}

// GetCollectionContent returns a collection's current content. The
// published OpenAPI schema for this endpoint documents a single free-form
// JSON object, which is inconsistent with UpdateCollectionContent's actual
// wire behavior (see its doc comment) — decoded generically here since
// neither this client nor the Node SDK it mirrors can resolve that
// discrepancy without a live response to check against.
// See https://docs.customer.io/api/app/#operation/getCollectionContent
func (c *APIClient) GetCollectionContent(ctx context.Context, collectionID string) (map[string]any, error) {
	if collectionID == "" {
		return nil, ParamError{Param: "collectionID"}
	}

	var resp map[string]any
	if err := c.doJSON(ctx, "GET", formatPath("/v1/collections/%s/content", collectionID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateCollectionContent replaces a collection's entire content with rows.
// This sends a top-level JSON array, matching the Node SDK's actual request
// behavior (createCollection's Data field is also an array of rows) — even
// though the App API's published OpenAPI schema for this endpoint currently
// documents a single flat object instead of an array. Verify against a
// live response if this matters for your use case.
// See https://docs.customer.io/api/app/#operation/updateCollectionContent
func (c *APIClient) UpdateCollectionContent(ctx context.Context, collectionID string, rows []map[string]any) error {
	if collectionID == "" {
		return ParamError{Param: "collectionID"}
	}
	if rows == nil {
		return ParamError{Param: "rows"}
	}

	return c.doJSON(ctx, "PUT", formatPath("/v1/collections/%s/content", collectionID), rows, nil, 200, 204)
}
