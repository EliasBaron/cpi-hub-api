package mapper

import (
	"cpi-hub-api/internal/core/domain/criteria"
	"testing"
)

func TestToPostgreSQLQuery(t *testing.T) {
	tests := []struct {
		name           string
		criteria       *criteria.Criteria
		expectedQuery  string
		expectedParams []interface{}
	}{
		{
			name: "Single equal filter",
			criteria: &criteria.Criteria{
				Filters: []criteria.Filter{
					{Field: "name", Value: "John", Operator: criteria.OperatorEqual},
				},
			},
			expectedQuery:  "WHERE name = $1",
			expectedParams: []interface{}{"John"},
		},
		{
			name: "Multiple filters",
			criteria: &criteria.Criteria{
				Filters: []criteria.Filter{
					{Field: "age", Value: 25, Operator: criteria.OperatorGte},
					{Field: "status", Value: "active", Operator: criteria.OperatorEqual},
				},
			},
			expectedQuery:  "WHERE age >= $1 AND status = $2",
			expectedParams: []interface{}{25, "active"},
		},
		{
			name: "IN operator with slice",
			criteria: &criteria.Criteria{
				Filters: []criteria.Filter{
					{Field: "status", Value: []interface{}{"active", "pending"}, Operator: criteria.OperatorIn},
				},
			},
			expectedQuery:  "WHERE status IN ($1, $2)",
			expectedParams: []interface{}{"active", "pending"},
		},
		{
			name: "Regex operator",
			criteria: &criteria.Criteria{
				Filters: []criteria.Filter{
					{Field: "name", Value: "%john%", Operator: criteria.OperatorRegex},
				},
			},
			expectedQuery:  "WHERE name ILIKE $1",
			expectedParams: []interface{}{"%john%"},
		},
		{
			name: "Exists operator true",
			criteria: &criteria.Criteria{
				Filters: []criteria.Filter{
					{Field: "email", Value: true, Operator: criteria.OperatorExists},
				},
			},
			expectedQuery:  "WHERE email IS NOT NULL",
			expectedParams: []interface{}{},
		},
		{
			name: "No filters",
			criteria: &criteria.Criteria{
				Filters: []criteria.Filter{},
			},
			expectedQuery:  "",
			expectedParams: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, params := ToPostgreSQLQuery(tt.criteria)

			if query != tt.expectedQuery {
				t.Errorf("Expected query '%s', got '%s'", tt.expectedQuery, query)
			}

			if len(params) != len(tt.expectedParams) {
				t.Errorf("Expected %d params, got %d", len(tt.expectedParams), len(params))
			}

			for i, param := range params {
				if i < len(tt.expectedParams) && param != tt.expectedParams[i] {
					t.Errorf("Expected param[%d] = %v, got %v", i, tt.expectedParams[i], param)
				}
			}
		})
	}
}

func TestBuildOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		sort     criteria.Sort
		expected string
	}{
		{
			name:     "Ascending sort",
			sort:     criteria.Sort{Field: "name", SortDirection: criteria.OrderDirectionAsc},
			expected: "ORDER BY name ASC",
		},
		{
			name:     "Descending sort",
			sort:     criteria.Sort{Field: "created_at", SortDirection: criteria.OrderDirectionDesc},
			expected: "ORDER BY created_at DESC",
		},
		{
			name:     "Empty field",
			sort:     criteria.Sort{Field: "", SortDirection: criteria.OrderDirectionAsc},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildOrderBy(tt.sort)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestBuildPagination(t *testing.T) {
	tests := []struct {
		name       string
		pagination criteria.Pagination
		expected   string
	}{
		{
			name:       "Normal pagination",
			pagination: criteria.Pagination{Page: 2, PageSize: 10},
			expected:   "LIMIT 10 OFFSET 10",
		},
		{
			name:       "First page",
			pagination: criteria.Pagination{Page: 1, PageSize: 20},
			expected:   "LIMIT 20 OFFSET 0",
		},
		{
			name:       "Zero page size",
			pagination: criteria.Pagination{Page: 1, PageSize: 0},
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildPagination(tt.pagination)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
