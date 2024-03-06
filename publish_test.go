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

func (h *TestNotificationHandler) Handle(ctx context.Context, notification TestNotification) error {
	h.Executed = true
	return nil
}

type TestNotificationHandler2 struct {
	Executed bool
}

func (h *TestNotificationHandler2) Handle(ctx context.Context, notification TestNotification) error {
	h.Executed = true
	return nil
}

func TestPublish(t *testing.T) {
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

}
