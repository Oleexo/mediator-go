package mediator_test

import (
	"context"
	"errors"
	"github.com/Oleexo/mediator-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

type synchronousHandler interface {
	Execute() error
	Executed() bool
}

type synchronousHandlerImpl struct {
	executed bool
	err      error
}

func newSynchronousHandler(err error) synchronousHandler {
	return &synchronousHandlerImpl{
		executed: false,
		err:      err,
	}
}

func (s *synchronousHandlerImpl) Execute() error {
	s.executed = true
	return s.err
}

func (s *synchronousHandlerImpl) Executed() bool {
	return s.executed
}

func TestSynchronousPublishStrategy(t *testing.T) {
	t.Run("no error will run all tests", func(t *testing.T) {
		strategy := mediator.NewSynchronousPublishStrategy()
		handlers := []interface{}{
			newSynchronousHandler(nil),
			newSynchronousHandler(nil),
		}

		notification := TestNotification{Value: "test"}

		result := strategy.Execute(context.Background(),
			[]mediator.Notification{notification},
			func(notification mediator.Notification) []any {
				return handlers
			},
			func(ctx context.Context, notification mediator.Notification, handler interface{}) error {
				return handler.(synchronousHandler).Execute()
			})

		assert.NoError(t, result)
		for _, handler := range handlers {
			assert.True(t, handler.(synchronousHandler).Executed())
		}
	})

	t.Run("error will stop the execution", func(t *testing.T) {
		strategy := mediator.NewSynchronousPublishStrategy()
		handlers := []interface{}{
			newSynchronousHandler(nil),
			newSynchronousHandler(errors.New("an error happened")),
			newSynchronousHandler(nil),
		}
		notification := TestNotification{Value: "test"}

		result := strategy.Execute(context.Background(),
			[]mediator.Notification{notification},
			func(notification mediator.Notification) []any {
				return handlers
			},
			func(ctx context.Context, notification mediator.Notification, handler interface{}) error {
				return handler.(synchronousHandler).Execute()
			})

		assert.Error(t, result)
		assert.True(t, handlers[0].(synchronousHandler).Executed())
		assert.True(t, handlers[1].(synchronousHandler).Executed())
		assert.False(t, handlers[2].(synchronousHandler).Executed())
	})
}
