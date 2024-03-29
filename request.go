package mediator

import "context"

type BaseRequest interface {
	String() string
}

// Request is a marker interface for requests
// A request is a message that is sent to a single handler
type Request[TResponse interface{}] interface {
	BaseRequest
}

// RequestHandler is a marker interface for request handlers
// A request handler is a handler that handles a request
type RequestHandler[TRequest Request[TResponse], TResponse interface{}] interface {
	Handle(ctx context.Context, request TRequest) (TResponse, error)
}
