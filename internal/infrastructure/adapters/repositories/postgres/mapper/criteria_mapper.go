package mapper

import (
	"cpi-hub-api/internal/core/domain/criteria"
	"fmt"
	"strings"
)

type FilterClause struct {
	Query  string
	Params []interface{}
}

func ToPostgreSQLQuery(c *criteria.Criteria) (string, []interface{}) {
	var clauses []string
	var params []interface{}
	paramIndex := 1

	for _, filter := range c.Filters {
		clause, filterParams := buildFilterClause(filter, paramIndex)
		if clause != "" {
			clauses = append(clauses, clause)
			params = append(params, filterParams...)
			paramIndex += len(filterParams)
		}
	}

	whereClause := ""
	if len(clauses) > 0 {
		whereClause = "WHERE " + strings.Join(clauses, " AND ")
	}

	return whereClause, params
}

func buildFilterClause(filter criteria.Filter, startIndex int) (string, []interface{}) {
	var params []interface{}

	switch filter.Operator {
	case criteria.OperatorEqual:
		return fmt.Sprintf("%s = $%d", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorNotEqual:
		return fmt.Sprintf("%s != $%d", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorIn:
		if values, ok := filter.Value.([]interface{}); ok {
			placeholders := make([]string, len(values))
			for i, _ := range values {
				placeholders[i] = fmt.Sprintf("$%d", startIndex+i)
				params = append(params, values[i])
			}
			return fmt.Sprintf("%s IN (%s)", filter.Field, strings.Join(placeholders, ", ")), params
		}
		return fmt.Sprintf("%s IN ($%d)", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorNotIn:
		if values, ok := filter.Value.([]interface{}); ok {
			placeholders := make([]string, len(values))
			for i, _ := range values {
				placeholders[i] = fmt.Sprintf("$%d", startIndex+i)
				params = append(params, values[i])
			}
			return fmt.Sprintf("%s NOT IN (%s)", filter.Field, strings.Join(placeholders, ", ")), params
		}
		return fmt.Sprintf("%s NOT IN ($%d)", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorGt:
		return fmt.Sprintf("%s > $%d", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorGte:
		return fmt.Sprintf("%s >= $%d", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorLt:
		return fmt.Sprintf("%s < $%d", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorLte:
		return fmt.Sprintf("%s <= $%d", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorRegex:
		if pattern, ok := filter.Value.(string); ok {
			return fmt.Sprintf("%s ILIKE $%d", filter.Field, startIndex), []interface{}{pattern}
		}
		return "", nil

	case criteria.OperatorExists:
		if exists, ok := filter.Value.(bool); ok {
			if exists {
				return fmt.Sprintf("%s IS NOT NULL", filter.Field), nil
			} else {
				return fmt.Sprintf("%s IS NULL", filter.Field), nil
			}
		}
		return "", nil

	default:
		return "", nil
	}
}

func BuildOrderBy(sort criteria.Sort) string {
	if sort.Field == "" {
		return ""
	}

	direction := "ASC"
	if sort.SortDirection == criteria.OrderDirectionDesc {
		direction = "DESC"
	}

	return fmt.Sprintf("ORDER BY %s %s", sort.Field, direction)
}

func BuildPagination(pagination criteria.Pagination) string {
	if pagination.PageSize <= 0 {
		return ""
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	return fmt.Sprintf("LIMIT %d OFFSET %d", pagination.PageSize, offset)
}
