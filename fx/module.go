package mediatorfx

import (
	"github.com/Oleexo/mediator-go"
	"github.com/Oleexo/mediator-go/pipelines"
	"go.uber.org/fx"
)

func AsRequestHandler[TRequest mediator.Request[TResponse], TResponse interface{}](f any) []interface{} {
	return []interface{}{
		fx.Annotate(
			f,
			fx.As(new(mediator.RequestHandler[TRequest, TResponse])),
		),
		fx.Annotate(
			mediator.NewRequestHandlerDefinition[TRequest, TResponse],
			fx.As(new(mediator.RequestHandlerDefinition)),
			fx.ResultTags(`group:"mediator_request_handlers"`),
		),
	}
}

func AsNotificationHandler[TNotification mediator.Notification, TNotificationHandler mediator.NotificationHandler[TNotification]](f any) []interface{} {
	return []interface{}{
		f,
		fx.Annotate(mediator.NewNotificationHandlerDefinition[TNotification],
			fx.As(new(mediator.NotificationHandlerDefinition)),
			fx.From(new(TNotificationHandler)),
			fx.ResultTags(`group:"mediator_notification_handlers"`),
		),
	}
}

func AsPipelineBehavior(f any) interface{} {
	return fx.Annotate(f,
		fx.As(new(mediator.PipelineBehavior)),
		fx.ResultTags(`group:"mediator_pipelines"`),
	)
}

func AddValidationPipeline(optFns ...func(options *pipelines.Options)) interface{} {
	return AsPipelineBehavior(func() *pipelines.ValidationPipeline {
		return pipelines.NewValidationPipeline(optFns...)
	})
}

type SendContainerParams struct {
	fx.In

	RequestHandlers []mediator.RequestHandlerDefinition `group:"mediator_request_handlers"`
	Pipelines       []mediator.PipelineBehavior         `group:"mediator_pipelines"`
}

type PublishContainerParams struct {
	fx.In

	NotificationHandlers []mediator.NotificationHandlerDefinition `group:"mediator_notification_handlers"`
}

func NewSend(param SendContainerParams) mediator.SendContainer {
	return mediator.NewSendContainer(mediator.WithRequestDefinitionHandlers(param.RequestHandlers),
		mediator.WithPipelineBehaviors(param.Pipelines))
}

func NewNotification(param PublishContainerParams) mediator.NotificationContainer {
	return mediator.NewNotificationContainer(mediator.WithNotificationDefinitionHandlers(param.NotificationHandlers))
}

func NewModule() fx.Option {
	return fx.Module("mediatorfx",
		fx.Provide(NewSend),
		fx.Provide(NewNotification),
	)
}
