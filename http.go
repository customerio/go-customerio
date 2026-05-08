package customerio

import (
	"net/http"
	"time"
)

// DefaultHTTPTimeout is the timeout used by clients created without
// WithHTTPClient. For requests with large attachments, you may need
// a longer timeout via WithHTTPClient.
const DefaultHTTPTimeout = 30 * time.Second

func newDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   DefaultHTTPTimeout,
		Transport: newDefaultTransport(),
	}
}

func newDefaultTransport() http.RoundTripper {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		// Non-*http.Transport (e.g. OpenTelemetry instrumented transport).
		// Return as-is; MaxIdleConnsPerHost tuning is skipped.
		return http.DefaultTransport
	}

	transport = transport.Clone()
	transport.MaxIdleConnsPerHost = 100
	return transport
}
