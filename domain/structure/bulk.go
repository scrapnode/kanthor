package structure

type BulkRes[T any] struct {
	Entity T
	Error  error
}
