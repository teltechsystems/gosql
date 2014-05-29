package gosql

import (
	"strings"
)

type table struct {
	tableName string
	columns   []string
}

func (t table) GetAlias() string {
	parts := strings.Split(t.tableName, " ")

	return parts[len(parts)-1]
}
