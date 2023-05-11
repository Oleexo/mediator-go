package mediator

import "context"

type Notification interface {
}

type NotificationHandler[TNotification Notification] interface {
	Handle(ctx context.Context, request TNotification) error
}
