package customerio

import "context"

// ReportingWebhook notifies an external endpoint of message-lifecycle events.
type ReportingWebhook struct {
	ID             int      `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Type           string   `json:"type,omitempty"`
	Endpoint       string   `json:"endpoint,omitempty"`
	Disabled       bool     `json:"disabled,omitempty"`
	FullResolution bool     `json:"full_resolution,omitempty"`
	WithContent    bool     `json:"with_content,omitempty"`
	Events         []string `json:"events,omitempty"`
}

// ListReportingWebhooksResponse is the decoded shape of GET /v1/reporting_webhooks.
// The wrapping key isn't shown in an example in the OpenAPI doc; inferred
// by convention with every other list endpoint in this API (plural of the
// URL's resource segment).
type ListReportingWebhooksResponse struct {
	ReportingWebhooks []ReportingWebhook `json:"reporting_webhooks"`
}

// ListReportingWebhooks returns every reporting webhook configured in the workspace.
// See https://docs.customer.io/api/app/#operation/listReportingWebhooks
func (c *APIClient) ListReportingWebhooks(ctx context.Context) (*ListReportingWebhooksResponse, error) {
	var resp ListReportingWebhooksResponse
	if err := c.doJSON(ctx, "GET", "/v1/reporting_webhooks", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReportingWebhookResponse wraps a single ReportingWebhook.
type ReportingWebhookResponse struct {
	ReportingWebhook ReportingWebhook `json:"reporting_webhook"`
}

// ReportingWebhookInput creates a reporting webhook. The App API's schema
// documents Name/Events as required alongside Endpoint, but the Node SDK
// this is ported from only validates Endpoint client-side — this client
// mirrors that (server-side validation still applies to the others).
type ReportingWebhookInput struct {
	Endpoint       string   `json:"endpoint"`
	Events         []string `json:"events,omitempty"`
	Name           string   `json:"name,omitempty"`
	FullResolution *bool    `json:"full_resolution,omitempty"`
	WithContent    *bool    `json:"with_content,omitempty"`
	Disabled       *bool    `json:"disabled,omitempty"`
}

// CreateReportingWebhook creates a reporting webhook.
// See https://docs.customer.io/api/app/#operation/createReportingWebhook
func (c *APIClient) CreateReportingWebhook(ctx context.Context, input ReportingWebhookInput) (*ReportingWebhookResponse, error) {
	if input.Endpoint == "" {
		return nil, ParamError{Param: "endpoint"}
	}

	var resp ReportingWebhookResponse
	if err := c.doJSON(ctx, "POST", "/v1/reporting_webhooks", input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetReportingWebhook returns one reporting webhook by id.
// See https://docs.customer.io/api/app/#operation/getReportingWebhook
func (c *APIClient) GetReportingWebhook(ctx context.Context, webhookID string) (*ReportingWebhookResponse, error) {
	if webhookID == "" {
		return nil, ParamError{Param: "webhookID"}
	}

	var resp ReportingWebhookResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/reporting_webhooks/%s", webhookID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateReportingWebhook updates a reporting webhook. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateReportingWebhook
func (c *APIClient) UpdateReportingWebhook(ctx context.Context, webhookID string, data map[string]any) (*ReportingWebhookResponse, error) {
	if webhookID == "" {
		return nil, ParamError{Param: "webhookID"}
	}

	var resp ReportingWebhookResponse
	if err := c.doJSON(ctx, "PUT", formatPath("/v1/reporting_webhooks/%s", webhookID), data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteReportingWebhook deletes a reporting webhook.
// See https://docs.customer.io/api/app/#operation/deleteReportingWebhook
func (c *APIClient) DeleteReportingWebhook(ctx context.Context, webhookID string) error {
	if webhookID == "" {
		return ParamError{Param: "webhookID"}
	}

	return c.doJSON(ctx, "DELETE", formatPath("/v1/reporting_webhooks/%s", webhookID), nil, nil, 200, 204)
}
