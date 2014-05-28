package gosql

import (
	"strings"
)

type Query struct {
	from       *table
	joins      []join
	whereParts []wherePart
}

func (q *Query) From(tableName string, columns []string) *Query {
	q.from = &table{tableName, columns}
	return q
}

func (q *Query) Join(joinType string, tableName string, predicate string, columns []string) *Query {
	q.joins = append(q.joins, join{
		joinType:  joinType,
		table:     table{tableName, columns},
		predicate: predicate,
	})
	return q
}

func (q *Query) InnerJoin(tableName string, predicate string, columns []string) *Query {
	return q.Join(INNER_JOIN, tableName, predicate, columns)
}

func (q *Query) LeftJoin(tableName string, predicate string, columns []string) *Query {
	return q.Join(LEFT_JOIN, tableName, predicate, columns)
}

func (q *Query) String() string {
	if q.from == nil {
		return ""
	}

	// Generate the array of columns including those from joins
	columns := q.from.columns
	for i := range q.joins {
		for j := range q.joins[i].table.columns {
			columns = append(columns, q.joins[i].table.columns[j])
		}
	}

	query := "SELECT " + strings.Join(columns, ", ") +
		" FROM " + q.from.tableName

	// Build up the joins
	for i := range q.joins {
		query += " " + q.joins[i].joinType + " " + q.joins[i].table.tableName + " ON " + q.joins[i].predicate
	}

	// Build up the where conditions
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
