package gosql

import (
	"strings"
)

type Query struct {
	from       *table
	whereParts []wherePart
}

func (q *Query) From(tableName string, columns []string) *Query {
	q.from = &table{tableName, columns}
	return q
}

func (q *Query) String() string {
	if q.from == nil {
		return ""
	}

	query := "SELECT " + strings.Join(q.from.columns, ",") +
		" FROM " + q.from.tableName

	if len(q.whereParts) > 0 {
		predicates := make([]string, len(q.whereParts))
		for i := range q.whereParts {
			predicates[i] = q.whereParts[i].predicate
		}

		query += " WHERE " + strings.Join(predicates, " AND ")
	}

	return query
}

func (q *Query) Where(predicate string, args ...interface{}) *Query {
	q.whereParts = append(q.whereParts, wherePart{predicate, args})
	return q
}
