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

func notificationWithContainer(ctx context.Context, container mediator.PublishContainer) {

	notification := MyNotification{}

	err := mediator.Publish(ctx, container, notification)
	if err != nil {
		// todo: handle error
		panic(err)
	}
}

func notificationWithPublisher(ctx context.Context, container mediator.PublishContainer) {
	// Create a publisher with the notification definitions
	publisher := mediator.NewPublisher(container)

	notification := MyNotification{}

	err := publisher.Publish(ctx, notification)
	if err != nil {
		// todo: handle error
		panic(err)
	}

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
	// Create a new container with the notification definitions
	container := mediator.NewPublishContainer(
		mediator.WithNotificationDefinitionHandlers(notificationDefinitions...),
	)

	notificationWithContainer(externalContext, container)
	notificationWithPublisher(externalContext, container)
}
