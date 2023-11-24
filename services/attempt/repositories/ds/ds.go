package ds

type Datastore interface {
	Message() Message
	Request() Request
	Response() Response
	Attempt() Attempt
}

type ScanResults[T any] struct {
	Data  T
	Error error
}
