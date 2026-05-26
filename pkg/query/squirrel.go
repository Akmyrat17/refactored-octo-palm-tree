package query

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

func ApplyFilters(qb squirrel.SelectBuilder, filters []Filter, allowedFields map[string]string) (squirrel.SelectBuilder, error) {
	for _, f := range filters {
		column, ok := allowedFields[f.Field]
		if !ok {
			// Skip unknown fields instead of failing fully, or we can choose to return an error based on requirements.
			// Let's return error as suggested.
			return qb, fmt.Errorf("unknown filter field: %s", f.Field)
		}

		switch f.Operator {
		case "eq", "":
			qb = qb.Where(squirrel.Eq{column: f.Value})
		case "neq":
			qb = qb.Where(squirrel.NotEq{column: f.Value})
		case "gt":
			qb = qb.Where(squirrel.Gt{column: f.Value})
		case "gte":
			qb = qb.Where(squirrel.GtOrEq{column: f.Value})
		case "lt":
			qb = qb.Where(squirrel.Lt{column: f.Value})
		case "lte":
			qb = qb.Where(squirrel.LtOrEq{column: f.Value})
		case "like":
			qb = qb.Where(squirrel.Like{column: "%" + f.Value + "%"})
		case "ilike":
			qb = qb.Where(squirrel.ILike{column: "%" + f.Value + "%"})
		case "in":
			vals := strings.Split(f.Value, ",")
			qb = qb.Where(squirrel.Eq{column: vals})
		case "nin":
			vals := strings.Split(f.Value, ",")
			qb = qb.Where(squirrel.NotEq{column: vals})
		default:
			return qb, fmt.Errorf("unknown operator: %s", f.Operator)
		}
	}

	return qb, nil
}

func ApplySort(qb squirrel.SelectBuilder, sorts []SortField, allowedFields map[string]string) squirrel.SelectBuilder {
	for _, s := range sorts {
		column, ok := allowedFields[s.Field]
		if !ok {
			continue // ignore unallowed sort fields
		}
		if s.Desc {
			qb = qb.OrderBy(column + " DESC")
		} else {
			qb = qb.OrderBy(column + " ASC")
		}
	}

	return qb
}
