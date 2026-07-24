package customerio_test

import (
	"context"
	"testing"
)

func TestListWorkspaces(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"workspaces": []map[string]any{{"id": 1, "name": "Production"}}})
	got, err := api.ListWorkspaces(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/workspaces")
	if len(got.Workspaces) != 1 || got.Workspaces[0].Name != "Production" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetIPAddresses(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"ip_addresses": []string{"192.0.2.1", "192.0.2.2"}})
	got, err := api.GetIPAddresses(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/info/ip_addresses")
	if len(got.IPAddresses) != 2 {
		t.Errorf("unexpected response: %#v", got)
	}
}
