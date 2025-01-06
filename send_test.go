package mediator_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/Oleexo/mediator-go"
	"github.com/stretchr/testify/assert"
	"strings"
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
	Order int
}

func (t TestPipelineBehavior) Handle(ctx context.Context, request mediator.BaseRequest, next mediator.RequestHandlerFunc) (interface{}, error) {
	r, ok := request.(*TestRequest)
	if !ok {
		return nil, errors.New("request is not a TestRequest")
	}
	if strings.Contains(r.Value, "pipeline") {
		r.Value += fmt.Sprintf(".%d", t.Order)
	} else {
		r.Value = fmt.Sprintf("pipeline.%d", t.Order)
	}
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
			mediator.WithPipelineBehavior(TestPipelineBehavior{
				Order: 1,
			}),
		)

		request := TestRequest{Value: "test"}

		response, err := mediator.Send[*TestRequest, string](context.Background(), container, &request)
		assert.NoError(t, err)

		assert.Equal(t, "pipeline.1", response)
	})

	t.Run("should apply pipeline in reverse order", func(t *testing.T) {
		handler := TestRequestHandler{}
		container := mediator.NewSendContainer(
			mediator.WithRequestDefinitionHandler(mediator.NewRequestHandlerDefinition[*TestRequest, string](handler)),
			mediator.WithPipelineBehavior(TestPipelineBehavior{
				Order: 1,
			}),
			mediator.WithPipelineBehavior(TestPipelineBehavior{
				Order: 2,
			}),
		)

		request := TestRequest{Value: "test"}

		response, err := mediator.Send[*TestRequest, string](context.Background(), container, &request)
		assert.NoError(t, err)

		assert.Equal(t, "pipeline.2.1", response)
	})
}
