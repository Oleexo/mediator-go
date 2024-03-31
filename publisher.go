package mediator

import (
	"context"
	"github.com/pkg/errors"
	"reflect"
)

type publisher struct {
	notificationHandlers map[reflect.Type][]interface{}
	strategy             PublishStrategy
}

func NewPublisher(optFns ...func(*PublishOptions)) Publisher {
	options := &PublishOptions{}
	for _, optFn := range optFns {
		optFn(options)
	}
	notificationDefinitionHandlers := options.NotificationDefinitionHandlers
	notificationHandlers := make(map[reflect.Type][]interface{}, len(notificationDefinitionHandlers))
	for _, notificationHandler := range notificationDefinitionHandlers {
		if handlers, ok := notificationHandlers[notificationHandler.NotificationType()]; ok {
			handlers = append(handlers, notificationHandler.Handler())
			notificationHandlers[notificationHandler.NotificationType()] = handlers
		} else {
			notificationHandlers[notificationHandler.NotificationType()] = []interface{}{notificationHandler.Handler()}
		}
	}
	strategy := options.PublishStrategy
	if strategy == nil {
		strategy = NewSynchronousPublishStrategy()
	}
	return &publisher{
		notificationHandlers: notificationHandlers,
		strategy:             strategy,
	}

}

func (s publisher) Publish(ctx context.Context, notification interface{}) error {
	handlers := resolve(notification, s.notificationHandlers)
	if handlers == nil {
		return nil
	}

	return s.strategy.Execute(ctx, handlers, func(handlerCtx context.Context, handler interface{}) error {
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
