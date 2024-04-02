package mediator

import (
	"context"
)

type PublishOptions struct {
	NotificationDefinitionHandlers []NotificationHandlerDefinition
	PublishStrategy                PublishStrategy
}

// WithNotificationDefinitionHandler adds a notification handler to the container
func WithNotificationDefinitionHandler(notificationHandler NotificationHandlerDefinition) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.NotificationDefinitionHandlers = append(options.NotificationDefinitionHandlers, notificationHandler)
	}
}

// WithNotificationDefinitionHandlers adds notification handlers to the container
func WithNotificationDefinitionHandlers(notificationHandlers ...NotificationHandlerDefinition) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.NotificationDefinitionHandlers = append(options.NotificationDefinitionHandlers, notificationHandlers...)
	}
}

// WithPublishStrategy sets the strategy to publish notifications
func WithPublishStrategy(strategy PublishStrategy) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.PublishStrategy = strategy
	}
}

// Publisher is the interface to publish notifications
type Publisher interface {
	Publish(ctx context.Context, notification interface{}) error
}

type LaunchHandler func(context.Context, interface{}) error

// PublishStrategy is the strategy to publish notifications
type PublishStrategy interface {
	Execute(ctx context.Context,
		handlers []interface{},
		launcher LaunchHandler) error
}
