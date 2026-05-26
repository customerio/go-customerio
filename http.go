package customerio

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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

// doHTTP is the shared HTTP execution path for both CustomerIO (Track) and
// APIClient (App API). Auth header injection is caller-supplied via setAuth.
func doHTTP(ctx context.Context, client HTTPClient, method, url, userAgent string, body any, preflight func(*http.Request)) ([]byte, int, error) {
	var req *http.Request
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(b))
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		var err error
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, 0, err
		}
	}

	req.Header.Set("User-Agent", userAgent)
	preflight(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
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
