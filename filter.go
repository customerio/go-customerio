package customerio

// FilterOperator is a comparison operator supported in attribute and object
// attribute conditions.
type FilterOperator string

const (
	FilterOperatorEq     FilterOperator = "eq"
	FilterOperatorExists FilterOperator = "exists"
)

// AudienceFilter composes a boolean filter over customers, used by
// SearchCustomers and CreateCustomersExport. Build a leaf with
// FilterBySegment or FilterByAttribute; compose leaves with FilterAnd,
// FilterOr, or FilterNot. Exactly one of the fields is set on any given
// node — this mirrors the AudienceFilter union type in customerio-node's
// lib/types.ts, flattened to a single struct since Go has no union types.
type AudienceFilter struct {
	Segment   *SegmentCondition   `json:"segment,omitempty"`
	Attribute *AttributeCondition `json:"attribute,omitempty"`
	And       []AudienceFilter    `json:"and,omitempty"`
	Or        []AudienceFilter    `json:"or,omitempty"`
	Not       *AudienceFilter     `json:"not,omitempty"`
}

// SegmentCondition matches everyone in a given segment.
type SegmentCondition struct {
	ID int `json:"id"`
}

// AttributeCondition matches people whose attribute satisfies an operator
// against a value. Value is only meaningful for FilterOperatorEq; omit it
// for FilterOperatorExists.
type AttributeCondition struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    string         `json:"value,omitempty"`
}

// FilterBySegment builds a leaf AudienceFilter matching everyone in segment id.
func FilterBySegment(id int) AudienceFilter {
	return AudienceFilter{Segment: &SegmentCondition{ID: id}}
}

// FilterByAttribute builds a leaf AudienceFilter matching people whose
// attribute field satisfies operator (and value, for FilterOperatorEq).
func FilterByAttribute(field string, operator FilterOperator, value string) AudienceFilter {
	return AudienceFilter{Attribute: &AttributeCondition{Field: field, Operator: operator, Value: value}}
}

// FilterAnd composes filters that must all match.
func FilterAnd(filters ...AudienceFilter) AudienceFilter {
	return AudienceFilter{And: filters}
}

// FilterOr composes filters where at least one must match.
func FilterOr(filters ...AudienceFilter) AudienceFilter {
	return AudienceFilter{Or: filters}
}

// FilterNot negates a filter.
func FilterNot(filter AudienceFilter) AudienceFilter {
	return AudienceFilter{Not: &filter}
}

// ObjectFilter composes a boolean filter over objects, used by FindObjects.
// Build a leaf with ObjectFilterByAttribute; compose leaves with
// ObjectFilterAnd, ObjectFilterOr, or ObjectFilterNot.
type ObjectFilter struct {
	ObjectAttribute *ObjectAttributeCondition `json:"object_attribute,omitempty"`
	And             []ObjectFilter            `json:"and,omitempty"`
	Or              []ObjectFilter            `json:"or,omitempty"`
	Not             *ObjectFilter             `json:"not,omitempty"`
}

// ObjectAttributeCondition matches objects of type TypeID whose attribute
// field satisfies operator (and value, for FilterOperatorEq).
type ObjectAttributeCondition struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    string         `json:"value,omitempty"`
	TypeID   int            `json:"type_id"`
}

// ObjectFilterByAttribute builds a leaf ObjectFilter.
func ObjectFilterByAttribute(field string, operator FilterOperator, value string, typeID int) ObjectFilter {
	return ObjectFilter{ObjectAttribute: &ObjectAttributeCondition{Field: field, Operator: operator, Value: value, TypeID: typeID}}
}

// ObjectFilterAnd composes filters that must all match.
func ObjectFilterAnd(filters ...ObjectFilter) ObjectFilter {
	return ObjectFilter{And: filters}
}

// ObjectFilterOr composes filters where at least one must match.
func ObjectFilterOr(filters ...ObjectFilter) ObjectFilter {
	return ObjectFilter{Or: filters}
}

// ObjectFilterNot negates a filter.
func ObjectFilterNot(filter ObjectFilter) ObjectFilter {
	return ObjectFilter{Not: &filter}
}
