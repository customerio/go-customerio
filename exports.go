package customerio

import "context"

// Export reports a bulk data-export job's progress.
type Export struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id,omitempty"`
	UserEmail     string `json:"user_email,omitempty"`
	Total         int    `json:"total,omitempty"`
	DeduplicateID string `json:"deduplicate_id,omitempty"`
	Type          string `json:"type,omitempty"`
	Failed        bool   `json:"failed,omitempty"`
	Description   string `json:"description,omitempty"`
	Downloads     int    `json:"downloads,omitempty"`
	CreatedAt     int64  `json:"created_at,omitempty"`
	UpdatedAt     int64  `json:"updated_at,omitempty"`
	Status        string `json:"status,omitempty"`
}

// ListExportsResponse is the decoded shape of GET /v1/exports.
type ListExportsResponse struct {
	Exports []Export `json:"exports"`
}

// ListExports returns every export job in the workspace.
// See https://docs.customer.io/api/app/#operation/listExports
func (c *APIClient) ListExports(ctx context.Context) (*ListExportsResponse, error) {
	var resp ListExportsResponse
	if err := c.doJSON(ctx, "GET", "/v1/exports", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExportResponse wraps a single Export.
type ExportResponse struct {
	Export Export `json:"export"`
}

// GetExport returns one export job's progress by id.
// See https://docs.customer.io/api/app/#operation/getExport
func (c *APIClient) GetExport(ctx context.Context, exportID string) (*ExportResponse, error) {
	if exportID == "" {
		return nil, ParamError{Param: "exportID"}
	}

	var resp ExportResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/exports/%s", exportID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DownloadExportResponse is the decoded shape of GET /v1/exports/{id}/download.
type DownloadExportResponse struct {
	URL string `json:"url"`
}

// DownloadExport returns a signed download URL for a completed export job.
// The URL expires after 15 minutes.
// See https://docs.customer.io/api/app/#operation/downloadExport
func (c *APIClient) DownloadExport(ctx context.Context, exportID string) (*DownloadExportResponse, error) {
	if exportID == "" {
		return nil, ParamError{Param: "exportID"}
	}

	var resp DownloadExportResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/exports/%s/download", exportID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type createCustomersExportRequest struct {
	Filters AudienceFilter `json:"filters"`
}

// CreateCustomersExport starts an export job of every customer matching filter.
// See https://docs.customer.io/api/app/#operation/createCustomersExport
func (c *APIClient) CreateCustomersExport(ctx context.Context, filters AudienceFilter) (*ExportResponse, error) {
	var resp ExportResponse
	if err := c.doJSON(ctx, "POST", "/v1/exports/customers", createCustomersExportRequest{Filters: filters}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeliveryExportMetric filters CreateDeliveriesExport to deliveries that
// reached a specific lifecycle state.
type DeliveryExportMetric string

const (
	DeliveryExportMetricCreated       DeliveryExportMetric = "created"
	DeliveryExportMetricAttempted     DeliveryExportMetric = "attempted"
	DeliveryExportMetricSent          DeliveryExportMetric = "sent"
	DeliveryExportMetricDelivered     DeliveryExportMetric = "delivered"
	DeliveryExportMetricOpened        DeliveryExportMetric = "opened"
	DeliveryExportMetricClicked       DeliveryExportMetric = "clicked"
	DeliveryExportMetricConverted     DeliveryExportMetric = "converted"
	DeliveryExportMetricBounced       DeliveryExportMetric = "bounced"
	DeliveryExportMetricSpammed       DeliveryExportMetric = "spammed"
	DeliveryExportMetricUnsubscribed  DeliveryExportMetric = "unsubscribed"
	DeliveryExportMetricDropped       DeliveryExportMetric = "dropped"
	DeliveryExportMetricFailed        DeliveryExportMetric = "failed"
	DeliveryExportMetricUndeliverable DeliveryExportMetric = "undeliverable"
)

// DeliveryExportOptions filters CreateDeliveriesExport.
type DeliveryExportOptions struct {
	Start      int64
	End        int64
	Attributes []string
	Metric     DeliveryExportMetric
	Drafts     *bool
}

type createDeliveriesExportRequest struct {
	NewsletterID int                  `json:"newsletter_id"`
	Start        int64                `json:"start,omitempty"`
	End          int64                `json:"end,omitempty"`
	Attributes   []string             `json:"attributes,omitempty"`
	Metric       DeliveryExportMetric `json:"metric,omitempty"`
	Drafts       *bool                `json:"drafts,omitempty"`
}

// CreateDeliveriesExport starts an export job of a newsletter's deliveries.
// See https://docs.customer.io/api/app/#operation/createDeliveriesExport
func (c *APIClient) CreateDeliveriesExport(ctx context.Context, newsletterID int, opts DeliveryExportOptions) (*ExportResponse, error) {
	if newsletterID == 0 {
		return nil, ParamError{Param: "newsletterID"}
	}

	req := createDeliveriesExportRequest{
		NewsletterID: newsletterID,
		Start:        opts.Start,
		End:          opts.End,
		Attributes:   opts.Attributes,
		Metric:       opts.Metric,
		Drafts:       opts.Drafts,
	}

	var resp ExportResponse
	if err := c.doJSON(ctx, "POST", "/v1/exports/deliveries", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
