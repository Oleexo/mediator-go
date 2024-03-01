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

func AsNotificationHandler[TNotification mediator.Notification](handler mediator.NotificationHandler[TNotification]) any {
	definition := mediator.NewNotificationHandlerDefinition(handler)
	return fx.Annotate(definition,
		fx.As(new(mediator.NotificationHandlerDefinition)),
		fx.ResultTags(`group:"mediator_notification_handlers"`),
	)
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

type ContainerParams struct {
	fx.In

	RequestHandlers      []mediator.RequestHandlerDefinition      `group:"mediator_request_handlers"`
	NotificationHandlers []mediator.NotificationHandlerDefinition `group:"mediator_notification_handlers"`
	Pipelines            []mediator.PipelineBehavior              `group:"mediator_pipelines"`
}

func New(param ContainerParams) mediator.Container {
	return mediator.New(mediator.WithRequestDefinitionHandlers(param.RequestHandlers),
		mediator.WithNotificationDefinitionHandlers(param.NotificationHandlers),
		mediator.WithPipelineBehaviors(param.Pipelines))
}
