package mediator

import "context"

// Notification is a marker interface for notifications
// A notification is a message that is sent to multiple handlers
type Notification interface {
}

// NotificationHandler is a marker interface for notification handlers
// A notification handler is a handler that handles a notification
type NotificationHandler[TNotification Notification] interface {
	Handle(ctx context.Context, notification TNotification) error
}
