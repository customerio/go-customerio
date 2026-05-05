# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run with race detector (matches CI)
go test -race ./...

# Run a single test
go test -run TestName ./...
```

There is no Makefile and no lint step — CI only runs tests.

## Architecture

Single root package `customerio` (module `github.com/customerio/go-customerio/v3`), zero external dependencies.

### Two API clients

- **`CustomerIO`** ([customerio.go](customerio.go)) — Track API. Basic Auth (`siteID:apiKey`). Uses PUT, POST, and DELETE. Default transport sets `MaxIdleConnsPerHost: 100`.
- **`APIClient`** ([api.go](api.go)) — Transactional API (email, push, SMS, in-app, inbox). Bearer Token auth. POST only.

### Functional options

Both clients accept `...Option`. The `option` struct in [options.go](options.go) is shared. Common options: `WithRegion`, `WithHTTPClient`, `WithUserAgent`. Track-only options (`WithEventID`, `WithEventTimestamp`, `WithEventType`) wrap an inner `TrackOption`.

### HTTPClient interface

Defined in [http_client.go](http_client.go) as `Do(*http.Request) (*http.Response, error)`. Pass a custom implementation via `WithHTTPClient` to intercept or mock requests.

### Error types

- `ParamError` — invalid/missing parameter before the request is sent
- `CustomerIOError` — non-200 response from the Track API
- `TransactionalError` — non-200 response from the Transactional API (includes `meta.error` body)

### Regions

`Regions.US` and `Regions.EU` swap the base URLs for both clients. Default is US.

### Test pattern

Tests use `httptest.NewServer` set up in `TestMain` to capture the raw request (method, path, body). Assertions use `reflect.DeepEqual`. See [customerio_test.go](customerio_test.go) and [transactional_test.go](transactional_test.go) for examples.
