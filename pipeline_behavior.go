package mediator

import "context"

// RequestHandlerFunc is a function that handles a request
type RequestHandlerFunc func() (interface{}, error)

// PipelineBehavior is a marker interface for pipeline behaviors
// A pipeline behavior is a behavior that is executed as part of a pipeline
type PipelineBehavior interface {
	Handle(ctx context.Context, request BaseRequest, next RequestHandlerFunc) (interface{}, error)
}
