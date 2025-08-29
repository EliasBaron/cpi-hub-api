package criteria

import (
	"reflect"
	"testing"
)

func TestNewFilter(t *testing.T) {
	filter := NewFilter("testField", "testValue", OperatorEqual)
	expected := Filter{
		Field:    "testField",
		Value:    "testValue",
		Operator: OperatorEqual,
	}
	if !reflect.DeepEqual(filter, expected) {
		t.Errorf("Expected: %+v, got: %+v", expected, filter)
	}
}

func TestNewSort(t *testing.T) {
	sort := NewSort("testSort", OrderDirectionAsc)
	expected := Sort{
		Field:         "testSort",
		SortDirection: OrderDirectionAsc,
	}
	if !reflect.DeepEqual(sort, expected) {
		t.Errorf("Expected: %+v, got: %+v", expected, sort)
	}
}

func TestNewPagination(t *testing.T) {
	pag := NewPagination(1, 20)
	expected := Pagination{
		Page:     pag.Page,
		PageSize: pag.PageSize,
	}
	if !reflect.DeepEqual(pag, expected) {
		t.Errorf("Expected: %+v, got: %+v", expected, pag)
	}
}

func TestCriteriaBuilder_WithFilter(t *testing.T) {
	builder := NewCriteriaBuilder()
	builder.WithFilter("f1", "v1", OperatorEqual)
	builder.WithFilter("f2", "v2", OperatorNotEqual)
	criteria := builder.Build()

	if len(criteria.Filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(criteria.Filters))
	}
	expect := Filter{"f1", "v1", OperatorEqual}
	if !reflect.DeepEqual(criteria.Filters[0], expect) {
		t.Errorf("expected first filter %+v, got %+v", expect, criteria.Filters[0])
	}
	expect = Filter{"f2", "v2", OperatorNotEqual}
	if !reflect.DeepEqual(criteria.Filters[1], expect) {
		t.Errorf("expected second filter %+v, got %+v", expect, criteria.Filters[1])
	}
}

func TestCriteriaBuilder_WithSort(t *testing.T) {
	builder := NewCriteriaBuilder()
	builder.WithSort("sortField", OrderDirectionDesc)
	criteria := builder.Build()

	expectedSort := Sort{
		Field:         "sortField",
		SortDirection: OrderDirectionDesc,
	}
	if !reflect.DeepEqual(criteria.Sort, expectedSort) {
		t.Errorf("expected sort %+v, got %+v", expectedSort, criteria.Sort)
	}
}

func TestCriteriaBuilder_WithPagination(t *testing.T) {
	builder := NewCriteriaBuilder()
	builder.WithPagination(1, 20)
	criteria := builder.Build()

	expected := Pagination{
		Page:     1,
		PageSize: 20,
	}
	if !reflect.DeepEqual(criteria.Pagination, expected) {
		t.Errorf("expected pagination %+v, got %+v", expected, criteria.Pagination)
	}
}

func TestCriteriaBuilder_WithFilterAndSort(t *testing.T) {
	builder := NewCriteriaBuilder()
	builder.
		WithFilter("f", "v", OperatorEqual).
		WithSort("s", OrderDirectionAsc)
	criteria := builder.Build()

	expectedFilter := Filter{"f", "v", OperatorEqual}
	if len(criteria.Filters) != 1 {
		t.Errorf("expected 1 filter, got %d", len(criteria.Filters))
	}
	if !reflect.DeepEqual(criteria.Filters[0], expectedFilter) {
		t.Errorf("expected filter %+v, got %+v", expectedFilter, criteria.Filters[0])
	}

	expectedSort := Sort{"s", OrderDirectionAsc}
	if !reflect.DeepEqual(criteria.Sort, expectedSort) {
		t.Errorf("expected sort %+v, got %+v", expectedSort, criteria.Sort)
	}
}

func TestCriteriaBuilder_WithFilterAndSortAndPagination(t *testing.T) {
	builder := NewCriteriaBuilder()
	builder.
		WithFilter("f", "v", OperatorEqual).
		WithSort("s", OrderDirectionAsc).
		WithPagination(1, 20)

	criteria := builder.Build()

	expectedFilter := Filter{"f", "v", OperatorEqual}
	if len(criteria.Filters) != 1 {
		t.Errorf("expected 1 filter, got %d", len(criteria.Filters))
	}
	if !reflect.DeepEqual(criteria.Filters[0], expectedFilter) {
		t.Errorf("expected filter %+v, got %+v", expectedFilter, criteria.Filters[0])
	}

	expectedSort := Sort{"s", OrderDirectionAsc}
	if !reflect.DeepEqual(criteria.Sort, expectedSort) {
		t.Errorf("expected sort %+v, got %+v", expectedSort, criteria.Sort)
	}

	expected := Pagination{Page: 1, PageSize: 20}
	if !reflect.DeepEqual(criteria.Pagination, expected) {
		t.Errorf("expected pagination %+v, got %+v", expected, criteria.Pagination)
	}
}

func TestCriteriaBuilder_Empty(t *testing.T) {
	builder := NewCriteriaBuilder()
	criteria := builder.Build()

	if len(criteria.Filters) != 0 {
		t.Errorf("expected 0 filters, got %d", len(criteria.Filters))
	}
	if criteria.Sort != (Sort{}) {
		t.Errorf("expected empty sort, got %+v", criteria.Sort)
	}
}

func TestCriteriaBuilder_Chainability(t *testing.T) {
	builder := NewCriteriaBuilder()
	b := builder.WithFilter("field", "value", OperatorEqual)
	if b != builder {
		t.Error("WithFilter should return the builder itself for chaining")
	}
	b = builder.WithSort("sortField", OrderDirectionDesc)
	if b != builder {
		t.Error("WithSort should return the builder itself for chaining")
	}
}

func TestCriteriaBuilder_WithFilterAndCondition_True(t *testing.T) {
	builder := NewCriteriaBuilder()
	builder.WithFilterAndCondition("testField", "testValue", OperatorEqual, true)
	criteria := builder.Build()

	if len(criteria.Filters) != 1 {
		t.Fatalf("expected 1 filter when condition is true, got %d", len(criteria.Filters))
	}

	expectedFilter := Filter{
		Field:    "testField",
		Value:    "testValue",
		Operator: OperatorEqual,
	}
	if !reflect.DeepEqual(criteria.Filters[0], expectedFilter) {
		t.Errorf("expected filter %+v, got %+v", expectedFilter, criteria.Filters[0])
	}
}

func TestCriteriaBuilder_WithFilterAndCondition_False(t *testing.T) {
	builder := NewCriteriaBuilder()
	builder.WithFilterAndCondition("testField", "testValue", OperatorEqual, false)
	criteria := builder.Build()

	if len(criteria.Filters) != 0 {
		t.Errorf("expected 0 filters when condition is false, got %d", len(criteria.Filters))
	}
}

func TestCriteriaBuilder_WithFilterAndCondition_Chainability(t *testing.T) {
	builder := NewCriteriaBuilder()
	b := builder.WithFilterAndCondition("field", "value", OperatorEqual, true)
	if b != builder {
		t.Error("WithFilterAndCondition should return the builder itself for chaining")
	}
}

func TestCriteriaBuilder_WithFilterAndCondition_Multiple(t *testing.T) {
	builder := NewCriteriaBuilder()
	builder.
		WithFilterAndCondition("field1", "value1", OperatorEqual, true).
		WithFilterAndCondition("field2", "value2", OperatorNotEqual, false).
		WithFilterAndCondition("field3", "value3", OperatorIn, true)

	criteria := builder.Build()

	if len(criteria.Filters) != 2 {
		t.Fatalf("expected 2 filters (only those with condition=true), got %d", len(criteria.Filters))
	}

	expectedFilter1 := Filter{"field1", "value1", OperatorEqual}
	if !reflect.DeepEqual(criteria.Filters[0], expectedFilter1) {
		t.Errorf("expected first filter %+v, got %+v", expectedFilter1, criteria.Filters[0])
	}

	expectedFilter2 := Filter{"field3", "value3", OperatorIn}
	if !reflect.DeepEqual(criteria.Filters[1], expectedFilter2) {
		t.Errorf("expected second filter %+v, got %+v", expectedFilter2, criteria.Filters[1])
	}
}
