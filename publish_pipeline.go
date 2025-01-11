package mediator

import (
	"context"
	"log/slog"
	"time"
)

// RecoverStrategyPipelineBehavior is a type that provides panic recovery behavior within a strategy pipeline execution flow.
type RecoverStrategyPipelineBehavior struct {
}

func (r RecoverStrategyPipelineBehavior) Handle(_ context.Context, _ []Notification, next StrategyHandlerFunc) (retError error) {
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

// LogNotificationPipelineBehaviorFunc defines a function to log notifications within a pipeline's behavior context.
// It accepts a context and LogNotificationParameters holding details about the notification and its processing.
type LogNotificationPipelineBehaviorFunc func(
	ctx context.Context,
	param LogNotificationParameters,
)

// LogNotificationParameters holds data about a single notification and its processing state in the pipeline.
// It includes the notification instance, the handler that processed it, elapsed time, and any occurred error.
type LogNotificationParameters struct {
	Notification Notification
	Handler      any
	Elapsed      time.Duration
	Error        error
}

// LogNotificationPipelineBehavior is a pipeline behavior for logging notification handling processes.
// It uses a user-defined LogFunc to log details such as elapsed time and potential errors.
type LogNotificationPipelineBehavior struct {
	LogFunc LogNotificationPipelineBehaviorFunc
}

// NewLogNotificationPipelineBehavior creates a new instance of LogNotificationPipelineBehavior with the provided LogFunc.
func NewLogNotificationPipelineBehavior(logFunc LogNotificationPipelineBehaviorFunc) LogNotificationPipelineBehavior {
	return LogNotificationPipelineBehavior{
		LogFunc: logFunc,
	}
}

// NewSlogNotificationPipelineBehavior creates a LogNotificationPipelineBehavior to log notification process details.
func NewSlogNotificationPipelineBehavior() LogNotificationPipelineBehavior {
	return LogNotificationPipelineBehavior{
		LogFunc: func(ctx context.Context, param LogNotificationParameters) {
			if err := param.Error; err != nil {
				slog.Error("Notification process failed",
					"notification", param.Notification,
					"handler", param.Handler,
					"elapsed", param.Elapsed,
					"error", err,
				)
			} else {
				slog.Info("Notification processed successfully",
					"notification", param.Notification,
					"handler", param.Handler,
					"elapsed", param.Elapsed,
				)
			}
		},
	}
}

func (l LogNotificationPipelineBehavior) Handle(ctx context.Context, notification Notification, handler any, next NotificationHandlerFunc) error {
	startedAt := time.Now()

	err := next(ctx, notification, handler)

	if l.LogFunc != nil {
		elapsed := time.Since(startedAt)
		l.LogFunc(ctx, LogNotificationParameters{
			Notification: notification,
			Handler:      handler,
			Elapsed:      elapsed,
			Error:        err,
		})
	}

	return err
}

// LogStrategyPipelineBehaviorFunc defines a function type for handling logging of strategy pipeline behaviors.
// It processes context, parameters, and details of the strategy execution within a pipeline.
type LogStrategyPipelineBehaviorFunc func(
	ctx context.Context,
	param LogStrategyParameters,
)

// LogStrategyParameters encapsulates details for logging strategy pipeline behavior, including notifications, handlers, elapsed time, and errors.
type LogStrategyParameters struct {
	Notification Notification
	Handlers     []any
	Elapsed      time.Duration
	Error        error
}

// LogStrategyPipelineBehavior defines a structure for handling and logging strategy processing behavior in pipelines.
// It wraps a logging function of type LogStrategyPipelineBehaviorFunc to log details of the strategy execution.
type LogStrategyPipelineBehavior struct {
	LogFunc LogStrategyPipelineBehaviorFunc
}

// NewLogStrategyPipelineBehavior creates a LogStrategyPipelineBehavior using the provided logging function.
func NewLogStrategyPipelineBehavior(logFunc LogStrategyPipelineBehaviorFunc) LogStrategyPipelineBehavior {
	return LogStrategyPipelineBehavior{
		LogFunc: logFunc,
	}
}

// NewSlogStrategyPipelineBehavior creates a LogStrategyPipelineBehavior with a logging function for strategy operations.
func NewSlogStrategyPipelineBehavior() LogStrategyPipelineBehavior {
	return LogStrategyPipelineBehavior{
		LogFunc: func(ctx context.Context, param LogStrategyParameters) {
			if err := param.Error; err != nil {
				slog.Error("Strategy process failed",
					"notification", param.Notification,
					"handlers", param.Handlers,
					"elapsed", param.Elapsed,
					"error", err,
				)
			} else {
				slog.Info("Strategy processed successfully",
					"notification", param.Notification,
					"handlers", param.Handlers,
					"elapsed", param.Elapsed,
				)
			}
		},
	}
}

func (l LogStrategyPipelineBehavior) Handle(ctx context.Context, notification Notification, handlers []any, next StrategyHandlerFunc) error {
	startedAt := time.Now()

	err := next()

	if l.LogFunc != nil {
		elapsed := time.Since(startedAt)
		l.LogFunc(ctx, LogStrategyParameters{
			Notification: notification,
			Handlers:     handlers,
			Elapsed:      elapsed,
			Error:        err,
		})
	}

	return err
}
