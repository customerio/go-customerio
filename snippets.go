package customerio

import "context"

// Snippet is a reusable named liquid fragment, keyed by Name rather than a
// numeric id.
type Snippet struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}

// GetSnippetsResponse is the decoded shape of GET /v1/snippets.
type GetSnippetsResponse struct {
	Snippets []Snippet `json:"snippets"`
}

// GetSnippets returns every snippet defined in the workspace.
// See https://docs.customer.io/api/app/#operation/getSnippets
func (c *APIClient) GetSnippets(ctx context.Context) (*GetSnippetsResponse, error) {
	var resp GetSnippetsResponse
	if err := c.doJSON(ctx, "GET", "/v1/snippets", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SnippetResponse wraps a single Snippet.
type SnippetResponse struct {
	Snippet Snippet `json:"snippet"`
}

// SnippetInput creates or updates a snippet.
type SnippetInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CreateSnippet creates a snippet.
// See https://docs.customer.io/api/app/#operation/createSnippet
func (c *APIClient) CreateSnippet(ctx context.Context, input SnippetInput) (*SnippetResponse, error) {
	if input.Name == "" {
		return nil, ParamError{Param: "name"}
	}
	if input.Value == "" {
		return nil, ParamError{Param: "value"}
	}

	var resp SnippetResponse
	if err := c.doJSON(ctx, "POST", "/v1/snippets", input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateSnippet creates or updates a snippet (an upsert, keyed by input.Name).
// See https://docs.customer.io/api/app/#operation/updateSnippet
func (c *APIClient) UpdateSnippet(ctx context.Context, input SnippetInput) (*SnippetResponse, error) {
	if input.Name == "" {
		return nil, ParamError{Param: "name"}
	}
	if input.Value == "" {
		return nil, ParamError{Param: "value"}
	}

	var resp SnippetResponse
	if err := c.doJSON(ctx, "PUT", "/v1/snippets", input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteSnippet deletes a snippet by name.
// See https://docs.customer.io/api/app/#operation/deleteSnippet
func (c *APIClient) DeleteSnippet(ctx context.Context, name string) error {
	if name == "" {
		return ParamError{Param: "name"}
	}

	return c.doJSON(ctx, "DELETE", formatPath("/v1/snippets/%s", name), nil, nil, 200, 204)
}
