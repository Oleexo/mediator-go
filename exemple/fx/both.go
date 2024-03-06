package main

import (
	"context"
	"fmt"
	"github.com/Oleexo/mediator-go"
	mediatorfx "github.com/Oleexo/mediator-go/fx"
	"go.uber.org/fx"
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
	var constructor []any

	constructor = append(constructor, mediatorfx.AsRequestHandler[MyRequest, MyResponse](NewMyRequestHandler)...)
	constructor = append(constructor, mediatorfx.AsNotificationHandler[MyNotification, *MyNotificationHandler1](NewMyNotificationHandler1)...)
	constructor = append(constructor, mediatorfx.AsNotificationHandler[MyNotification, *MyNotificationHandler2](NewMyNotificationHandler2)...)

	fx.New(
		fx.Provide(constructor...),
		mediatorfx.NewModule(),
		fx.Invoke(func(container mediator.SendContainer) {
			request := MyRequest{}
			response, err := mediator.SendWithoutContext[MyRequest, MyResponse](container, request)
			if err != nil {
				// todo: handle error
				panic(err)
			}

			fmt.Printf("Response: %s", response.Result)
		}),
	).Run()
}
