package customerio

import "testing"

func TestFormatPath(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []any
		want   string
	}{
		{
			name:   "no args",
			format: "/api/v1/events",
			want:   "/api/v1/events",
		},
		{
			name:   "string arg is escaped",
			format: "/api/v1/customers/%s",
			args:   []any{"john doe"},
			want:   "/api/v1/customers/john%20doe",
		},
		{
			name:   "slash inside a string arg is encoded as %2F",
			format: "/api/v1/customers/%s",
			args:   []any{"abc/def"},
			want:   "/api/v1/customers/abc%2Fdef",
		},
		{
			name:   "multiple string args",
			format: "/api/v1/customers/%s/devices/%s",
			args:   []any{"cust/1", "dev 2"},
			want:   "/api/v1/customers/cust%2F1/devices/dev%202",
		},
		{
			name:   "non-string arg passes through",
			format: "/api/v1/segments/%d/%s",
			args:   []any{42, "add_customers"},
			want:   "/api/v1/segments/42/add_customers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatPath(tt.format, tt.args...); got != tt.want {
				t.Errorf("formatPath() = %q, want %q", got, tt.want)
			}
		})
	}
}
