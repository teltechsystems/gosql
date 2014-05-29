package gosql

var (
	INNER_JOIN = "INNER JOIN"
	LEFT_JOIN  = "LEFT JOIN"
)

type join struct {
	joinType  string
	table     table
	wherePart wherePart
}
