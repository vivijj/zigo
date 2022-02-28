package state

type QueryResult[T any] struct {
	// if the query result is not none
	Ok bool
	// fail when query.
	Err error
}
