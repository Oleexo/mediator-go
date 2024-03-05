package mediator

import (
	"context"

	"github.com/pkg/errors"
)

// PublishWithoutContext publishes a notification to multiple handlers without a context
func PublishWithoutContext[TNotification Notification](container NotificationContainer, notification TNotification) error {
	return Publish[TNotification](context.Background(), container, notification)
}

// Publish publishes a notification to multiple handlers
func Publish[TNotification Notification](ctx context.Context, container NotificationContainer, notification TNotification) error {
	handlers := container.resolve(notification)
	if handlers == nil {
		return nil
	}
	for _, handler := range handlers {
		handlerValue, ok := handler.(NotificationHandler[TNotification])
		if !ok {
			return errors.Errorf("handler for notification %T is not a NotificationHandler", notification)
		}
		err := handlerValue.Handle(ctx, notification)
		if err != nil {
			return err
		}
	}
	return nil
}
