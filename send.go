package mediator

import (
	"context"
	"github.com/pkg/errors"
)

func SendWithoutContext[TRequest Request[TResponse], TResponse interface{}](container Container,
	request TRequest) (TResponse, error) {
	return Send[TRequest, TResponse](context.Background(), container, request)
}

func Send[TRequest Request[TResponse], TResponse interface{}](ctx context.Context,
	container Container,
	request TRequest) (TResponse, error) {

	handler, exists := container.resolve(request)
	if !exists {
		return *new(TResponse), errors.Errorf("no handlers for request %T",
			request)
	}
	handlerValue, ok := handler.(RequestHandler[TRequest, TResponse])
	if !ok {
		return *new(TResponse), errors.Errorf("handler for request %T is not a Handle",
			request)
	}
	return handlerValue.Handle(ctx, request)
}
