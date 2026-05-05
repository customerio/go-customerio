package customerio

// Option configures Customer.io API and Track clients.
type Option struct {
	api   func(*APIClient)
	track func(*CustomerIO)
}

// Region configures the Customer.io API endpoints for a workspace region.
type Region struct {
	ApiURL   string
	TrackURL string
}

var (
	// RegionUS configures clients for Customer.io US endpoints.
	RegionUS = Region{
		ApiURL:   "https://api.customer.io",
		TrackURL: "https://track.customer.io",
	}
	// RegionEU configures clients for Customer.io EU endpoints.
	RegionEU = Region{
		ApiURL:   "https://api-eu.customer.io",
		TrackURL: "https://track-eu.customer.io",
	}
)

func WithRegion(r Region) Option {
	return Option{
		api: func(a *APIClient) {
			a.URL = r.ApiURL
		},
		track: func(c *CustomerIO) {
			c.URL = r.TrackURL
		},
	}
}

func WithHTTPClient(client HTTPClient) Option {
	return Option{
		api: func(a *APIClient) {
			a.Client = client
		},
		track: func(c *CustomerIO) {
			c.Client = client
		},
	}
}

func WithUserAgent(ua string) Option {
	return Option{
		api: func(a *APIClient) {
			a.UserAgent = ua
		},
		track: func(c *CustomerIO) {
			c.UserAgent = ua
		},
	}
}
