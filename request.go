package mediator

import "context"

type Request[TResponse interface{}] interface {
}

type RequestHandler[TRequest Request[TResponse], TResponse interface{}] interface {
	Handle(ctx context.Context, request TRequest) (TResponse, error)
}
