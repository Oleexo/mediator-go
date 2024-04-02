package mediator

import (
	"context"
	"reflect"
)

// SendContainer is the mediator container for request and notification handlers
// It is responsible for resolving handlers and pipeline behaviors
type SendContainer interface {
	resolve(request interface{}) (interface{}, bool)
	executeWithPipeline(ctx context.Context,
		request BaseRequest,
		requestHandlerBehavior RequestHandlerFunc) (interface{}, error)
}

type sendContainer struct {
	requestHandlers map[reflect.Type]interface{}
	pipelines       []PipelineBehavior
}

func (c sendContainer) resolve(request interface{}) (interface{}, bool) {
	handler, ok := c.requestHandlers[reflect.TypeOf(request)]
	return handler, ok
}

func (c sendContainer) executeWithPipeline(ctx context.Context,
	request BaseRequest,
	requestHandlerBehavior RequestHandlerFunc) (interface{}, error) {
	if len(c.pipelines) > 0 {
		var reversPipes = reversOrder(c.pipelines)

		v := buildPipeline(reversPipes, requestHandlerBehavior,
			func(next RequestHandlerFunc, pipe PipelineBehavior) RequestHandlerFunc {
				pipeValue := pipe
				nexValue := next

				var handlerFunc RequestHandlerFunc = func() (interface{}, error) {
					return pipeValue.Handle(ctx, request, nexValue)
				}

				return handlerFunc
			})

		response, err := v()

		if err != nil {
			return response, err
		}

		return response, nil
	} else {
		return requestHandlerBehavior()
	}
}

type SendContainerOptions struct {
	RequestDefinitionHandlers []RequestHandlerDefinition
	PipelineBehaviors         []PipelineBehavior
}

// WithRequestDefinitionHandler adds a request handler to the container
func WithRequestDefinitionHandler(requestHandler RequestHandlerDefinition) func(*SendContainerOptions) {
	return func(options *SendContainerOptions) {
		options.RequestDefinitionHandlers = append(options.RequestDefinitionHandlers, requestHandler)
	}
}

// WithRequestDefinitionHandlers adds request handlers to the container
func WithRequestDefinitionHandlers(requestHandlers ...RequestHandlerDefinition) func(*SendContainerOptions) {
	return func(options *SendContainerOptions) {
		options.RequestDefinitionHandlers = append(options.RequestDefinitionHandlers, requestHandlers...)
	}
}

// WithPipelineBehavior adds a pipeline behavior to the container
func WithPipelineBehavior(pipelineBehavior PipelineBehavior) func(*SendContainerOptions) {
	return func(options *SendContainerOptions) {
		options.PipelineBehaviors = append(options.PipelineBehaviors, pipelineBehavior)
	}
}

// WithPipelineBehaviors adds pipeline behaviors to the container
func WithPipelineBehaviors(pipelineBehaviors []PipelineBehavior) func(*SendContainerOptions) {
	return func(options *SendContainerOptions) {
		options.PipelineBehaviors = append(options.PipelineBehaviors, pipelineBehaviors...)
	}
}

func NewSendContainer(optFns ...func(*SendContainerOptions)) SendContainer {
	options := &SendContainerOptions{}
	for _, optFn := range optFns {
		optFn(options)
	}
	requestDefinitionHandlers := options.RequestDefinitionHandlers
	requestHandlers := make(map[reflect.Type]interface{}, len(requestDefinitionHandlers))
	for _, requestHandler := range requestDefinitionHandlers {
		requestHandlers[requestHandler.RequestType()] = requestHandler.Handler()
	}
	return sendContainer{
		requestHandlers: requestHandlers,
		pipelines:       options.PipelineBehaviors,
	}
}
