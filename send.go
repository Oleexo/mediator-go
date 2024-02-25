package mediator

import (
	"context"

	"github.com/pkg/errors"
)

// SendWithoutContext sends a request to a single handler without a context
func SendWithoutContext[TRequest Request[TResponse], TResponse interface{}](container Container,
	request TRequest) (TResponse, error) {
	return Send[TRequest, TResponse](context.Background(), container, request)
}

// Send sends a request to a single handler
func Send[TRequest Request[TResponse], TResponse interface{}](ctx context.Context,
	container Container,
	request TRequest) (TResponse, error) {

	handler, exists := container.resolve(request)
	if !exists {
		return *new(TResponse), errors.Errorf("no handlers for request %T", request)
	}
	handlerValue, ok := handler.(RequestHandler[TRequest, TResponse])
	if !ok {
		return *new(TResponse), errors.Errorf("handler for request %T is not a Handle", request)
	}
	return executeWithPipeline[TRequest, TResponse](ctx, container.pipelineBehaviors(), handlerValue, request)
}

func executeWithPipeline[TRequest Request[TResponse], TResponse interface{}](ctx context.Context,
	pipelineBehaviors []PipelineBehavior,
	handler RequestHandler[TRequest, TResponse],
	request TRequest) (TResponse, error) {
	if len(pipelineBehaviors) > 0 {
		var reversPipes = reversOrder(pipelineBehaviors)

		var lastHandler RequestHandlerFunc = func() (interface{}, error) {
			return handler.Handle(ctx, request)
		}
		v := buildPipeline(reversPipes, lastHandler,
			func(next RequestHandlerFunc, pipe PipelineBehavior) RequestHandlerFunc {
				pipeValue := pipe
				nexValue := next

				var handlerFunc RequestHandlerFunc = func() (interface{}, error) {
					return pipeValue.Handle(ctx, request, nexValue)
				}

				return handlerFunc
			})

		response, err := v()

		if err != nil {
			if r, ok := response.(TResponse); ok {
				return r, err
			}
			return *new(TResponse), err
		}

		return response.(TResponse), nil
	} else {
		return handler.Handle(ctx, request)
	}
}

func buildPipeline(a []PipelineBehavior, seed RequestHandlerFunc,
	f func(RequestHandlerFunc, PipelineBehavior) RequestHandlerFunc) RequestHandlerFunc {
	result := seed
	for _, pipelineBehavior := range a {
		result = f(result, pipelineBehavior)
	}
	return result
}

func reversOrder(values []PipelineBehavior) []PipelineBehavior {
	var reverseValues []PipelineBehavior

	for i := len(values) - 1; i >= 0; i-- {
		reverseValues = append(reverseValues, values[i])
	}

	return reverseValues
}
