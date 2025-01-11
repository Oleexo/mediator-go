package mediator

import "context"

// synchronousPublishStrategy is a struct implementing the PublishStrategy to execute handlers sequentially.
type synchronousPublishStrategy struct {
}

func (s synchronousPublishStrategy) Execute(ctx context.Context,
	notifications []Notification,
	resolver Resolver,
	launcher NotificationHandlerFunc) error {

	for _, notification := range notifications {
		handlers := resolver(notification)
		if handlers == nil {
			continue
		}
		for _, handler := range handlers {
			err := launcher(ctx, notification, handler)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// NewSynchronousPublishStrategy returns a PublishStrategy instance that executes handlers sequentially and stops on error.
func NewSynchronousPublishStrategy() PublishStrategy {
	return synchronousPublishStrategy{}
}

// WithSynchronousPublishStrategy sets the PublishStrategy to a synchronous execution model for the given PublishOptions.
func WithSynchronousPublishStrategy() func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.Strategy = NewSynchronousPublishStrategy()
	}
}
