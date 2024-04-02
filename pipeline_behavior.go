package mediator

import "context"

// PipelineBehavior is a marker interface for pipeline behaviors
// A pipeline behavior is a behavior that is executed as part of a pipeline
type PipelineBehavior interface {
	Handle(ctx context.Context, request BaseRequest, next RequestHandlerFunc) (interface{}, error)
}
