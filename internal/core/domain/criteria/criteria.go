package criteria

type Criteria struct {
	Filters    []Filter
	Sort       Sort
	Pagination Pagination
}

type Filter struct {
	Field    string
	Value    any
	Operator Operator
}

type Operator string
type Direction string

const (
	OperatorEqual    Operator = "eq"
	OperatorNotEqual Operator = "ne"
	OperatorIn       Operator = "in"
	OperatorNotIn    Operator = "nin"
	OperatorGt       Operator = "gt"
	OperatorGte      Operator = "gte"
	OperatorLt       Operator = "lt"
	OperatorLte      Operator = "lte"
	OperatorRegex    Operator = "regex"
	OperatorExists   Operator = "exists"

	OrderDirectionDesc Direction = "desc"
	OrderDirectionAsc  Direction = "asc"

	DefaultPage     int = 1
	DefaultPageSize int = 50
)

type Sort struct {
	Field         string
	SortDirection Direction
}

type Pagination struct {
	Page     int
	PageSize int
}

type CriteriaBuilder struct {
	filters    []Filter
	sort       Sort
	pagination Pagination
}

func NewCriteriaBuilder() *CriteriaBuilder {
	return &CriteriaBuilder{
		filters:    make([]Filter, 0),
		sort:       Sort{},
		pagination: Pagination{},
	}
}

func (b *CriteriaBuilder) WithFilter(field string, value any, operator Operator) *CriteriaBuilder {
	b.filters = append(b.filters, Filter{
		Field:    field,
		Value:    value,
		Operator: operator,
	})
	return b
}

func (b *CriteriaBuilder) WithFilterAndCondition(field string, value any, operator Operator, condition bool) *CriteriaBuilder {
	if condition {
		b.WithFilter(field, value, operator)
	}
	return b
}

func (b *CriteriaBuilder) WithSort(field string, direction Direction) *CriteriaBuilder {
	b.sort = Sort{
		Field:         field,
		SortDirection: direction,
	}
	return b
}

func (b *CriteriaBuilder) WithPagination(page int, pageSize int) *CriteriaBuilder {
	b.pagination = Pagination{
		Page:     page,
		PageSize: pageSize,
	}
	return b
}

func (b *CriteriaBuilder) Build() *Criteria {
	return &Criteria{
		Filters:    b.filters,
		Sort:       b.sort,
		Pagination: b.pagination,
	}
}

func NewFilter(field string, value string, operator Operator) Filter {
	return Filter{
		Field:    field,
		Value:    value,
		Operator: operator,
	}
}

func NewSort(field string, direction Direction) Sort {
	return Sort{
		Field:         field,
		SortDirection: direction,
	}
}

func NewPagination(page int, pageSize int) Pagination {
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}
