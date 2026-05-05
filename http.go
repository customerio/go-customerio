package customerio

import (
	"net/http"
	"time"
)

const DefaultHTTPTimeout = 30 * time.Second

func newDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   DefaultHTTPTimeout,
		Transport: newDefaultTransport(),
	}
}

func newDefaultTransport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 100
	return transport
}
