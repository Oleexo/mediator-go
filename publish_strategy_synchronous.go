package mediator

import "context"

// synchronousPublishStrategy is a struct implementing the PublishStrategy to execute handlers sequentially.
type synchronousPublishStrategy struct {
}

func (s synchronousPublishStrategy) Execute(ctx context.Context, handlers []interface{}, launcher NotificationHandlerFunc) error {
	for _, handler := range handlers {
		err := launcher(ctx, handler)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewSynchronousPublishStrategy returns a PublishStrategy instance that executes handlers sequentially and stops on error.
func NewSynchronousPublishStrategy() PublishStrategy {
	return synchronousPublishStrategy{}
}

// WithSynchronousPublishStrategy sets the publish strategy to a synchronous execution model for the given PublishOptions.
func WithSynchronousPublishStrategy() func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.Strategy = NewSynchronousPublishStrategy()
	}
}
