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
	var whereParts []string
	var params []interface{}
	paramIndex := 1

	// WHERE
	for _, f := range c.Filters {
		clause, clauseParams := buildFilterClause(f, paramIndex)
		if clause != "" {
			whereParts = append(whereParts, clause)
			params = append(params, clauseParams...)
			paramIndex += len(clauseParams)
		}
	}

	query := ""
	if len(whereParts) > 0 {
		logicalOp := " AND "
		if c.LogicalOperator == criteria.LogicalOperatorOr {
			logicalOp = " OR "
		}

		if len(whereParts) > 1 {
			query += " WHERE (" + strings.Join(whereParts, logicalOp) + ")"
		} else {
			query += " WHERE " + whereParts[0]
		}
	}

	// ORDER BY
	if c.Sort.Field != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", c.Sort.Field, c.Sort.SortDirection)
	}

	// PAGINATE
	if c.Pagination.PageSize > 0 {
		offset := (c.Pagination.Page - 1) * c.Pagination.PageSize
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", c.Pagination.PageSize, offset)
	}

	return query, params
}

func buildFilterClause(filter criteria.Filter, startIndex int) (string, []interface{}) {
	var params []interface{}

	switch filter.Operator {
	case criteria.OperatorEqual:
		return fmt.Sprintf("%s = $%d", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorNotEqual:
		return fmt.Sprintf("%s != $%d", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorLike:
		return fmt.Sprintf("%s LIKE $%d", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorILike:
		return fmt.Sprintf("%s ILIKE $%d", filter.Field, startIndex), []interface{}{filter.Value}

	case criteria.OperatorIn:
		switch v := filter.Value.(type) {
		case []int:
			placeholders := make([]string, len(v))
			args := make([]interface{}, len(v))
			for i, val := range v {
				placeholders[i] = fmt.Sprintf("$%d", startIndex+i)
				args[i] = interface{}(val)
			}
			return fmt.Sprintf("%s IN (%s)", filter.Field, strings.Join(placeholders, ",")), args
		default:
			return "", nil
		}

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
