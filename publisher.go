package mediator

import (
	"context"
	"github.com/pkg/errors"
	"reflect"
)

type publisher struct {
	container PublishContainer
}

func NewPublisher(container PublishContainer) Publisher {
	return &publisher{
		container: container,
	}

}

func (s publisher) Publish(ctx context.Context, notification interface{}) error {
	handlers := s.container.resolve(notification)
	if handlers == nil {
		return nil
	}

	return s.container.getStrategy().Execute(ctx, handlers, func(handlerCtx context.Context, handler interface{}) error {
		handlerMethod := reflect.ValueOf(handler).
			MethodByName("Handle")
		if !handlerMethod.IsValid() {
			return errors.Errorf("handler for notification %T is not a NotificationHandler", notification)
		}
		// Create a slice of reflect.Value with ctx and notification as arguments
		args := []reflect.Value{reflect.ValueOf(handlerCtx), reflect.ValueOf(notification)}

		// Call the method with ctx and notification as arguments and get the returned error
		result := handlerMethod.Call(args)

		// The first (and only) return value of the Handle method is the error
		methodResult := result[0].Interface()
		if methodResult != nil {
			return methodResult.(error)
		}
		return nil
	})
}
