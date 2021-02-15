package customerio

import "net/http"

type option struct {
	api   func(*APIClient)
	track func(*CustomerIO)
}

type region int

const (
	Region_US region = 0
	Region_EU region = 1

	apiURL_US = "https://api.customer.io"
	apiURL_EU = "https://api-eu.customer.io"

	trackURL_US = "https://track.customer.io"
	trackURL_EU = "https://track-eu.customer.io"
)

func WithRegion(r region) option {
	return option{
		api: func(a *APIClient) {
			switch r {
			case Region_US:
				a.URL = apiURL_US
			case Region_EU:
				a.URL = apiURL_EU
			}
		},
		track: func(c *CustomerIO) {
			switch r {
			case Region_US:
				c.URL = trackURL_US
			case Region_EU:
				c.URL = trackURL_EU
			}
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

// WithURL
func WithURL(url string) option {
	return option{
		api: func(a *APIClient) {
			a.URL = url
		},
		track: func(c *CustomerIO) {
			c.URL = url
		},
	}
}
