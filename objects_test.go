package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestGetObjectAttributes(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	ctx := context.Background()

	if _, err := api.GetObjectAttributes(ctx, "", "obj1", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "objectTypeID")
	}
	if _, err := api.GetObjectAttributes(ctx, "1", "", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "objectID")
	}

	api, rec := apiServer(t, 200, map[string]any{"object": map[string]any{
		"object_type_id": "1",
		"identifiers":    map[string]any{"object_id": "ae3000"},
		"attributes":     map[string]any{"name": "Acme"},
	}})
	got, err := api.GetObjectAttributes(ctx, "1", "ae3000", customerio.ObjectIdentifierTypeObjectID)
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/objects/1/ae3000/attributes?id_type=object_id")
	if got.Object.Identifiers.ObjectID != "ae3000" || got.Object.Attributes["name"] != "Acme" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetObjectRelationships(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	ctx := context.Background()

	if _, err := api.GetObjectRelationships(ctx, "", "obj1", customerio.ObjectRelationshipsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "objectTypeID")
	}
	if _, err := api.GetObjectRelationships(ctx, "1", "", customerio.ObjectRelationshipsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "objectID")
	}

	api, rec := apiServer(t, 200, map[string]any{"cio_relationships": []map[string]any{
		{"identifiers": map[string]any{"cio_id": "cio1", "id": "acmeInc"}, "object_type_id": "1", "object_type_disabled": false},
	}})
	got, err := api.GetObjectRelationships(ctx, "1", "ae3000", customerio.ObjectRelationshipsOptions{
		PaginationOptions: customerio.PaginationOptions{Limit: 20},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/objects/1/ae3000/relationships?limit=20")
	if len(got.Relationships) != 1 || got.Relationships[0].Identifiers.ID != "acmeInc" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestFindObjects(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	ctx := context.Background()

	if _, err := api.FindObjects(ctx, "", customerio.ObjectFilter{}, customerio.PaginationOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "objectTypeID")
	}

	api, rec := apiServer(t, 200, map[string]any{
		"identifiers": []map[string]any{{"cio_object_id": "ob020101", "object_id": "ae3000"}},
		"ids":         []string{"ae3000"},
	})
	filter := customerio.ObjectFilterByAttribute("name", customerio.FilterOperatorEq, "Acme", 1)
	got, err := api.FindObjects(ctx, "1", filter, customerio.PaginationOptions{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/objects?limit=10")

	wantBody := map[string]any{
		"object_type_id": "1",
		"filter": map[string]any{
			"object_attribute": map[string]any{"field": "name", "operator": "eq", "value": "Acme", "type_id": float64(1)},
		},
	}
	assertJSONEqual(t, rec.body, wantBody)
	if len(got.Identifiers) != 1 || got.Identifiers[0].ObjectID != "ae3000" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestFindObjectsComposedFilter(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"identifiers": []any{}, "ids": []any{}})

	filter := customerio.ObjectFilterOr(
		customerio.ObjectFilterByAttribute("tier", customerio.FilterOperatorEq, "gold", 1),
		customerio.ObjectFilterNot(customerio.ObjectFilterByAttribute("archived", customerio.FilterOperatorExists, "", 1)),
	)
	if _, err := api.FindObjects(context.Background(), "1", filter, customerio.PaginationOptions{}); err != nil {
		t.Fatal(err)
	}

	wantBody := map[string]any{
		"object_type_id": "1",
		"filter": map[string]any{
			"or": []any{
				map[string]any{"object_attribute": map[string]any{"field": "tier", "operator": "eq", "value": "gold", "type_id": float64(1)}},
				map[string]any{"not": map[string]any{"object_attribute": map[string]any{"field": "archived", "operator": "exists", "type_id": float64(1)}}},
			},
		},
	}
	assertJSONEqual(t, rec.body, wantBody)
}

func TestListObjectTypes(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"types": []map[string]any{
		{"id": "1", "name": "Concerts", "singular_name": "Concert", "enabled": true},
	}})
	got, err := api.ListObjectTypes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/object_types")
	if len(got.Types) != 1 || got.Types[0].Name != "Concerts" {
		t.Errorf("unexpected response: %#v", got)
	}
}
