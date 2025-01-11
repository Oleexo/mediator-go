package mediator

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
)

// PublishWithoutContext publishes a notification to multiple handlers without a context
func PublishWithoutContext[TNotification Notification](container PublishContainer, notification TNotification) error {
	return Publish[TNotification](context.Background(), container, notification)
}

// Publish publishes a notification to multiple handlers
func Publish[TNotification Notification](ctx context.Context, container PublishContainer, notification TNotification) error {
	return container.execute(ctx, []Notification{notification}, func(ctx context.Context, notification Notification, handler any) error {
		handlerValue, ok := handler.(NotificationHandler[TNotification])
		if !ok {
			return fmt.Errorf("handler for notification %T is not a NotificationHandler", notification)
		}
		n, ok := notification.(TNotification)
		if !ok {
			return errors.New("notification is not of type TNotification")
		}
		err := handlerValue.Handle(ctx, n)
		if err != nil {
			return err
		}
		return nil
	})
}

// PublishContainer is the mediator container for request and notification handlers
// It is responsible for resolving handlers and pipeline behaviors
type PublishContainer interface {
	execute(ctx context.Context, notification []Notification, seed NotificationHandlerFunc) error
}

type publishContainer struct {
	notificationHandlers map[reflect.Type][]interface{}
	strategy             PublishStrategy
	pipelines            []NotificationPipelineBehavior
	strategyPipelines    []StrategyPipelineBehavior
}

func (n publishContainer) resolve(notification interface{}) []interface{} {
	notificationType := reflect.TypeOf(notification)
	results, ok := n.notificationHandlers[notificationType]
	if !ok {
		return nil
	}
	return results
}

type Resolver func(notification Notification) []any

func (n publishContainer) execute(ctx context.Context, notifications []Notification, seed NotificationHandlerFunc) error {
	resolver := func(notification Notification) []any {
		return n.resolve(notification)
	}
	if len(n.pipelines) == 0 && len(n.strategyPipelines) == 0 {
		return n.strategy.Execute(ctx, notifications, resolver, seed)
	}

	return n.executeWithPipelines(ctx, notifications, resolver, seed)
}

func (n publishContainer) executeWithPipelines(ctx context.Context,
	notifications []Notification,
	resolver Resolver,
	seed NotificationHandlerFunc) error {
	var f = seed
	if len(n.pipelines) > 0 {
		f = buildPipeline[NotificationPipelineBehavior, NotificationHandlerFunc](n.pipelines,
			seed,
			func(next NotificationHandlerFunc, pipe NotificationPipelineBehavior) NotificationHandlerFunc {
				pipeValue := pipe
				nexValue := next

				var handlerFunc NotificationHandlerFunc = func(ctx context.Context, notification Notification, handler any) error {
					return pipeValue.Handle(ctx, notification, handler, nexValue)
				}

				return handlerFunc
			})
	}

	if len(n.strategyPipelines) == 0 {
		return n.strategy.Execute(ctx, notifications, resolver, f)
	}

	s := buildPipeline[StrategyPipelineBehavior, StrategyHandlerFunc](n.strategyPipelines,
		func() error {
			return n.strategy.Execute(ctx, notifications, resolver, f)
		},
		func(next StrategyHandlerFunc, pipe StrategyPipelineBehavior) StrategyHandlerFunc {
			pipeValue := pipe
			nexValue := next

			var handlerFunc StrategyHandlerFunc = func() error {
				return pipeValue.Handle(ctx, notifications, nexValue)
			}

			return handlerFunc
		},
	)

	return s()
}

// WithDefaultPublishOptions sets the default publish strategy and adds default strategy pipeline behaviors.
// The default strategy is Synchronous.
func WithDefaultPublishOptions() func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.Strategy = NewSynchronousPublishStrategy()
		options.StrategyPipelines = append(options.StrategyPipelines, NewRecoverStrategyPipelineBehavior())
	}
}

// NewPublishContainer creates and returns a new PublishContainer with customizable publish options and strategies.
// It accepts optional functional arguments to configure notification handlers, pipelines, and publish strategy.
func NewPublishContainer(optFns ...func(*PublishOptions)) PublishContainer {
	options := &PublishOptions{}
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
	strategy := options.Strategy
	if strategy == nil {
		strategy = NewSynchronousPublishStrategy()
	}

	pipelines := options.Pipelines
	slices.Reverse(pipelines)
	strategyPipelines := options.StrategyPipelines
	slices.Reverse(strategyPipelines)
	return &publishContainer{
		notificationHandlers: notificationHandlers,
		strategy:             strategy,
		pipelines:            pipelines,
		strategyPipelines:    strategyPipelines,
	}
}
