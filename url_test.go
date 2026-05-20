package customerio

import (
	"net/url"
	"testing"
)

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		query    url.Values
		segments []string
		want     string
	}{
		{
			name:     "no query",
			base:     "https://track.customer.io",
			segments: []string{"api", "v1", "events"},
			want:     "https://track.customer.io/api/v1/events",
		},
		{
			name:     "base with trailing slash",
			base:     "https://track.customer.io/",
			segments: []string{"api", "v1", "events"},
			want:     "https://track.customer.io/api/v1/events",
		},
		{
			name:     "nil query produces no query string",
			base:     "https://track.customer.io",
			query:    nil,
			segments: []string{"api", "v1", "segments", "42", "add_customers"},
			want:     "https://track.customer.io/api/v1/segments/42/add_customers",
		},
		{
			name:     "with query",
			base:     "https://track.customer.io",
			query:    url.Values{"id_type": {"email"}},
			segments: []string{"api", "v1", "segments", "42", "add_customers"},
			want:     "https://track.customer.io/api/v1/segments/42/add_customers?id_type=email",
		},
		{
			name:     "dynamic segment with special characters is escaped",
			base:     "https://track.customer.io",
			segments: []string{"api", "v1", "customers", "john doe"},
			want:     "https://track.customer.io/api/v1/customers/john%20doe",
		},
		{
			name:     "slash inside a segment is encoded as %2F",
			base:     "https://track.customer.io",
			segments: []string{"api", "v1", "customers", "abc/def"},
			want:     "https://track.customer.io/api/v1/customers/abc%2Fdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildURL(tt.base, tt.query, tt.segments...)
			if err != nil {
				t.Fatalf("buildURL() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("buildURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
