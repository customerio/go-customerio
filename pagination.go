package customerio

// PaginationOptions controls paging on App API list/search endpoints.
// Start is an opaque cursor returned by a previous response's "next" field,
// not a page number or offset — pass it back verbatim to fetch the next
// page. Limit caps the number of results per page; 0 lets the API apply its
// default.
type PaginationOptions struct {
	Start string
	Limit int
}

func (o PaginationOptions) apply(q *queryBuilder) *queryBuilder {
	return q.setString("start", o.Start).setInt("limit", o.Limit)
}

// ObjectIdentifierType identifies which kind of id a request supplies for an
// object (as opposed to IdentifierType, which is for customers).
type ObjectIdentifierType string

const (
	// ObjectIdentifierTypeObjectID selects the caller-supplied object_id.
	ObjectIdentifierTypeObjectID ObjectIdentifierType = "object_id"
	// ObjectIdentifierTypeCioObjectID selects Customer.io's generated cio_object_id.
	ObjectIdentifierTypeCioObjectID ObjectIdentifierType = "cio_object_id"
)

// MetricType filters metrics/message reports to one delivery channel.
type MetricType string

const (
	MetricTypeEmail    MetricType = "email"
	MetricTypeWebhook  MetricType = "webhook"
	MetricTypeTwilio   MetricType = "twilio"
	MetricTypeWhatsApp MetricType = "whatsapp"
	MetricTypeSlack    MetricType = "slack"
	MetricTypePush     MetricType = "push"
	MetricTypeInApp    MetricType = "in_app"
)

// NewsletterChannelType filters newsletter message reports to one delivery
// channel — the same idea as MetricType, but newsletters support "inbox"
// and not "whatsapp"/"slack", so it's a distinct type rather than a
// restriction of MetricType's values.
type NewsletterChannelType string

const (
	NewsletterChannelEmail   NewsletterChannelType = "email"
	NewsletterChannelWebhook NewsletterChannelType = "webhook"
	NewsletterChannelTwilio  NewsletterChannelType = "twilio"
	NewsletterChannelPush    NewsletterChannelType = "push"
	NewsletterChannelInApp   NewsletterChannelType = "in_app"
	NewsletterChannelInbox   NewsletterChannelType = "inbox"
)

// MetricResolution buckets a metrics or journey-metrics report. Which of
// the short (hours/days/weeks/months) or long (hourly/daily/weekly/
// monthly) form an endpoint accepts is validated server-side, not here.
type MetricResolution string

// CampaignMetricsVersion selects the metrics schema for campaign/action
// metrics: "2" adds bot-click-filtered human_*/prefetch_* series on top of "1".
type CampaignMetricsVersion string

const (
	CampaignMetricsVersion1 CampaignMetricsVersion = "1"
	CampaignMetricsVersion2 CampaignMetricsVersion = "2"
)

// SortDirection orders a list endpoint's results.
type SortDirection string

const (
	SortAscending  SortDirection = "asc"
	SortDescending SortDirection = "desc"
)

// DesignStudioFilter is a tri-state filter used by Design Studio list
// endpoints (is_template/has_translations/is_linked): DesignStudioFilterAny
// is sent explicitly (unlike the zero value "", which is omitted from the
// query entirely) to request no filtering — both have the same effect
// server-side (the API defaults to "any"), but explicit "any" round-trips
// onto the wire while "" doesn't send the parameter at all.
type DesignStudioFilter string

const (
	DesignStudioFilterTrue  DesignStudioFilter = "true"
	DesignStudioFilterFalse DesignStudioFilter = "false"
	DesignStudioFilterAny   DesignStudioFilter = "any"
)

// DesignStudioSortBy selects the sort field for Design Studio list endpoints.
type DesignStudioSortBy string

const (
	DesignStudioSortByCreated DesignStudioSortBy = "created"
	DesignStudioSortByUpdated DesignStudioSortBy = "updated"
	DesignStudioSortByName    DesignStudioSortBy = "name"
)

// DesignStudioListOptions filters/sorts/paginates the Design Studio folder,
// email, and component list endpoints. Field names are Go-idiomatic
// camelCase but map to snake_case wire query params (e.g. ParentFolderID ->
// parent_folder_id) — a real mismatch in the underlying API, not
// inconsistency introduced here.
type DesignStudioListOptions struct {
	ParentFolderID        string
	DirectDescendantsOnly *bool
	SortBy                DesignStudioSortBy
	SortOrder             SortDirection
	CreatedBefore         int64
	CreatedAfter          int64
	UpdatedBefore         int64
	UpdatedAfter          int64
	Page                  int
	Limit                 int
}

func (o DesignStudioListOptions) apply(q *queryBuilder) *queryBuilder {
	return q.setString("parent_folder_id", o.ParentFolderID).
		setBool("direct_descendants_only", o.DirectDescendantsOnly).
		setString("sort_by", string(o.SortBy)).
		setString("sort_order", string(o.SortOrder)).
		setInt64("created_before", o.CreatedBefore).
		setInt64("created_after", o.CreatedAfter).
		setInt64("updated_before", o.UpdatedBefore).
		setInt64("updated_after", o.UpdatedAfter).
		setInt("page", o.Page).
		setInt("limit", o.Limit)
}
