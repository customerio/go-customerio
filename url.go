package customerio

import (
	"net/url"
	"strings"
)

// buildURL joins the given path segments onto base and appends the optional
// query parameters, returning a fully-formed URL string. Each segment is
// percent-escaped, so dynamic values (customer IDs, device IDs) can be passed
// without pre-escaping; a "/" inside a segment is encoded as %2F rather than
// being treated as a path separator.
func buildURL(base string, query url.Values, segments ...string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	rawPath := strings.TrimRight(u.EscapedPath(), "/")
	path := strings.TrimRight(u.Path, "/")
	for _, s := range segments {
		rawPath += "/" + url.PathEscape(s)
		path += "/" + s
	}
	u.Path = path
	u.RawPath = rawPath

	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	return u.String(), nil
}
