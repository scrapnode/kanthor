package ds

type Datastore interface {
	Request() Request
	Response() Response
}
