package mediator_test

import (
	"context"
	"github.com/Oleexo/mediator-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestRequest struct {
	Value string
}

func (t TestRequest) String() string {
	return t.Value
}

type TestRequestHandler struct {
}

func (t TestRequestHandler) Handle(ctx context.Context, request *TestRequest) (string, error) {
	return request.Value, nil
}

type TestPipelineBehavior struct {
}

func (t TestPipelineBehavior) Handle(ctx context.Context, request mediator.BaseRequest, next mediator.RequestHandlerFunc) (interface{}, error) {
	request.(*TestRequest).Value = "pipeline"
	return next()
}

func TestSend(t *testing.T) {
	t.Run("should send a request to a single handler", func(t *testing.T) {
		handler := TestRequestHandler{}
		container := mediator.NewSendContainer(
			mediator.WithRequestDefinitionHandler(mediator.NewRequestHandlerDefinition[*TestRequest, string](handler)),
		)

		request := TestRequest{Value: "test"}

		response, err := mediator.Send[*TestRequest, string](context.Background(), container, &request)
		assert.NoError(t, err)

		assert.Equal(t, request.Value, response)
	})

	t.Run("should apply pipeline behavior", func(t *testing.T) {
		handler := TestRequestHandler{}
		container := mediator.NewSendContainer(
			mediator.WithRequestDefinitionHandler(mediator.NewRequestHandlerDefinition[*TestRequest, string](handler)),
			mediator.WithPipelineBehavior(TestPipelineBehavior{}),
		)

		request := TestRequest{Value: "test"}

		response, err := mediator.Send[*TestRequest, string](context.Background(), container, &request)
		assert.NoError(t, err)

		assert.Equal(t, "pipeline", response)
	})
}
