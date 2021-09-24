package customerio

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CustomerIO wraps the customer.io track API, see: https://customer.io/docs/api/#apitrackintroduction
type CustomerIO struct {
	siteID string
	apiKey string
	URL    string
	Client *http.Client
}

// CustomerIOError is returned by any method that fails at the API level
type CustomerIOError struct {
	status int
	url    string
	body   []byte
}

func (e *CustomerIOError) Error() string {
	return fmt.Sprintf("%v: %v %v", e.status, e.url, string(e.body))
}

// ParamError is an error returned if a parameter to the track API is invalid.
type ParamError struct {
	Param string // Param is the name of the parameter.
}

func (e ParamError) Error() string { return e.Param + ": missing" }

// NewTrackClient prepares a client for use with the Customer.io track API, see: https://customer.io/docs/api/#apitrackintroduction
// using a Tracking Site ID and API Key pair from https://fly.customer.io/settings/api_credentials
func NewTrackClient(siteID, apiKey string, opts ...option) *CustomerIO {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
		},
	}
	c := &CustomerIO{
		siteID: siteID,
		apiKey: apiKey,
		URL:    "https://track.customer.io",
		Client: client,
	}

	for _, opt := range opts {
		opt.track(c)
	}

	return c
}

// NewCustomerIO prepares a client for use with the Customer.io track API, see: https://customer.io/docs/api/#apitrackintroduction
// deprecated in favour of NewTrackClient
func NewCustomerIO(siteID, apiKey string) *CustomerIO {
	return NewTrackClient(siteID, apiKey)
}

// Identify identifies a customer and sets their attributes
func (c *CustomerIO) Identify(customerID string, attributes map[string]interface{}) error {
	if customerID == "" {
		return ParamError{Param: "customerID"}
	}
	return c.request("PUT",
		fmt.Sprintf("%s/api/v1/customers/%s", c.URL, url.PathEscape(customerID)),
		attributes)
}

// Track sends a single event to Customer.io for the supplied user
func (c *CustomerIO) Track(customerID string, eventName string, data map[string]interface{}) error {
	if customerID == "" {
		return ParamError{Param: "customerID"}
	}
	if eventName == "" {
		return ParamError{Param: "eventName"}
	}
	return c.request("POST",
		fmt.Sprintf("%s/api/v1/customers/%s/events", c.URL, url.PathEscape(customerID)),
		map[string]interface{}{
			"name": eventName,
			"data": data,
		})
}

// TrackAnonymous sends a single event to Customer.io for the anonymous user
func (c *CustomerIO) TrackAnonymous(anonymousID, eventName string, data map[string]interface{}) error {
	if anonymousID == "" {
		return ParamError{Param: "anonymousID"}
	}
	if eventName == "" {
		return ParamError{Param: "eventName"}
	}
	return c.request("POST",
		fmt.Sprintf("%s/api/v1/events", c.URL),
		map[string]interface{}{
			"name":         eventName,
			"anonymous_id": anonymousID,
			"data":         data,
		})
}

// Delete deletes a customer
func (c *CustomerIO) Delete(customerID string) error {
	if customerID == "" {
		return ParamError{Param: "customerID"}
	}
	return c.request("DELETE",
		fmt.Sprintf("%s/api/v1/customers/%s", c.URL, url.PathEscape(customerID)),
		nil)
}

// AddDevice adds a device for a customer
func (c *CustomerIO) AddDevice(customerID string, deviceID string, platform string, data map[string]interface{}) error {
	if customerID == "" {
		return ParamError{Param: "customerID"}
	}
	if deviceID == "" {
		return ParamError{Param: "deviceID"}
	}
	if platform == "" {
		return ParamError{Param: "platform"}
	}

	body := map[string]map[string]interface{}{
		"device": {
			"id":       deviceID,
			"platform": platform,
		},
	}
	for k, v := range data {
		body["device"][k] = v
	}
	return c.request("PUT",
		fmt.Sprintf("%s/api/v1/customers/%s/devices", c.URL, url.PathEscape(customerID)),
		body)
}

// DeleteDevice deletes a device for a customer
func (c *CustomerIO) DeleteDevice(customerID string, deviceID string) error {
	if customerID == "" {
		return ParamError{Param: "customerID"}
	}
	if deviceID == "" {
		return ParamError{Param: "deviceID"}
	}
	return c.request("DELETE",
		fmt.Sprintf("%s/api/v1/customers/%s/devices/%s", c.URL, url.PathEscape(customerID), url.PathEscape(deviceID)),
		nil)
}

func (c *CustomerIO) auth() string {
	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", c.siteID, c.apiKey)))
}

func (c *CustomerIO) request(method, url string, body interface{}) error {
	var req *http.Request
	if body != nil {
		j, err := json.Marshal(body)
		if err != nil {
			return err
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(j))
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(j)))
	} else {
		var err error
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return err
		}
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", c.auth()))

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return &CustomerIOError{
			status: resp.StatusCode,
			url:    url,
			body:   responseBody,
		}
	}

	return nil
}

const (
	// Identifier Types
	IdentifierTypeID    = "id"
	IdentifierTypeEmail = "email"
	IdentifierTypeCioID = "cio_id"
)

// MergeCustomers sends a request to Customer.io to merge two customer profiles together.
func (c *CustomerIO) MergeCustomers(primaryIDType, primaryID, secondaryIDType, secondaryID string) error {
	if !isValidIDType(primaryIDType) {
		return ParamError{Param: "primaryIDType"}
	}
	if strings.TrimSpace(primaryID) == "" {
		return ParamError{Param: "primaryID"}
	}

	if !isValidIDType(secondaryIDType) {
		return ParamError{Param: "secondaryIDType"}
	}
	if strings.TrimSpace(secondaryID) == "" {
		return ParamError{Param: "secondaryID"}
	}

	return c.request("POST",
		fmt.Sprintf("%s/api/v1/merge_customers", c.URL),
		map[string]interface{}{
			"primary": map[string]string{
				primaryIDType: primaryID,
			},
			"secondary": map[string]string{
				secondaryIDType: secondaryID,
			},
		})
}

func isValidIDType(input string) bool {
	input = strings.TrimSpace(input)
	if input == IdentifierTypeID ||
		input == IdentifierTypeEmail ||
		input == IdentifierTypeCioID {
		return true
	}
	return false
}
