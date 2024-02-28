package main

import (
	"context"
	"fmt"

	"github.com/Oleexo/mediator-go"
)

type MyRequest struct {
	Name string
}

func (r MyRequest) String() string {
	return fmt.Sprintf("MyRequest{Name=%s}", r.Name)
}

type MyResponse struct {
	Result string
}

type MyRequestHandler struct {
}

func NewMyRequestHandler() *MyRequestHandler {
	return &MyRequestHandler{}
}

func (h *MyRequestHandler) Handle(ctx context.Context, cmd MyRequest) (MyResponse, error) {
	// todo: your request code
	return MyResponse{
		Result: "Hello " + cmd.Name,
	}, nil
}

func main() {
	handler := NewMyRequestHandler()
	def := mediator.NewRequestHandlerDefinition[MyRequest, MyResponse](handler)

	requestDefinitions := []mediator.RequestHandlerDefinition{
		def,
	}
	notificationDefinitions := make([]mediator.NotificationHandlerDefinition, 0)
	container := mediator.New(
		mediator.WithRequestDefinitionHandlers(requestDefinitions),
		mediator.WithNotificationDefinitionHandlers(notificationDefinitions),
	)

	request := MyRequest{}

	response, err := mediator.SendWithoutContext[MyRequest, MyResponse](container, request)
	if err != nil {
		// todo: handle error
		panic(err)
	}

	fmt.Printf("Response: %s", response.Result)
}
