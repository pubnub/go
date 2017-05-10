package pubnub

type Endpoint interface {
	BuildPath() string
	BuildQuery() map[string]string
	// or bytes[]?
	BuildBody() string
}

type TransactionalEndpoint interface {
	Sync() (interface{}, error)
	Async()
}
