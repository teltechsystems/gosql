package gosql

type wherePart struct {
	predicate string
	args      []interface{}
}
