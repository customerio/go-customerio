package customerio

import (
	"context"
	"fmt"
	"net/url"
)

type deviceV1 struct {
	ID         string         `json:"id"`
	Platform   string         `json:"platform"`
	LastUsed   string         `json:"last_used,omitempty"`
	Attributes map[string]any `json:"attributes"`
}

// Device identifies a push notification device for transactional sends.
type Device struct {
	Token      string         `json:"token"`
	Platform   string         `json:"platform"`
	LastUsed   string         `json:"last_used,omitempty"`
	Attributes map[string]any `json:"attributes"`
}

func newDeviceV1(deviceID, platform string, data map[string]any) (*deviceV1, error) {
	if deviceID == "" {
		return nil, ParamError{Param: "deviceID"}
	}
	if platform == "" {
		return nil, ParamError{Param: "platform"}
	}
	d := &deviceV1{
		ID:       deviceID,
		Platform: platform,
	}

	if len(data) > 0 {
		d.Attributes = make(map[string]any)
	}

	for k, v := range data {
		if k == "last_used" {
			d.LastUsed = fmt.Sprintf("%v", v)
			continue
		}
		d.Attributes[k] = v
	}

	return d, nil
}

// NewDevice prepares a push notification device for transactional sends.
func NewDevice(deviceID, platform string, data map[string]any) (*Device, error) {
	d, err := newDeviceV1(deviceID, platform, data)
	if err != nil {
		return nil, err
	}
	return &Device{
		Token:      d.ID,
		Platform:   d.Platform,
		Attributes: d.Attributes,
		LastUsed:   d.LastUsed,
	}, nil
}

// Delete deletes a customer
func (c *CustomerIO) Delete(customerID string) error {
	return c.DeleteCtx(context.Background(), customerID)
}

// AddDeviceCtx adds a device for a customer
func (c *CustomerIO) AddDeviceCtx(ctx context.Context, customerID string, deviceID string, platform string, data map[string]any) error {
	if customerID == "" {
		return ParamError{Param: "customerID"}
	}

	d, err := newDeviceV1(deviceID, platform, data)
	if err != nil {
		return err
	}

	body := map[string]any{
		"device": d,
	}

	return c.request(ctx, "PUT",
		fmt.Sprintf("%s/api/v1/customers/%s/devices", c.URL, url.PathEscape(customerID)),
		body)
}

// AddDevice adds a device for a customer
func (c *CustomerIO) AddDevice(customerID string, deviceID string, platform string, data map[string]any) error {
	return c.AddDeviceCtx(context.Background(), customerID, deviceID, platform, data)
}

// DeleteDeviceCtx deletes a device for a customer
func (c *CustomerIO) DeleteDeviceCtx(ctx context.Context, customerID string, deviceID string) error {
	if customerID == "" {
		return ParamError{Param: "customerID"}
	}
	if deviceID == "" {
		return ParamError{Param: "deviceID"}
	}
	return c.request(ctx, "DELETE",
		fmt.Sprintf("%s/api/v1/customers/%s/devices/%s", c.URL, url.PathEscape(customerID), url.PathEscape(deviceID)),
		nil)
}

// DeleteDevice deletes a device for a customer
func (c *CustomerIO) DeleteDevice(customerID string, deviceID string) error {
	return c.DeleteDeviceCtx(context.Background(), customerID, deviceID)
}
