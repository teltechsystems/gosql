package gosql

var (
	INNER_JOIN = "INNER JOIN"
)

type join struct {
	joinType  string
	table     table
	predicate string
}
