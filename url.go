package customerio

import (
	"fmt"
	"net/url"
	"strconv"
)

// formatPath builds a request path from a printf-style format string, escaping
// any string argument with url.PathEscape so dynamic values (customer IDs,
// device IDs) are safe to interpolate without pre-escaping; a "/" inside such a
// value is encoded as %2F rather than being treated as a path separator.
// Non-string arguments (e.g. integer IDs) are passed through unchanged.
//
// The returned path does not include the base URL; callers prepend it.
func formatPath(format string, args ...any) string {
	for i, a := range args {
		if s, ok := a.(string); ok {
			args[i] = url.PathEscape(s)
		}
	}
	return fmt.Sprintf(format, args...)
}

// queryBuilder accumulates optional query parameters, skipping any left at
// their zero value — mirroring the App API's own buildQueryString helper
// (see customerio-node lib/utils.ts), which omits null/undefined/empty
// values rather than sending them as empty query params.
type queryBuilder struct {
	values url.Values
}

func newQuery() *queryBuilder {
	return &queryBuilder{values: url.Values{}}
}

func (q *queryBuilder) setString(key, value string) *queryBuilder {
	if value != "" {
		q.values.Set(key, value)
	}
	return q
}

func (q *queryBuilder) setInt(key string, value int) *queryBuilder {
	if value != 0 {
		q.values.Set(key, strconv.Itoa(value))
	}
	return q
}

func (q *queryBuilder) setBool(key string, value *bool) *queryBuilder {
	if value != nil {
		q.values.Set(key, strconv.FormatBool(*value))
	}
	return q
}

func (q *queryBuilder) setInt64(key string, value int64) *queryBuilder {
	if value != 0 {
		q.values.Set(key, strconv.FormatInt(value, 10))
	}
	return q
}

// String renders the accumulated parameters as a "?"-prefixed query string,
// or "" if no parameter was ever set.
func (q *queryBuilder) String() string {
	encoded := q.values.Encode()
	if encoded == "" {
		return ""
	}
	return "?" + encoded
}
