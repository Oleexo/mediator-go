package mediator

import (
	"context"
	"sync"
)

// parallelPublishStrategy is a struct implementing the PublishStrategy interface for concurrent handler execution.
type parallelPublishStrategy struct {
}

func (p parallelPublishStrategy) Execute(ctx context.Context, handlers []interface{}, launcher NotificationHandlerFunc) error {
	errChan := make(chan error, len(handlers))
	var wg sync.WaitGroup

	for _, handler := range handlers {
		go func(handler interface{}) {
			wg.Add(1)
			defer wg.Done()
			errChan <- launcher(ctx, handler)
		}(handler)
	}

	wg.Wait()

	for i := 0; i < len(handlers); i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}
	return nil
}

// NewParallelPublishStrategy creates a new instance of a parallel publish strategy for concurrent execution of handlers.
func NewParallelPublishStrategy() PublishStrategy {
	return parallelPublishStrategy{}
}

// WithParallelPublishStrategy sets the publish strategy to handle notifications with parallel execution of handlers.
func WithParallelPublishStrategy() func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.Strategy = NewParallelPublishStrategy()
	}
}
