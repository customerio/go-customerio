package customerio_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestSearchCustomers(t *testing.T) {
	ctx := context.Background()

	resp := map[string]any{
		"identifiers": []map[string]any{{"id": "1", "email": "a@example.com", "cio_id": "c1"}},
		"ids":         []string{"1"},
		"next":        "MDox",
	}
	api, rec := apiServer(t, 200, resp)

	filter := customerio.FilterByAttribute("email", customerio.FilterOperatorEq, "a@example.com")
	got, err := api.SearchCustomers(ctx, filter, customerio.PaginationOptions{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/customers?limit=10")

	want := &customerio.SearchCustomersResponse{
		Identifiers: []customerio.CustomerIdentifiers{{ID: "1", Email: "a@example.com", CioID: "c1"}},
		IDs:         []string{"1"},
		Next:        "MDox",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}

	wantBody := map[string]any{
		"filter": map[string]any{"attribute": map[string]any{"field": "email", "operator": "eq", "value": "a@example.com"}},
	}
	assertJSONEqual(t, rec.body, wantBody)
}

func TestSearchCustomersComposedFilter(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"identifiers": []any{}, "ids": []any{}})

	filter := customerio.FilterAnd(
		customerio.FilterBySegment(7),
		customerio.FilterNot(customerio.FilterByAttribute("plan", customerio.FilterOperatorExists, "")),
	)
	if _, err := api.SearchCustomers(context.Background(), filter, customerio.PaginationOptions{}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/customers")

	wantBody := map[string]any{
		"filter": map[string]any{
			"and": []any{
				map[string]any{"segment": map[string]any{"id": float64(7)}},
				map[string]any{"not": map[string]any{"attribute": map[string]any{"field": "plan", "operator": "exists"}}},
			},
		},
	}
	assertJSONEqual(t, rec.body, wantBody)
}

func TestGetCustomersByEmail(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCustomersByEmail(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "email")
	}

	api, rec := apiServer(t, 200, map[string]any{"results": []map[string]any{{"id": "1", "email": "a@example.com"}}})
	got, err := api.GetCustomersByEmail(context.Background(), "a@example.com")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/customers?email=a%40example.com")
	want := &customerio.GetCustomersByEmailResponse{Results: []customerio.CustomerIdentifiers{{ID: "1", Email: "a@example.com"}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestGetAttributes(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetAttributes(context.Background(), "", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "id")
	}

	api, rec := apiServer(t, 200, map[string]any{"customer": map[string]any{
		"id":           "1",
		"attributes":   map[string]any{"plan": "gold"},
		"unsubscribed": false,
	}})
	got, err := api.GetAttributes(context.Background(), "1", "")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/customers/1/attributes?id_type=id")
	if got.Customer.Attributes["plan"] != "gold" {
		t.Errorf("unexpected attributes: %#v", got.Customer.Attributes)
	}

	api2, rec2 := apiServer(t, 200, map[string]any{"customer": map[string]any{"id": "a@example.com"}})
	if _, err := api2.GetAttributes(context.Background(), "a@example.com", customerio.IdentifierTypeEmail); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec2, "GET", "/v1/customers/a@example.com/attributes?id_type=email")
}

func TestGetCustomerActivities(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCustomerActivities(context.Background(), "", customerio.CustomerActivitiesOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}

	api, rec := apiServer(t, 200, map[string]any{"activities": []map[string]any{
		{"id": "act1", "type": "sent_email", "timestamp": 1397566226},
	}, "next": "abc"})
	got, err := api.GetCustomerActivities(context.Background(), "1", customerio.CustomerActivitiesOptions{
		IDType:            customerio.IdentifierTypeID,
		Type:              "sent_email",
		PaginationOptions: customerio.PaginationOptions{Limit: 5},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/customers/1/activities?id_type=id&limit=5&type=sent_email")
	if len(got.Activities) != 1 || got.Activities[0].ID != "act1" || got.Next != "abc" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetCustomerMessages(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCustomerMessages(context.Background(), "", customerio.CustomerMessagesOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}

	api, rec := apiServer(t, 200, map[string]any{"messages": []map[string]any{
		{"id": "m1", "type": "email", "subject": "hi", "created": 1700000000},
	}})
	got, err := api.GetCustomerMessages(context.Background(), "a@example.com", customerio.CustomerMessagesOptions{
		IDType:  customerio.IdentifierTypeEmail,
		StartTS: 1690000000,
		EndTS:   1700000000,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/customers/a@example.com/messages?end_ts=1700000000&id_type=email&start_ts=1690000000")
	if len(got.Messages) != 1 || got.Messages[0].Subject != "hi" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetCustomerRelationships(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCustomerRelationships(context.Background(), "", customerio.PaginationOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}

	api, rec := apiServer(t, 200, map[string]any{"cio_relationships": []map[string]any{
		{"object_type_id": "1", "identifiers": map[string]any{"object_id": "ae3000"}},
	}, "next": ""})
	got, err := api.GetCustomerRelationships(context.Background(), "1", customerio.PaginationOptions{Start: "cursor1"})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/customers/1/relationships?start=cursor1")
	if len(got.Relationships) != 1 || got.Relationships[0].Identifiers.ObjectID != "ae3000" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetCustomerSegments(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCustomerSegments(context.Background(), "", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}

	api, rec := apiServer(t, 200, map[string]any{"segments": []map[string]any{{"id": 3, "name": "VIP"}}})
	got, err := api.GetCustomerSegments(context.Background(), "1", "")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/customers/1/segments?id_type=id")
	if len(got.Segments) != 1 || got.Segments[0].Name != "VIP" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetCustomerSubscriptionPreferences(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCustomerSubscriptionPreferences(context.Background(), "", customerio.CustomerSubscriptionPreferencesOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}

	api, rec := apiServer(t, 200, map[string]any{"customer": map[string]any{
		"id": "1",
		"topics": []map[string]any{
			{"id": 1, "subscribed": true, "name": "Offers"},
		},
		"unsubscribed": false,
	}})
	got, err := api.GetCustomerSubscriptionPreferences(context.Background(), "1", customerio.CustomerSubscriptionPreferencesOptions{Language: "en"})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/customers/1/subscription_preferences?language=en")
	if len(got.Customer.Topics) != 1 || got.Customer.Topics[0].Name != "Offers" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetCustomersAttributes(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetCustomersAttributes(context.Background(), nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "ids")
	}

	api, rec := apiServer(t, 200, map[string]any{"customers": []map[string]any{
		{"customer": map[string]any{"id": "1", "attributes": map[string]any{"plan": "gold"}}},
	}})
	got, err := api.GetCustomersAttributes(context.Background(), []string{"1", "2"})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/customers/attributes")
	wantBody := map[string]any{"ids": []any{"1", "2"}}
	assertJSONEqual(t, rec.body, wantBody)
	if len(got.Customers) != 1 || got.Customers[0].Customer.ID != "1" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetSubscriptionCenterToken(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetSubscriptionCenterToken(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "customerID")
	}

	api, rec := apiServer(t, 200, map[string]any{"token": "tok123", "url": "https://track.customer.io/u/i/tok123/"})
	got, err := api.GetSubscriptionCenterToken(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/subscription_center/1/token")
	if got.Token != "tok123" {
		t.Errorf("unexpected token: %#v", got)
	}
}

func TestCustomersEndpointsError(t *testing.T) {
	api, _ := apiServer(t, 404, map[string]any{"errors": []map[string]any{{"detail": "not found"}}})

	_, err := api.GetAttributes(context.Background(), "missing", "")
	if err == nil {
		t.Fatal("expected error")
	}
	cioErr, ok := err.(*customerio.CustomerIOError)
	if !ok {
		t.Fatalf("expected *CustomerIOError, got %T", err)
	}
	if cioErr.StatusCode() != 404 {
		t.Errorf("expected status 404, got %d", cioErr.StatusCode())
	}
}
