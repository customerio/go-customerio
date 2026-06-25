package customerio

import (
	"fmt"
	"net/url"
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
