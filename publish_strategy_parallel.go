package mediator

import (
	"context"
	"sync"
)

// parallelPublishStrategy is a struct implementing the PublishStrategy interface for concurrent handler execution.
type parallelPublishStrategy struct {
}

func (p parallelPublishStrategy) Execute(ctx context.Context,
	notifications []Notification,
	resolver Resolver,
	launcher NotificationHandlerFunc) error {

	ec := 0
	groups := make(map[Notification][]any)
	for _, notification := range notifications {
		handlers := resolver(notification)
		groups[notification] = handlers
		ec += len(handlers)
	}

	errChan := make(chan error, ec)
	var wg sync.WaitGroup

	c := 0
	for notification, handlers := range groups {
		for _, handler := range handlers {
			wg.Add(1)
			go func(ctx context.Context, notification Notification, handler interface{}) {
				defer wg.Done()
				errChan <- launcher(ctx, notification, handler)
			}(ctx, notification, handler)
			c++
		}
	}

	wg.Wait()

	for i := 0; i < c; i++ {
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

// WithParallelPublishStrategy sets the PublishStrategy to handle notifications with parallel execution of handlers.
func WithParallelPublishStrategy() func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.Strategy = NewParallelPublishStrategy()
	}
}
