package customerio

import "context"

// Workspace is a Customer.io workspace's identity and usage counters.
type Workspace struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name,omitempty"`
	MessagesSent         int64  `json:"messages_sent,omitempty"`
	BillableMessagesSent int64  `json:"billable_messages_sent,omitempty"`
	People               int64  `json:"people,omitempty"`
	ObjectTypes          int    `json:"object_types,omitempty"`
	Objects              int64  `json:"objects,omitempty"`
}

// ListWorkspacesResponse is the decoded shape of GET /v1/workspaces.
type ListWorkspacesResponse struct {
	Workspaces []Workspace `json:"workspaces"`
}

// ListWorkspaces returns every workspace the App API key can see.
// See https://docs.customer.io/api/app/#operation/listWorkspaces
func (c *APIClient) ListWorkspaces(ctx context.Context) (*ListWorkspacesResponse, error) {
	var resp ListWorkspacesResponse
	if err := c.doJSON(ctx, "GET", "/v1/workspaces", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IPAddressesResponse is the decoded shape of GET /v1/info/ip_addresses.
type IPAddressesResponse struct {
	IPAddresses []string `json:"ip_addresses"`
}

// GetIPAddresses returns the IP addresses Customer.io sends email from for
// this workspace's region.
// See https://docs.customer.io/api/app/#operation/getIpAddresses
func (c *APIClient) GetIPAddresses(ctx context.Context) (*IPAddressesResponse, error) {
	var resp IPAddressesResponse
	if err := c.doJSON(ctx, "GET", "/v1/info/ip_addresses", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
