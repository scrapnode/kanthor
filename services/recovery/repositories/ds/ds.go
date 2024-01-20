package ds

type Datastore interface {
	Message() Message
	Request() Request
}
