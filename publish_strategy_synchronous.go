package mediator

import "context"

type synchronousPublishStrategy struct {
}

func (s synchronousPublishStrategy) Execute(ctx context.Context, handlers []interface{}, launcher LaunchHandler) error {
	for _, handler := range handlers {
		err := launcher(ctx, handler)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewSynchronousPublishStrategy() PublishStrategy {
	return synchronousPublishStrategy{}
}

func WithSynchronousPublishStrategy() func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.PublishStrategy = NewSynchronousPublishStrategy()
	}
}
