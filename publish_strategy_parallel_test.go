package mediator_test

import (
	"context"
	"github.com/Oleexo/mediator-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type parallelHandler interface {
	Execute() error
	Executed() bool
}

type parallelHandlerImpl struct {
	executed bool
	err      error
}

func newParallelHandler(err error) parallelHandler {
	return &parallelHandlerImpl{
		executed: false,
		err:      err,
	}
}

func (s *parallelHandlerImpl) Execute() error {
	s.executed = true
	return s.err
}

func (s *parallelHandlerImpl) Executed() bool {
	return s.executed
}

func TestParallelPublishStrategy(t *testing.T) {
	t.Run("no error will run all tests", func(t *testing.T) {
		strategy := mediator.NewParallelPublishStrategy()
		handlers := []interface{}{
			newParallelHandler(nil),
			newParallelHandler(nil),
		}

		result := strategy.Execute(context.Background(),
			handlers,
			func(ctx context.Context, handler interface{}) error {
				return handler.(parallelHandler).Execute()
			})

		assert.NoError(t, result)
		for _, handler := range handlers {
			assert.True(t, handler.(parallelHandler).Executed())
		}
	})

	t.Run("error will not stop the execution", func(t *testing.T) {
		strategy := mediator.NewParallelPublishStrategy()

		handlers := []interface{}{
			newParallelHandler(nil),
			newParallelHandler(errors.New("error")),
			newParallelHandler(nil),
		}

		result := strategy.Execute(context.Background(),
			handlers,
			func(ctx context.Context, handler interface{}) error {
				return handler.(parallelHandler).Execute()

			})

		assert.Error(t, result)
		for _, handler := range handlers {
			assert.True(t, handler.(parallelHandler).Executed())
		}
	})
}
