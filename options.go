package customerio

import "time"

// Option configures Customer.io API and Track clients.
type Option interface {
	applyAPI(*APIClient)
	applyTrack(*CustomerIO)
}

type option struct {
	api   func(*APIClient)
	track func(*CustomerIO)
}

func (o option) applyAPI(a *APIClient) {
	if o.api != nil {
		o.api(a)
	}
}

func (o option) applyTrack(c *CustomerIO) {
	if o.track != nil {
		o.track(c)
	}
}

// Region configures the Customer.io API endpoints for a workspace region.
type Region string

const (
	// RegionUS configures clients for Customer.io US endpoints.
	RegionUS Region = "us"
	// RegionEU configures clients for Customer.io EU endpoints.
	RegionEU Region = "eu"
)

func (r Region) APIURL() string {
	if r == RegionEU {
		return "https://api-eu.customer.io"
	}
	return "https://api.customer.io"
}

func (r Region) TrackURL() string {
	if r == RegionEU {
		return "https://track-eu.customer.io"
	}
	return "https://track.customer.io"
}

func WithRegion(r Region) Option {
	return option{
		api: func(a *APIClient) {
			a.URL = r.APIURL()
		},
		track: func(c *CustomerIO) {
			c.URL = r.TrackURL()
		},
	}
}

func WithHTTPClient(client HTTPClient) Option {
	return option{
		api: func(a *APIClient) {
			a.Client = client
		},
		track: func(c *CustomerIO) {
			c.Client = client
		},
	}
}

func WithUserAgent(ua string) Option {
	return option{
		api: func(a *APIClient) {
			a.UserAgent = ua
		},
		track: func(c *CustomerIO) {
			c.UserAgent = ua
		},
	}
}

// TrackOption sets optional top-level fields on tracked events.
type TrackOption func(map[string]any)

// TrackType is the type of event being tracked.
type TrackType string

const (
	TrackTypeEvent  TrackType = "event"
	TrackTypePage   TrackType = "page"
	TrackTypeScreen TrackType = "screen"
)

// WithEventID sets the id field on tracked events.
func WithEventID(id string) TrackOption {
	return func(payload map[string]any) {
		payload["id"] = id
	}
}

// WithEventTimestamp sets the timestamp field on tracked events.
func WithEventTimestamp(timestamp time.Time) TrackOption {
	return func(payload map[string]any) {
		payload["timestamp"] = timestamp.Unix()
	}
}

// WithEventType sets the type field on tracked events.
func WithEventType(typ TrackType) TrackOption {
	return func(payload map[string]any) {
		payload["type"] = typ
	}
}

func trackPayload(eventName string, data map[string]any, opts ...TrackOption) map[string]any {
	payload := map[string]any{
		"name": eventName,
		"data": data,
	}

	for _, opt := range opts {
		opt(payload)
	}

	return payload
}
