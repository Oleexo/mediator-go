package main

import (
	"context"
	"fmt"
	"github.com/Oleexo/mediator-go"
)

type MyNotification struct {
	Name string
}

type MyNotificationHandler1 struct {
}

func (*MyNotificationHandler1) Handle(ctx context.Context, request MyNotification) error {
	fmt.Printf("Handler 1\n")
	return nil
}

type MyNotificationHandler2 struct {
}

func (*MyNotificationHandler2) Handle(ctx context.Context, request MyNotification) error {
	fmt.Printf("Handler 2\n")
	return nil
}

func NewMyNotificationHandler1() *MyNotificationHandler1 {
	return &MyNotificationHandler1{}
}

func NewMyNotificationHandler2() *MyNotificationHandler2 {
	return &MyNotificationHandler2{}
}

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
	container mediator.PublishContainer
}

func NewMyRequestHandler(container mediator.PublishContainer) *MyRequestHandler {
	return &MyRequestHandler{
		container: container,
	}
}

func (h *MyRequestHandler) Handle(ctx context.Context, cmd MyRequest) (MyResponse, error) {
	// todo: your request code
	notification := MyNotification{
		Name: "MyNotification",
	}

	// Publish a notification
	if err := mediator.Publish(ctx, h.container, notification); err != nil {
		return MyResponse{}, err
	}

	// Return a response
	return MyResponse{
		Result: "Hello " + cmd.Name,
	}, nil
}

func main() {
	externalContext := context.Background()
	handler1 := NewMyNotificationHandler1()
	handler2 := NewMyNotificationHandler2()
	def1 := mediator.NewNotificationHandlerDefinition[MyNotification](handler1)
	def2 := mediator.NewNotificationHandlerDefinition[MyNotification](handler2)

	notificationDefinitions := []mediator.NotificationHandlerDefinition{
		def1,
		def2,
	}
	notificationContainer := mediator.NewPublishContainer(
		mediator.WithNotificationDefinitionHandlers(notificationDefinitions...),
	)

	handler := NewMyRequestHandler(notificationContainer)
	def := mediator.NewRequestHandlerDefinition[MyRequest, MyResponse](handler)

	requestDefinitions := []mediator.RequestHandlerDefinition{
		def,
	}
	container := mediator.NewSendContainer(
		mediator.WithRequestDefinitionHandlers(requestDefinitions),
	)

	request := MyRequest{}

	response, err := mediator.Send[MyRequest, MyResponse](externalContext, container, request)
	if err != nil {
		// todo: handle error
		panic(err)
	}

	fmt.Printf("Response: %s", response.Result)
}
