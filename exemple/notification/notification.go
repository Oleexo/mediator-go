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
	container := mediator.NewNotificationContainer(
		mediator.WithNotificationDefinitionHandlers(notificationDefinitions),
	)

	notification := MyNotification{}

	err := mediator.Publish(externalContext, container, notification)
	if err != nil {
		// todo: handle error
		panic(err)
	}
}
