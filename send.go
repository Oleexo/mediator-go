package mediator

import (
	"context"
	"fmt"
	"reflect"
)

// SendWithoutContext sends a request using the provided container without requiring a context.
// It defaults to using context.Background().
func SendWithoutContext[TRequest Request[TResponse], TResponse interface{}](container SendContainer,
	request TRequest) (TResponse, error) {
	return Send[TRequest, TResponse](context.Background(), container, request)
}

// Send sends a request to a handler via the specified SendContainer and returns the handler's response or an error.
func Send[TRequest Request[TResponse], TResponse interface{}](ctx context.Context,
	container SendContainer,
	request TRequest) (TResponse, error) {

	handler, exists := container.resolve(request)
	if !exists {
		return *new(TResponse), fmt.Errorf("no handlers for request %T", request)
	}
	handlerValue, ok := handler.(RequestHandler[TRequest, TResponse])
	if !ok {
		return *new(TResponse), fmt.Errorf("handler for request %T is not a Handle", request)
	}
	var requestHandlerBehavior RequestHandlerFunc = func() (interface{}, error) {
		return handlerValue.Handle(ctx, request)
	}
	response, err := container.executeWithPipeline(ctx, request, requestHandlerBehavior)
	if err != nil {
		if r, ok := response.(TResponse); ok {
			return r, err
		}
		return *new(TResponse), err
	}
	return response.(TResponse), nil
}

// RequestPipelineBehavior is a marker interface for pipeline behaviors
// A pipeline behavior is a behavior that is executed as part of a pipeline
type RequestPipelineBehavior interface {

	// Handle processes the BaseRequest within the pipeline and calls the next RequestHandlerFunc.
	Handle(ctx context.Context, request BaseRequest, next RequestHandlerFunc) (interface{}, error)
}

// Sender defines an interface for sending requests and retrieving responses within a given context.
type Sender interface {

	// Send executes the provided `BaseRequest` within the given context and returns the response or an error.
	Send(ctx context.Context, request BaseRequest) (interface{}, error)
}

// RequestHandlerFunc is a function that handles a request
type RequestHandlerFunc func() (interface{}, error)

type sender struct {
	container SendContainer
}

// NewSender creates a new Sender instance with the provided SendContainer for request handling and execution.
func NewSender(container SendContainer) Sender {
	return &sender{
		container: container,
	}
}

func (s sender) Send(ctx context.Context, request BaseRequest) (interface{}, error) {
	handler, exists := s.container.resolve(request)
	if !exists {
		return nil, fmt.Errorf("no handlers for request %T", request)
	}
	var requestHandlerBehavior RequestHandlerFunc = func() (interface{}, error) {
		handlerMethod := reflect.ValueOf(handler).
			MethodByName("Handle")
		if !handlerMethod.IsValid() {
			return nil, fmt.Errorf("handler for request %T is not a RequestHandler", request)
		}
		// Create a slice of reflect.Value with ctx and notification as arguments
		args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(request)}

		// Call the method with ctx and notification as arguments and get the returned error
		result := handlerMethod.Call(args)

		valueResult := result[0].Interface()
		errorResult := result[1].Interface()
		if errorResult != nil {
			return valueResult, errorResult.(error)
		}
		return valueResult, nil
	}

	return s.container.executeWithPipeline(ctx, request, requestHandlerBehavior)
}
