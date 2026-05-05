package customerio

import (
	"net/http"
)

func newDefaultHTTPClient() *http.Client {
	return &http.Client{
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
