package mediator

import "reflect"

// PublishContainer is the mediator container for request and notification handlers
// It is responsible for resolving handlers and pipeline behaviors
type PublishContainer interface {
	resolve(request interface{}) []interface{}
}

type notificationContainer struct {
	notificationHandlers map[reflect.Type][]interface{}
}

func (c notificationContainer) resolve(request interface{}) []interface{} {
	handlers, ok := c.notificationHandlers[reflect.TypeOf(request)]
	if !ok {
		return nil
	}
	return handlers
}

type PublishContainerOptions struct {
	NotificationDefinitionHandlers []NotificationHandlerDefinition
}

// WithNotificationDefinitionHandler adds a notification handler to the container
func WithNotificationDefinitionHandler(notificationHandler NotificationHandlerDefinition) func(*PublishContainerOptions) {
	return func(options *PublishContainerOptions) {
		options.NotificationDefinitionHandlers = append(options.NotificationDefinitionHandlers, notificationHandler)
	}
}

// WithNotificationDefinitionHandlers adds notification handlers to the container
func WithNotificationDefinitionHandlers(notificationHandlers ...NotificationHandlerDefinition) func(*PublishContainerOptions) {
	return func(options *PublishContainerOptions) {
		options.NotificationDefinitionHandlers = append(options.NotificationDefinitionHandlers, notificationHandlers...)
	}
}

func NewPublishContainer(optFns ...func(*PublishContainerOptions)) PublishContainer {
	options := &PublishContainerOptions{}
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
	return &notificationContainer{
		notificationHandlers: notificationHandlers,
	}
}
