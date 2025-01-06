package mediator_test

import (
	"context"
	"errors"
	"github.com/Oleexo/mediator-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecoverStrategyPipelineBehavior(t *testing.T) {
	panicErr := errors.New("test")
	handler := &TestNotificationPanicHandler{
		Error: panicErr,
	}

	container := mediator.NewPublishContainer(
		mediator.WithStrategyPipelineBehavior(
			mediator.NewRecoverStrategyPipelineBehavior(),
		),
		mediator.WithNotificationDefinitionHandler(
			mediator.NewNotificationHandlerDefinition[TestNotification](handler),
		),
	)

	publisher := mediator.NewPublisher(container)
	notif := TestNotification{Value: "test"}

	err := publisher.Publish(context.Background(), notif)
	assert.Error(t, err)
	assert.Equal(t, panicErr, err)
}
