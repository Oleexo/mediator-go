package mediator

import "reflect"

type NotificationHandlerDefinition interface {
	NotificationType() reflect.Type
	Handler() interface{}
}

func NewNotificationHandlerDefinition[TNotification Notification](handler NotificationHandler[TNotification]) NotificationHandlerDefinition {
	var notification Notification
	notificationType := reflect.TypeOf(notification)

	return TypedNotificationHandlerDefinition[TNotification]{
		notificationType: notificationType,
		handler:          handler,
	}
}

type TypedNotificationHandlerDefinition[TNotification Notification] struct {
	notificationType reflect.Type
	handler          NotificationHandler[TNotification]
}

func (t TypedNotificationHandlerDefinition[TNotification]) NotificationType() reflect.Type {
	return t.notificationType
}

func (t TypedNotificationHandlerDefinition[TNotification]) Handler() interface{} {
	return t.handler
}
