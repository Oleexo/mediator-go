package mediator

import (
	"context"
	"log/slog"
	"time"
)

// RecoverRequestPipelineBehavior is a middleware type that recovers from panics during request handling in a pipeline.
type RecoverRequestPipelineBehavior struct {
}

func (r RecoverRequestPipelineBehavior) Handle(_ context.Context, _ BaseRequest, next RequestHandlerFunc) (response interface{}, retError error) {
	defer func() {
		if err := recover(); err != nil {
			retError = err.(error)
			response = nil
		}
	}()
	return next()
}

// NewRecoverRequestPipelineBehavior creates and returns a new instance of RecoverRequestPipelineBehavior middleware.
func NewRecoverRequestPipelineBehavior() RecoverRequestPipelineBehavior {
	return RecoverRequestPipelineBehavior{}
}

// LogRequestPipelineBehaviorFunc defines a function signature for logging request details during pipeline processing.
// It receives a context, the request, the elapsed time, and an error if one occurred during processing.
type LogRequestPipelineBehaviorFunc func(
	ctx context.Context,
	param LogRequestParameters,
)

type LogRequestParameters struct {
	Request BaseRequest
	Elapsed time.Duration
	Error   error
}

// LogRequestPipelineBehavior is a middleware for logging request processing duration and errors in a pipeline.
// It uses a provided LogRequestPipelineBehaviorFunc to handle logging logic during request processing.
type LogRequestPipelineBehavior struct {
	LogFunc LogRequestPipelineBehaviorFunc
}

// NewLogRequestPipelineBehavior creates a new LogRequestPipelineBehavior with the provided logging function.
func NewLogRequestPipelineBehavior(logFunc LogRequestPipelineBehaviorFunc) LogRequestPipelineBehavior {
	return LogRequestPipelineBehavior{
		LogFunc: logFunc,
	}
}

// NewSlogRequestPipelineBehavior creates a LogRequestPipelineBehavior with logging for request processing and errors using slog package.
func NewSlogRequestPipelineBehavior() LogRequestPipelineBehavior {
	return NewLogRequestPipelineBehavior(func(ctx context.Context, param LogRequestParameters) {
		if param.Error != nil {
			slog.Error("Processing request failed",
				"request", param.Request.String(),
				"elapsed", param.Elapsed.String(),
				"error", param.Error,
			)
		} else {
			slog.Info("Processing request succeeded",
				"request", param.Request.String(),
				"elapsed", param.Elapsed.String(),
			)
		}
	})
}

func (l LogRequestPipelineBehavior) Handle(ctx context.Context, request BaseRequest, next RequestHandlerFunc) (interface{}, error) {
	startedAt := time.Now()

	response, err := next()

	if l.LogFunc != nil {
		elapsed := time.Since(startedAt)
		l.LogFunc(ctx, LogRequestParameters{
			Request: request,
			Elapsed: elapsed,
			Error:   err,
		})
	}

	return response, err
}
