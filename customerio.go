package customerio

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
)

// CustomerIO wraps the customer.io API, see: https://learn.customer.io/api/
type CustomerIO struct {
	siteID  string
	apiKey  string
	Host    string
	APIHost string
	SSL     bool
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

// NewCustomerIO creates a new CustomerIO object to perform requests on the supplied credentials
func NewCustomerIO(siteID, apiKey string) *CustomerIO {
	return &CustomerIO{siteID, apiKey, "track.customer.io", "api.customer.io", true}
}

// Identify identifies a customer and sets their attributes
func (c *CustomerIO) Identify(customerID string, attributes map[string]interface{}) error {
	j, err := json.Marshal(attributes)

	if err != nil {
		return err
	}

	status, responseBody, err := c.request("PUT", c.customerURL(customerID), j)

	if err != nil {
		return err
	} else if status != 200 {
		return &CustomerIOError{status, c.customerURL(customerID), responseBody}
	}

	return nil
}

// Track sends a single event to Customer.io for the supplied user
func (c *CustomerIO) Track(customerID string, eventName string, data map[string]interface{}) error {

	body := map[string]interface{}{"name": eventName, "data": data}
	j, err := json.Marshal(body)

	if err != nil {
		return err
	}

	status, responseBody, err := c.request("POST", c.eventURL(customerID), j)

	if err != nil {
		return err
	} else if status != 200 {
		return &CustomerIOError{status, c.eventURL(customerID), responseBody}
	}

	return nil
}

// Trigger triggers a broadcast campaign for the supplied campaign ID
func (c *CustomerIO) Trigger(campaignID string, data map[string]interface{}) error {

	body := make(map[string]interface{})
	if recipients, ok := data["recipients"]; ok {
		delete(data, "recipients")
		body = map[string]interface{}{"data": data, "recipients": recipients}
	} else {
		body = map[string]interface{}{"data": data}
	}

	j, err := json.Marshal(body)

	if err != nil {
		return err
	}

	status, responseBody, err := c.request("POST", c.campaignURL(campaignID), j)

	if err != nil {
		return err
	} else if status != 200 {
		return &CustomerIOError{status, c.campaignURL(campaignID), responseBody}
	}

	return nil
}

// Delete deletes a customer
func (c *CustomerIO) Delete(customerID string) error {
	status, responseBody, err := c.request("DELETE", c.customerURL(customerID), []byte{})

	if err != nil {
		return err
	} else if status != 200 {
		return &CustomerIOError{status, c.customerURL(customerID), responseBody}
	}

	return nil
}

func (c *CustomerIO) auth() string {
	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", c.siteID, c.apiKey)))
}

func (c *CustomerIO) protocol() string {
	if !c.SSL {
		return "http://"
	}
	return "https://"
}

func (c *CustomerIO) customerURL(customerID string) string {
	return c.protocol() + path.Join(c.Host, "api/v1", "customers", customerID)
}

func (c *CustomerIO) campaignURL(campaignID string) string {
	return c.protocol() + path.Join(c.APIHost, "/v1/api", "campaigns", campaignID, "triggers")
}

func (c *CustomerIO) eventURL(customerID string) string {
	return c.protocol() + path.Join(c.Host, "api/v1", "customers", customerID, "events")
}

func (c *CustomerIO) request(method, url string, body []byte) (status int, responseBody []byte, err error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", c.auth()))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(body)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	status = resp.StatusCode
	if resp.ContentLength >= 0 {
		responseBody = make([]byte, resp.ContentLength)
		resp.Body.Read(responseBody)
	}

	return status, responseBody, nil

}
