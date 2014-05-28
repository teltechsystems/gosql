package gosql

import (
	"strings"
)

type Query struct {
	from *table
}

func (q *Query) From(tableName string, columns []string) *Query {
	q.from = &table{tableName, columns}
	return q
}

func (q *Query) String() string {
	query := "SELECT " + strings.Join(q.from.columns, ",") +
		" FROM " + q.from.tableName

	return query
}
