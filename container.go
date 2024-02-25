package mediator

import (
	"reflect"
)

// Container is the mediator container for request and notification handlers
// It is responsible for resolving handlers and pipeline behaviors
type Container interface {
	resolve(request interface{}) (interface{}, bool)
	pipelineBehaviors() []PipelineBehavior
}

type container struct {
	requestHandlers      map[reflect.Type]interface{}
	notificationHandlers map[reflect.Type]interface{}
	pipelines            []PipelineBehavior
}

func (c container) pipelineBehaviors() []PipelineBehavior {
	return c.pipelines
}

func (c container) resolve(request interface{}) (interface{}, bool) {
	handler, ok := c.requestHandlers[reflect.TypeOf(request)]
	return handler, ok
}

type ContainerOptions struct {
	RequestDefinitionHandlers      []RequestHandlerDefinition
	NotificationDefinitionHandlers []NotificationHandlerDefinition
	PipelineBehaviors              []PipelineBehavior
}

// WithRequestDefinitionHandler adds a request handler to the container
func WithRequestDefinitionHandler(requestHandler RequestHandlerDefinition) func(*ContainerOptions) {
	return func(options *ContainerOptions) {
		options.RequestDefinitionHandlers = append(options.RequestDefinitionHandlers, requestHandler)
	}
}

// WithRequestDefinitionHandlers adds request handlers to the container
func WithRequestDefinitionHandlers(requestHandlers []RequestHandlerDefinition) func(*ContainerOptions) {
	return func(options *ContainerOptions) {
		options.RequestDefinitionHandlers = append(options.RequestDefinitionHandlers, requestHandlers...)
	}
}

// WithNotificationDefinitionHandler adds a notification handler to the container
func WithNotificationDefinitionHandler(notificationHandler NotificationHandlerDefinition) func(*ContainerOptions) {
	return func(options *ContainerOptions) {
		options.NotificationDefinitionHandlers = append(options.NotificationDefinitionHandlers, notificationHandler)
	}
}

// WithNotificationDefinitionHandlers adds notification handlers to the container
func WithNotificationDefinitionHandlers(notificationHandlers []NotificationHandlerDefinition) func(*ContainerOptions) {
	return func(options *ContainerOptions) {
		options.NotificationDefinitionHandlers = append(options.NotificationDefinitionHandlers, notificationHandlers...)
	}
}

// WithPipelineBehavior adds a pipeline behavior to the container
func WithPipelineBehavior(pipelineBehavior PipelineBehavior) func(*ContainerOptions) {
	return func(options *ContainerOptions) {
		options.PipelineBehaviors = append(options.PipelineBehaviors, pipelineBehavior)
	}
}

// WithPipelineBehaviors adds pipeline behaviors to the container
func WithPipelineBehaviors(pipelineBehaviors []PipelineBehavior) func(*ContainerOptions) {
	return func(options *ContainerOptions) {
		options.PipelineBehaviors = append(options.PipelineBehaviors, pipelineBehaviors...)
	}
}

// New creates a new mediator container
func New(optFns ...func(*ContainerOptions)) Container {
	options := &ContainerOptions{}
	for _, optFn := range optFns {
		optFn(options)
	}
	requestDefinitionHandlers := options.RequestDefinitionHandlers
	requestHandlers := make(map[reflect.Type]interface{}, len(requestDefinitionHandlers))
	for _, requestHandler := range requestDefinitionHandlers {
		requestHandlers[requestHandler.RequestType()] = requestHandler.Handler()
	}
	notificationDefinitionHandlers := options.NotificationDefinitionHandlers
	notificationHandlers := make(map[reflect.Type]interface{}, len(notificationDefinitionHandlers))
	for _, notificationHandler := range notificationDefinitionHandlers {
		notificationHandlers[notificationHandler.NotificationType()] = notificationHandler.Handler()
	}
	return &container{
		requestHandlers:      requestHandlers,
		notificationHandlers: notificationHandlers,
		pipelines:            options.PipelineBehaviors,
	}
}
