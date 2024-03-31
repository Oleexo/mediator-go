package mediator

import (
	"context"
	"sync"
)

type parallelPublishStrategy struct {
}

func (p parallelPublishStrategy) Execute(ctx context.Context, handlers []interface{}, launcher LaunchHandler) error {
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

func NewParallelPublishStrategy() PublishStrategy {
	return parallelPublishStrategy{}
}

func WithParallelPublishStrategy() func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.PublishStrategy = NewParallelPublishStrategy()
	}
}
