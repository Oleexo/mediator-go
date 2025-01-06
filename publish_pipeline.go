package mediator

import (
	"context"
)

// RecoverStrategyPipelineBehavior is a type that provides panic recovery behavior within a strategy pipeline execution flow.
type RecoverStrategyPipelineBehavior struct {
}

func (r RecoverStrategyPipelineBehavior) Handle(ctx context.Context, notification Notification, next StrategyHandlerFunc) (retError error) {
	defer func() {
		if err := recover(); err != nil {
			retError = err.(error)
		}
	}()
	return next()
}

// NewRecoverStrategyPipelineBehavior creates and returns a StrategyPipelineBehavior to recover from panics during execution.
func NewRecoverStrategyPipelineBehavior() RecoverStrategyPipelineBehavior {
	return RecoverStrategyPipelineBehavior{}
}
