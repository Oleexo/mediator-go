# Mediator for Go

Inspire from [MediatR](https://github.com/jbogard/MediatR)

Simple mediator implementation for Go with no dependencies.

## Summary

- [Installation](#installation)
- [Usage](#usage)
  - [Requests](#requests)
    - [Using `mediator.Send[]()`](#using-mediatorsend)
    - [Using `sender.Send()`](#using-sendersend)
    - [Pipeline behavior](#pipeline-behavior)
  - [Notifications](#notifications)
    - [Using `mediator.Publish[]()`](#using-mediatorpublish)
    - [Using `publisher.Publish()`](#using-publisherpublish)
    - [Publish strategy](#publish-strategy)
- [Contributing](#contributing)

## Installation

Use Go modules to install mediator-go in your application

```bash
go get github.com/Oleexo/mediator-go
```

## Usage

There are two main concepts in mediator-go: request and notification.

- Request is a message process by one handler and return a response.
- Notification is a message process by multiple handlers and return nothing.

### Requests

All request should implement `mediator.Request` interface and all request handler should
implement `mediator.RequestHandler` interface.
The request can be pass through `PipelineBehavior` if some are registered.

The first step is to define a request and its response.

```go
package mypackage

import (
	"fmt"
)

// MyRequest is an example of request. 
// All requests should implement mediator.Request interface.
type MyRequest struct {
	Name string
}

// String is a method to return a string representation of the request.
func (r MyRequest) String() string {
	return fmt.Sprintf("MyRequest{Name=%s}", r.Name)
}

// MyResponse is an example of response.
type MyResponse struct {
	Result string
}

```

The second step is to define a request handler.

```go
package mypackage

import (
	"context"
)

// MyRequestHandler is an example of request handler.
// All request handlers should implement mediator.RequestHandler interface.
type MyRequestHandler struct {
}

func NewMyRequestHandler() *MyRequestHandler {
	return &MyRequestHandler{}
}

// Handle is a method to handle the request.
func (h MyRequestHandler) Handle(_ context.Context, cmd MyRequest) (MyResponse, error) {
	// todo: your request code

	return MyResponse{
		Result: "Hello " + cmd.Name,
	}, nil
}
```

Now it's time to call your request through the mediator.

There is two methods to send the request to the handler:

- `mediator.Send[]()` is the generic method to send the request to the handler with the minimum of reflection.
- `sender.Send()` is the method to send the request to the handler with the reflection.

#### Using `mediator.Send[]()`

This method use the minimum of reflection to send the request to the handler.
The method is more complicate to mock or inject. Use (sender.Send())[#### Using `sender.Send()`] for more flexibility.

The first step is to create the `SendContainer` with the request handler definition.

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
	// Create the request handler
	requestHandler := NewMyRequestHandler()

	// Create a definition of the request handler associated with the request and the response
	requestDefinitions := []mediator.RequestHandlerDefinition{
		mediator.NewRequestHandlerDefinition[MyRequest, MyResponse](requestHandler),
	}

	// Create the send container with all the request definitions
	sendContainer := mediator.NewSendContainer(
		mediator.WithRequestDefinitionHandlers(requestDefinitions...),
	)
}
```

The second step is to send the request to the handler.

```go
package main

import (
	"context"
	"fmt"

	"github.com/Oleexo/mediator-go"
)

func main() {
	// registering 
	sendContainer := mediator.NewSendContainer(...)
	ctx := context.Background()

	response, err := mediator.Send[MyRequest, MyResponse](ctx, container, request)
	if err != nil {
		// todo: handle error
		panic(err)
	}

	fmt.Printf("Response: %s", response.Result)
}
```

#### Using `sender.Send()`

This method use the reflection to send the request to the handler.

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
	// Create the request handler
	requestHandler := NewMyRequestHandler()

	// Create a definition of the request handler associated with the request and the response
	requestDefinitions := []mediator.RequestHandlerDefinition{
		mediator.NewRequestHandlerDefinition[MyRequest, MyResponse](requestHandler),
	}

	// Create the send container with all the request definitions
	sendContainer := mediator.NewSendContainer(
		mediator.WithRequestDefinitionHandlers(requestDefinitions...),
	)
	
	sender := mediator.NewSender(sendContainer)
}
```

The second step is to send the request to the handler.

```go
package main

import (
	"context"
	"fmt"

	"github.com/Oleexo/mediator-go"
)

func main() {
	// registering 
    sender := mediator.NewSender(sendContainer)
	ctx := context.Background()

	r, err := sender.Send(ctx, request)
	if err != nil {
		// todo: handle error
		panic(err)
	}
	
	response := r.(MyResponse)

	fmt.Printf("Response: %s", response.Result)
}
```

#### Pipeline behavior

You can add pipeline behavior to the request.

```go
package main

import (
	"fmt"

	"github.com/Oleexo/mediator-go"
	"github.com/Oleexo/mediator-go/pipelines"
)

func main() {
	// registering 
	container := mediator.NewSendContainer(
	    ..., // registering request handler
	    pipelines.WithStructValidation(),
    )

	response, err := mediator.SendWithoutContext[MyRequest, MyResponse](container, request)
	if err != nil {
		// todo: handle error
		panic(err)
	}

	fmt.Printf("Response: %s", response.Result)
}

```

### Notifications

A notification have no base interface to implement.
The notification will not be pass through the `PipelineBehavior`.

The first step is to define a notification.

```go
package mypackage

// MyNotification is an example of notification.
type MyNotification struct {
	Name string
}
```

The second step is to define all handlers for this notification.

```go
package mypackage

import (
    "context"
    "fmt"
)

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
```

Like the request, there are two methods to publish the notification to handlers.
- 
- `mediator.Publish[]()` is the generic method to send the notification to handlers with the minimum of reflection.
- `publisher.Publish()` is the method to send the notification to handlers with the reflection.

#### Using `mediator.Publish[]()`

This method use the minimum of reflection to send the notification to the handlers.

The first step is to create the `PublishContainer` with the notification handler definition.

```go
package main

import (
	"context"
	
	"github.com/Oleexo/mediator-go"
)

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
	publishContainer := mediator.NewPublishContainer(
		mediator.WithNotificationDefinitionHandlers(definitions...),
	)
}
```

The second step is to publish the notification with the `publishContainer`.

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
	// Create a new container with the notification definitions
	publishContainer := mediator.NewPublishContainer(...)

	notification := MyNotification{}

	err := mediator.Publish(ctx, container, notification)
	if err != nil {
		// todo: handle error
		panic(err)
	}
}
```

#### Using `publisher.Publish()`

This method use the reflection to send the notification to the handlers.

This first step is to create the `Publisher` with the notification handler definition.

```go
package main

import (
	"context"

	"github.com/Oleexo/mediator-go"
)

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
  publishContainer := mediator.NewPublishContainer(
    mediator.WithNotificationDefinitionHandlers(definitions...),
  )

  // Create a new container with the notification definitions
	publisher := mediator.NewPublisher(publishContainer)
}
```

The second step is to publish the notification with the `publisher`.

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
	// Create a new container with the notification definitions
	publisher := mediator.NewPublisher()...)

	notification := MyNotification{}

	err := publisher.Publish(ctx, notification)
	if err != nil {
		// todo: handle error
		panic(err)
	}
}
```
#### Publish strategy

Publish strategies are the way to handle notification through the handlers.
There are two strategies available:
- Synchronous (Default): The handlers will be executed one by one and the process stop at the first error
- Parallel: The handlers will be executed in parallel and the process will return the first error

The strategy can be set at the creation of the `PublishContainer` or `Publisher`.

```go
package main

import (
    "github.com/Oleexo/mediator-go"
)

func main() {
    // Create a new container with the notification definitions
    publishContainer := mediator.NewPublishContainer(
        mediator.WithNotificationDefinitionHandlers(definitions...),
        mediator.WithPublishStrategy(mediator.Parallel),
    )
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.