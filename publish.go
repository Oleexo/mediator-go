package mediator

import (
	"context"
	"fmt"
	"reflect"
)

type PublishOptions struct {
	NotificationDefinitionHandlers []NotificationHandlerDefinition
	Pipelines                      []NotificationPipelineBehavior
	Strategy                       PublishStrategy
	StrategyPipelines              []StrategyPipelineBehavior
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
		options.Strategy = strategy
	}
}

// WithNotificationPipelineBehavior adds a notification pipeline behavior to the publish options.
func WithNotificationPipelineBehavior(pipelineBehavior NotificationPipelineBehavior) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.Pipelines = append(options.Pipelines, pipelineBehavior)
	}
}

// WithStrategyPipelineBehavior adds a strategy pipeline behavior to the publish options configuration.
func WithStrategyPipelineBehavior(pipelineBehavior StrategyPipelineBehavior) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.StrategyPipelines = append(options.StrategyPipelines, pipelineBehavior)
	}
}

// Publisher is the interface to publish notifications
type Publisher interface {
	Publish(ctx context.Context, notification interface{}) error
}

// NotificationHandlerFunc is a function type for handling notifications with a given context and handler object.
type NotificationHandlerFunc func(ctx context.Context, handler any) error

// NotificationPipelineBehavior defines a pipeline step for processing notifications before reaching their handlers.
// The Handle method processes a notification with a potentially custom behavior and delegates to the next pipeline step.
type NotificationPipelineBehavior interface {
	Handle(ctx context.Context, notification Notification, handler any, next NotificationHandlerFunc) error
}

// StrategyHandlerFunc represents a function used as part of a strategy pipeline, returning an error if one occurs.
type StrategyHandlerFunc func() error

// StrategyPipelineBehavior defines a behavior interface for handling strategies in a pipeline configuration.
// It processes a notification alongside a contextual handler function sequence.
type StrategyPipelineBehavior interface {
	Handle(ctx context.Context, notification Notification, handlers []any, next StrategyHandlerFunc) error
}

// PublishStrategy is the strategy to publish notifications
type PublishStrategy interface {
	Execute(ctx context.Context,
		handlers []interface{},
		launcher NotificationHandlerFunc) error
}

type publisher struct {
	container PublishContainer
}

// NewPublisher creates a Publisher instance using the provided PublishContainer to resolve notification handlers.
func NewPublisher(container PublishContainer) Publisher {
	return &publisher{
		container: container,
	}
}

// Publish sends a notification to all registered handlers in the container and executes their Handle method.
// Returns an error if resolving or executing handlers fails.
func (s publisher) Publish(ctx context.Context, notification interface{}) error {
	return s.container.execute(ctx, notification, func(ctx context.Context, handler any) error {
		handlerMethod := reflect.ValueOf(handler).
			MethodByName("Handle")
		if !handlerMethod.IsValid() {
			return fmt.Errorf("handler for notification %T is not a NotificationHandler", notification)
		}
		// Create a slice of reflect.Value with ctx and notification as arguments
		args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(notification)}

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
