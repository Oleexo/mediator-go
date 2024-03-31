package mediator_test

import (
	"context"
	"github.com/Oleexo/mediator-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPublisher(t *testing.T) {
	t.Run("With no handlers", func(t *testing.T) {
		publisher := mediator.NewPublisher()
		notif := TestNotification{Value: "test"}

		err := publisher.Publish(context.Background(), notif)
		assert.NoError(t, err)
	})

	t.Run("With one handler", func(t *testing.T) {
		handler := &TestNotificationHandler{}
		publisher := mediator.NewPublisher(
			mediator.WithNotificationDefinitionHandler(mediator.NewNotificationHandlerDefinition[TestNotification](handler)),
		)
		notif := TestNotification{Value: "test"}

		err := publisher.Publish(context.Background(), notif)
		assert.NoError(t, err)
		assert.True(t, handler.Executed)
	})

	t.Run("With multiple handlers", func(t *testing.T) {
		handler := &TestNotificationHandler{}
		handler2 := &TestNotificationHandler2{}
		publisher := mediator.NewPublisher(
			mediator.WithNotificationDefinitionHandlers(
				mediator.NewNotificationHandlerDefinition[TestNotification](handler),
				mediator.NewNotificationHandlerDefinition[TestNotification](handler2),
			),
		)
		notif := TestNotification{Value: "test"}

		err := publisher.Publish(context.Background(), notif)
		assert.NoError(t, err)
		assert.True(t, handler.Executed)
		assert.True(t, handler2.Executed)
	})

}
