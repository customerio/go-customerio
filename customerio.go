package customerio

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
)

// CustomerIO wraps the customer.io API, see: http://customer.io/docs/api/rest.html
type CustomerIO struct {
	siteID string
	apiKey string
	Host   string
	SSL    bool
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

var DefaultClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 100,
	},
}

// NewCustomerIO creates a new CustomerIO object to perform requests on the supplied credentials
func NewCustomerIO(siteID, apiKey string) *CustomerIO {
	return &CustomerIO{siteID, apiKey, "track.customer.io", true, DefaultClient}
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

// TrackAnonymous sends a single event to Customer.io for the anonymous user
func (c *CustomerIO) TrackAnonymous(eventName string, data map[string]interface{}) error {
	body := map[string]interface{}{"name": eventName, "data": data}
	j, err := json.Marshal(body)

	if err != nil {
		return err
	}

	status, responseBody, err := c.request("POST", c.anonURL(), j)

	if err != nil {
		return err
	} else if status != 200 {
		return &CustomerIOError{status, c.anonURL(), responseBody}
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

// AddDevice adds a device for a customer
func (c *CustomerIO) AddDevice(customerID string, deviceID string, platform string, data map[string]interface{}) error {
	if customerID == "" {
		return errors.New("customerID is a required field")
	}
	if deviceID == "" {
		return errors.New("deviceID is a required field")
	}
	if platform == "" {
		return errors.New("platform is a required field")
	}

	body := map[string]map[string]interface{}{"device": {"id": deviceID, "platform": platform}}
	for k, v := range data {
		body["device"][k] = v
	}
	j, err := json.Marshal(body)

	if err != nil {
		return err
	}

	status, responseBody, err := c.request("PUT", c.deviceURL(customerID), j)

	if err != nil {
		return err
	} else if status != 200 {
		return &CustomerIOError{status, c.deviceURL(customerID), responseBody}
	}

	return nil
}

// DeleteDevice deletes a device for a customer
func (c *CustomerIO) DeleteDevice(customerID string, deviceID string) error {
	status, responseBody, err := c.request("DELETE", c.deleteDeviceURL(customerID, deviceID), []byte{})

	if err != nil {
		return err
	} else if status != 200 {
		return &CustomerIOError{status, c.deleteDeviceURL(customerID, deviceID), responseBody}
	}

	return nil
}

func (c *CustomerIO) AddCustomersToSegment(segmentID int, customerIDs []string) error {
	if segmentID == 0 {
		return errors.New("segmentID is a required field")
	}
	if len(customerIDs) == 0 {
		return errors.New("customerIDs is a required field")
	}

	body := map[string]interface{}{"ids": customerIDs}
	j, err := json.Marshal(body)
	if err != nil {
		return err
	}

	status, responseBody, err := c.request("POST", c.addCustomersToManualSegmentURL(segmentID), j)
	if err != nil {
		return err
	} else if status != 200 {
		return &CustomerIOError{status, c.addCustomersToManualSegmentURL(segmentID), responseBody}
	}

	return nil
}

func (c *CustomerIO) RemoveCustomersFromSegment(segmentID int, customerIDs []string) error {
	if segmentID == 0 {
		return errors.New("segmentID is a required field")
	}
	if len(customerIDs) == 0 {
		return errors.New("customerIDs is a required field")
	}

	body := map[string]interface{}{"ids": customerIDs}
	j, err := json.Marshal(body)
	if err != nil {
		return err
	}

	status, responseBody, err := c.request("POST", c.removeCustomersFromManualSegmentURL(segmentID), j)
	if err != nil {
		return err
	} else if status != 200 {
		return &CustomerIOError{status, c.removeCustomersFromManualSegmentURL(segmentID), responseBody}
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
	return c.protocol() + path.Join(c.Host, "api/v1", "customers", encodeID(customerID))
}

func (c *CustomerIO) eventURL(customerID string) string {
	return c.protocol() + path.Join(c.Host, "api/v1", "customers", encodeID(customerID), "events")
}

func (c *CustomerIO) anonURL() string {
	return c.protocol() + path.Join(c.Host, "api/v1", "events")
}

func (c *CustomerIO) deviceURL(customerID string) string {
	return c.protocol() + path.Join(c.Host, "api/v1", "customers", encodeID(customerID), "devices")
}

func (c *CustomerIO) deleteDeviceURL(customerID string, deviceID string) string {
	return c.protocol() + path.Join(c.Host, "api/v1", "customers", encodeID(customerID), "devices", deviceID)
}

func (c *CustomerIO) addCustomersToManualSegmentURL(segmentID int) string {
	return c.protocol() + path.Join(c.Host, "api/v1/", "segments", strconv.Itoa(segmentID), "add_customers")
}

func (c *CustomerIO) removeCustomersFromManualSegmentURL(segmentID int) string {
	return c.protocol() + path.Join(c.Host, "api/v1/", "segments", strconv.Itoa(segmentID), "remove_customers")
}

func (c *CustomerIO) request(method, url string, body []byte) (status int, responseBody []byte, err error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", c.auth()))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(body)))

	resp, err := c.Client.Do(req)
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
