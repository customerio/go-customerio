package customerio

import "net/http"

type option struct {
	api   func(*APIClient)
	track func(*CustomerIO)
}

type region struct {
	apiURL   string
	trackURL string
}

var (
	RegionUS = region{
		apiURL:   "https://api.customer.io",
		trackURL: "https://track.customer.io",
	}
	RegionEU = region{
		apiURL:   "https://api-eu.customer.io",
		trackURL: "https://track-eu.customer.io",
	}
)

func WithRegion(r region) option {
	return option{
		api: func(a *APIClient) {
			a.URL = r.apiURL
		},
		track: func(c *CustomerIO) {
			c.URL = r.trackURL
		},
	}
}

func WithHTTPClient(client *http.Client) option {
	return option{
		api: func(a *APIClient) {
			a.Client = client
		},
		track: func(c *CustomerIO) {
			c.Client = client
		},
	}
}
