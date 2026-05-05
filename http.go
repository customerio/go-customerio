package customerio

import (
	"net/http"
	"time"
)

// DefaultHTTPTimeout is the timeout used by clients created without WithHTTPClient.
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
		return http.DefaultTransport
	}

	transport = transport.Clone()
	transport.MaxIdleConnsPerHost = 100
	return transport
}
