package mediator

import (
	"context"

	"github.com/pkg/errors"
)

// PublishWithoutContext publishes a notification to multiple handlers without a context
func PublishWithoutContext[TNotification Notification](container Container, notification TNotification) error {
	return Publish[TNotification](context.Background(), container, notification)
}

// Publish publishes a notification to multiple handlers
func Publish[TNotification Notification](ctx context.Context, container Container, notification TNotification) error {
	handler, exists := container.resolve(notification)
	if !exists {
		return nil
	}
	handlerValue, ok := handler.(NotificationHandler[TNotification])
	if !ok {
		return errors.Errorf("handler for notification %T is not a NotificationHandler",
			notification)
	}
	return handlerValue.Handle(ctx, notification)
}
