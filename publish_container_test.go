package mediator_test

import (
	"context"
	"github.com/Oleexo/mediator-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestNotification struct {
	Value string
}

type TestNotificationHandler struct {
	Executed bool
}

func (h *TestNotificationHandler) Handle(_ context.Context, _ TestNotification) error {
	h.Executed = true
	return nil
}

type TestNotificationHandler2 struct {
	Executed bool
}

func (h *TestNotificationHandler2) Handle(_ context.Context, _ TestNotification) error {
	h.Executed = true
	return nil
}

type TestNotificationPanicHandler struct {
	Error error
}

func (h *TestNotificationPanicHandler) Handle(_ context.Context, _ TestNotification) error {
	panic(h.Error)
}

type UselessNotificationPipelineBehavior struct {
	Executed bool
	Count    int
}

func (u *UselessNotificationPipelineBehavior) Handle(ctx context.Context,
	notification mediator.Notification,
	handler any,
	next mediator.NotificationHandlerFunc) error {
	u.Executed = true
	u.Count++
	return next(ctx, notification, handler)
}

type UselessStrategyPipelineBehavior struct {
	Executed bool
}

func (u *UselessStrategyPipelineBehavior) Handle(_ context.Context, _ []mediator.Notification, next mediator.StrategyHandlerFunc) error {
	u.Executed = true
	return next()
}

func TestPublishContainer(t *testing.T) {
	t.Run("With no handlers", func(t *testing.T) {
		container := mediator.NewPublishContainer()
		notif := TestNotification{Value: "test"}

		err := mediator.PublishWithoutContext(container, notif)
		assert.NoError(t, err)
	})

	t.Run("With one handler", func(t *testing.T) {
		handler := &TestNotificationHandler{}
		container := mediator.NewPublishContainer(
			mediator.WithNotificationDefinitionHandler(mediator.NewNotificationHandlerDefinition[TestNotification](handler)),
		)
		notif := TestNotification{Value: "test"}

		err := mediator.PublishWithoutContext(container, notif)
		assert.NoError(t, err)
		assert.True(t, handler.Executed)
	})

	t.Run("With multiple handlers", func(t *testing.T) {
		handler := &TestNotificationHandler{}
		handler2 := &TestNotificationHandler2{}
		container := mediator.NewPublishContainer(
			mediator.WithNotificationDefinitionHandlers(
				mediator.NewNotificationHandlerDefinition[TestNotification](handler),
				mediator.NewNotificationHandlerDefinition[TestNotification](handler2),
			),
		)
		notif := TestNotification{Value: "test"}

		err := mediator.PublishWithoutContext(container, notif)
		assert.NoError(t, err)
		assert.True(t, handler.Executed)
		assert.True(t, handler2.Executed)
	})

	t.Run("With notification pipeline behavior", func(t *testing.T) {
		handler := &TestNotificationHandler{}
		pipeline := &UselessNotificationPipelineBehavior{}
		container := mediator.NewPublishContainer(
			mediator.WithNotificationDefinitionHandlers(
				mediator.NewNotificationHandlerDefinition[TestNotification](handler),
			),
			mediator.WithNotificationPipelineBehavior(pipeline),
		)

		notif := TestNotification{Value: "test"}

		err := mediator.PublishWithoutContext(container, notif)
		assert.NoError(t, err)
		assert.True(t, pipeline.Executed)
	})

	t.Run("With notification pipeline behavior and multiple handlers", func(t *testing.T) {
		handler := &TestNotificationHandler{}
		handler2 := &TestNotificationHandler2{}
		pipeline := &UselessNotificationPipelineBehavior{}
		container := mediator.NewPublishContainer(
			mediator.WithNotificationDefinitionHandlers(
				mediator.NewNotificationHandlerDefinition[TestNotification](handler),
				mediator.NewNotificationHandlerDefinition[TestNotification](handler2),
			),
			mediator.WithNotificationPipelineBehavior(pipeline),
		)

		notif := TestNotification{Value: "test"}

		err := mediator.PublishWithoutContext(container, notif)
		assert.NoError(t, err)
		assert.True(t, pipeline.Executed)
		assert.Equal(t, 2, pipeline.Count)
	})

	t.Run("With strategy pipeline behavior", func(t *testing.T) {
		handler := &TestNotificationHandler{}
		pipeline := &UselessStrategyPipelineBehavior{}
		container := mediator.NewPublishContainer(
			mediator.WithNotificationDefinitionHandlers(
				mediator.NewNotificationHandlerDefinition[TestNotification](handler),
			),
			mediator.WithStrategyPipelineBehavior(pipeline),
		)

		notif := TestNotification{Value: "test"}

		err := mediator.PublishWithoutContext(container, notif)
		assert.NoError(t, err)
		assert.True(t, pipeline.Executed)
	})
}
