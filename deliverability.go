package customerio

import "context"

// SuppressionType is an ESP (email service provider) suppression category.
type SuppressionType string

const (
	SuppressionTypeBlocks        SuppressionType = "blocks"
	SuppressionTypeBounces       SuppressionType = "bounces"
	SuppressionTypeSpamReports   SuppressionType = "spam_reports"
	SuppressionTypeInvalidEmails SuppressionType = "invalid_emails"
)

// Suppression is one email address suppressed from a delivery category.
type Suppression struct {
	Created int64  `json:"created,omitempty"`
	Email   string `json:"email,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Status  string `json:"status,omitempty"`
}

// SuppressionSearchResponse is the decoded shape of
// GET /v1/esp/search_suppression/{email} and GET /v1/esp/suppression/{type}.
type SuppressionSearchResponse struct {
	Category     string        `json:"category,omitempty"`
	Suppressions []Suppression `json:"suppressions"`
}

// SearchSuppression returns every suppression category an email address is on.
// See https://docs.customer.io/api/app/#operation/searchSuppression
func (c *APIClient) SearchSuppression(ctx context.Context, email string) (*SuppressionSearchResponse, error) {
	if email == "" {
		return nil, ParamError{Param: "email"}
	}

	var resp SuppressionSearchResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/esp/search_suppression/%s", email), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SuppressionsOptions filters GetSuppressions. Pagination is offset-based
// (unlike GetDomainSuppressions, which is cursor-based) — the API doesn't
// document a "next" field for this endpoint, so callers track Offset themselves.
type SuppressionsOptions struct {
	Limit  int
	Offset int
	Email  string
	Domain string
}

// GetSuppressions returns every email address suppressed for one category,
// workspace-wide.
// See https://docs.customer.io/api/app/#operation/getSuppressions
func (c *APIClient) GetSuppressions(ctx context.Context, suppressionType SuppressionType, opts SuppressionsOptions) (*SuppressionSearchResponse, error) {
	if suppressionType == "" {
		return nil, ParamError{Param: "suppressionType"}
	}

	q := newQuery().
		setInt("limit", opts.Limit).
		setInt("offset", opts.Offset).
		setString("email", opts.Email).
		setString("domain", opts.Domain)
	requestPath := formatPath("/v1/esp/suppression/%s", string(suppressionType)) + q.String()

	var resp SuppressionSearchResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DomainSuppressionsOptions filters GetDomainSuppressions. Unlike
// GetSuppressions, pagination here is cursor-based via Start/the response's Next.
type DomainSuppressionsOptions struct {
	Limit int
	Email string
	Start string
}

// DomainSuppressionsResponse is the decoded shape of
// GET /v1/esp/domains/{domain}/suppression/{type}.
type DomainSuppressionsResponse struct {
	Category     string        `json:"category,omitempty"`
	Suppressions []Suppression `json:"suppressions"`
	Next         string        `json:"next,omitempty"`
}

// GetDomainSuppressions returns every email address suppressed for one
// category, scoped to a single sending domain.
// See https://docs.customer.io/api/app/#operation/getDomainSuppressions
func (c *APIClient) GetDomainSuppressions(ctx context.Context, domainName string, suppressionType SuppressionType, opts DomainSuppressionsOptions) (*DomainSuppressionsResponse, error) {
	if domainName == "" {
		return nil, ParamError{Param: "domainName"}
	}
	if suppressionType == "" {
		return nil, ParamError{Param: "suppressionType"}
	}

	q := newQuery().setInt("limit", opts.Limit).setString("email", opts.Email).setString("start", opts.Start)
	requestPath := formatPath("/v1/esp/domains/%s/suppression/%s", domainName, string(suppressionType)) + q.String()

	var resp DomainSuppressionsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateSuppression suppresses an email address from a delivery category.
// See https://docs.customer.io/api/app/#operation/createSuppression
func (c *APIClient) CreateSuppression(ctx context.Context, suppressionType SuppressionType, email string) error {
	if suppressionType == "" {
		return ParamError{Param: "suppressionType"}
	}
	if email == "" {
		return ParamError{Param: "email"}
	}

	requestPath := formatPath("/v1/esp/suppression/%s/%s", string(suppressionType), email)
	return c.doJSON(ctx, "POST", requestPath, nil, nil, 200, 204)
}

// DeleteSuppression removes an email address from a delivery category's suppression list.
// See https://docs.customer.io/api/app/#operation/deleteSuppression
func (c *APIClient) DeleteSuppression(ctx context.Context, suppressionType SuppressionType, email string) error {
	if suppressionType == "" {
		return ParamError{Param: "suppressionType"}
	}
	if email == "" {
		return ParamError{Param: "email"}
	}

	requestPath := formatPath("/v1/esp/suppression/%s/%s", string(suppressionType), email)
	return c.doJSON(ctx, "DELETE", requestPath, nil, nil, 200, 204)
}
