package mediator

import (
	"context"
	"github.com/pkg/errors"
	"reflect"
)

// PublishWithoutContext publishes a notification to multiple handlers without a context
func PublishWithoutContext[TNotification Notification](container PublishContainer, notification TNotification) error {
	return Publish[TNotification](context.Background(), container, notification)
}

// Publish publishes a notification to multiple handlers
func Publish[TNotification Notification](ctx context.Context, container PublishContainer, notification TNotification) error {
	handlers := resolve(notification, container.getHandlers())
	if handlers == nil {
		return nil
	}

	return container.getStrategy().Execute(ctx, handlers, func(handlerCtx context.Context, handler interface{}) error {
		handlerValue, ok := handler.(NotificationHandler[TNotification])
		if !ok {
			return errors.Errorf("handler for notification %T is not a NotificationHandler", notification)
		}
		err := handlerValue.Handle(handlerCtx, notification)
		if err != nil {
			return err
		}
		return nil
	})
}

// PublishContainer is the mediator container for request and notification handlers
// It is responsible for resolving handlers and pipeline behaviors
type PublishContainer interface {
	getHandlers() map[reflect.Type][]interface{}
	getStrategy() PublishStrategy
}

type notificationContainer struct {
	notificationHandlers map[reflect.Type][]interface{}
	strategy             PublishStrategy
}

func (n notificationContainer) getHandlers() map[reflect.Type][]interface{} {
	return n.notificationHandlers
}

func (n notificationContainer) getStrategy() PublishStrategy {
	return n.strategy
}

func NewPublishContainer(optFns ...func(*PublishOptions)) PublishContainer {
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

	return &notificationContainer{
		notificationHandlers: notificationHandlers,
		strategy:             strategy,
	}
}
