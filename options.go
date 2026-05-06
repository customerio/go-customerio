package customerio

import (
	"net/url"
	"time"
)

type option struct {
	api   func(*APIClient)
	track func(*CustomerIO)
}

type region struct {
	ApiURL   string
	TrackURL string
}

var (
	RegionUS = region{
		ApiURL:   "https://api.customer.io",
		TrackURL: "https://track.customer.io",
	}
	RegionEU = region{
		ApiURL:   "https://api-eu.customer.io",
		TrackURL: "https://track-eu.customer.io",
	}
)

func WithRegion(r region) option {
	return option{
		api: func(a *APIClient) {
			a.URL = r.ApiURL
		},
		track: func(c *CustomerIO) {
			c.URL = r.TrackURL
		},
	}
}

func WithHTTPClient(client HTTPClient) option {
	return option{
		api: func(a *APIClient) {
			a.Client = client
		},
		track: func(c *CustomerIO) {
			c.Client = client
		},
	}
}

func WithUserAgent(ua string) option {
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
type TrackOption func(map[string]interface{})

// TrackType is the type of event being tracked.
type TrackType string

const (
	TrackTypeEvent  TrackType = "event"
	TrackTypePage   TrackType = "page"
	TrackTypeScreen TrackType = "screen"
)

// WithEventID sets the id field on tracked events.
func WithEventID(id string) TrackOption {
	return func(payload map[string]interface{}) {
		payload["id"] = id
	}
}

// WithEventTimestamp sets the timestamp field on tracked events.
func WithEventTimestamp(timestamp time.Time) TrackOption {
	return func(payload map[string]interface{}) {
		payload["timestamp"] = timestamp.Unix()
	}
}

// WithEventType sets the type field on tracked events.
func WithEventType(typ TrackType) TrackOption {
	return func(payload map[string]interface{}) {
		payload["type"] = typ
	}
}

func trackPayload(eventName string, data map[string]interface{}, opts ...TrackOption) map[string]interface{} {
	payload := map[string]interface{}{
		"name": eventName,
		"data": data,
	}

	for _, opt := range opts {
		opt(payload)
	}

	return payload
}

// SegmentOption configures optional query parameters on segment membership requests.
type SegmentOption func(url.Values)

// WithSegmentIDType sets the id_type query parameter, identifying which kind of
// identifier the supplied ids represent. Defaults server-side to IdentifierTypeID.
func WithSegmentIDType(t IdentifierType) SegmentOption {
	return func(v url.Values) {
		v.Set("id_type", string(t))
	}
}
