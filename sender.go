package mediator

import (
	"context"
	"fmt"
	"reflect"
)

type sender struct {
	container SendContainer
}

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
