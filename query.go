package gosql

import (
	"database/sql"
	"strings"
)

type Query struct {
	from         *table
	joins        []join
	whereParts   []wherePart
	orderByParts []string
	using        *sql.DB
}

func (q *Query) From(tableName string, columns []string) *Query {
	q.from = &table{tableName, columns}
	return q
}

func (q *Query) getArgs() []interface{} {
	args := make([]interface{}, 0)

	for i := range q.whereParts {
		for j := range q.whereParts[i].args {
			args = append(args, q.whereParts[i].args[j])
		}
	}

	return args
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

func (q *Query) OrderBy(orderByParts []string) *Query {
	q.orderByParts = orderByParts
	return q
}

func (q *Query) Query() (*sql.Rows, error) {
	if q.using == nil {
		return nil, MissingDatabase
	}

	return q.using.Query(q.String(), q.getArgs()...)
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

		query += " WHERE (" + strings.Join(predicates, ") AND (") + ")"
	}

	// Build up the order by clause
	if len(q.orderByParts) > 0 {
		query += " ORDER BY " + strings.Join(q.orderByParts, ", ")
	}

	return query
}

func (q *Query) Use(db *sql.DB) *Query {
	q.using = db
	return q
}

func (q *Query) Where(predicate string, args ...interface{}) *Query {
	q.whereParts = append(q.whereParts, wherePart{predicate, args})
	return q
}
