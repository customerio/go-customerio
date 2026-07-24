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
