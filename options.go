package customerio

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
