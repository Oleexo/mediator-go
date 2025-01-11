# Mediator for Go ✨

Inspired by [MediatR](https://github.com/jbogard/MediatR)

Simple mediator implementation for Go with no dependencies.

**Why ?**

> *I mostly write web applications and strive to keep my code clean and simple.
By using CQRS through a mediator, I can effectively decouple the HTTP layer from the core application logic.*

# Table of Contents

- [Mediator for Go ✨](#mediator-for-go-)
- [Table of Contents](#table-of-contents)
- [Installation](#installation)
- [Usage](#usage)
	- [Send](#send)
		- [Define a request](#define-a-request)
		- [Define a handler](#define-a-handler)
		- [Send container](#send-container)
		- [Send a request](#send-a-request)
			- [Using request pipeline behavior](#using-request-pipeline-behavior)
	- [Publish](#publish)
		- [Define a notification](#define-a-notification)
		- [Define a handler](#define-a-handler-1)
		- [Publish container](#publish-container)
		- [Publish a notification](#publish-a-notification)
			- [Using a notification pipeline behavior](#using-a-notification-pipeline-behavior)
		- [Publish strategy](#publish-strategy)
			- [Using a strategy pipeline behavior](#using-a-strategy-pipeline-behavior)
- [📚 Modules](#-modules)
- [💡 Contributing](#-contributing)

---

# Installation

Use Go modules to install mediator-go in your application:

```bash
go get github.com/Oleexo/mediator-go
```

---

# Usage

## Send

Send is the mechanism to process a Request using a designated Handler,
which executes the action and returns a Response.
The Request represents the action to be performed,
while the Handler contains the logic to perform that action.

### Define a request

A request must implement the `mediator.Request` interface.
The request should include a `String() string` method to allow rapid serialization for logging purposes.

```go
package mypackage

import (
	"fmt"
)

// MyRequest is an example of a request. 
// All requests should implement mediator.Request interface.
type MyRequest struct {
	Name string
}

// String is a method to return a string representation of the request.
func (r MyRequest) String() string {
	return fmt.Sprintf("MyRequest{Name=%s}", r.Name)
}
```

✅ **Best practice:** Always ensure your request-response types are clearly documented and tested for edge cases.

### Define a handler

A request handler processes a previously defined request.
It should implement the `mediator.RequestHandler` interface
with the method `Handle(ctx context.Context, request TRequest) (TResponse, error)`
where `TRequest` is the type of the request and `TResponse` is the type of the response.
A request without a proper response type can use `mediator.Unit` to return "nothing".
*`Handle(ctx context.Context, request TRequest) (mediator.Unit, error)`*

The `Handle` method is invoked when `mediator.Send()` is called with the corresponding request type.
Ensure that the method implementation is thread-safe and efficiently handles the incoming request to produce the
expected response.

```go
package mypackage

import (
	"context"
)

// MyResponse is an example of a response.
type MyResponse struct {
	Result string
}

// MyRequestHandler is an example of a request handler.
// All request handlers should implement mediator.RequestHandler interface.
type MyRequestHandler struct {
}

// NewMyRequestHandler is a constructor function for MyRequestHandler
func NewMyRequestHandler() *MyRequestHandler {
	return &MyRequestHandler{}
}

// Handle is the method responsible for handling the request.
func (h MyRequestHandler) Handle(_ context.Context, cmd MyRequest) (MyResponse, error) {
	// 🚧 TODO: Implement your request processing logic here.

	return MyResponse{
		Result: "Hello " + cmd.Name,
	}, nil
}
```

### Send container

The send container is a struct that maintains the mapping between each `Request` type and its corresponding
`RequestHandler`. This is essential for routing the correct handler to process a given request.

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
	// Create the request handler
	requestHandler := NewMyRequestHandler()

	// Associate the handler with the request and response
	requestDefinitions := []mediator.RequestHandlerDefinition{
		mediator.NewRequestHandlerDefinition[MyRequest, MyResponse](requestHandler),
	}

	// Create the send container with all handler definitions
	sendContainer := mediator.NewSendContainer(
		mediator.WithRequestDefinitionHandlers(requestDefinitions...),
	)
}
```

### Send a request

Now it's time to send your request through the mediator!

There are currently two ways to send a request through the send container.

Using `mediator.Send()`:

```go
package main

import (
	"context"
	"fmt"
	"github.com/Oleexo/mediator-go"
)

func main() {
	// Send container creation ...
	var sendContainer mediator.SendContainer

	ctx := context.Background()
	request := MyRequest{
		Name: "Mediator !",
	}

	response, err := mediator.Send[MyRequest, MyResponse](ctx, sendContainer, request)
	if err != nil {
		fmt.Printf("The handler return an error: %s\n", err)
	} else {
		fmt.Printf("Response: %s", response)
    }
}
```

*This method also exists without context: `mediator.SendWithoutContext()`.
In this case `context.Background()` will be used*

Using `Sender`:

```go
package main

import (
	"context"
	"fmt"
	"github.com/Oleexo/mediator-go"
)

func main() {
	// Send container creation ...
	var sendContainer mediator.SendContainer
	
	sender := mediator.NewSender(sendContainer)

	ctx := context.Background()
	request := MyRequest{
		Name: "Mediator !",
	}

	response, err := sender.Send(ctx, request)
	if err != nil {
		fmt.Printf("The handler return an error: %s\n", err)
	} else {
		fmt.Printf("Response: %s", response)
    }
}
```

The `Sender` is more suitable for use in a testable context with mocking,
as it encapsulates the `SendContainer` logic, providing a cleaner and more
isolated interface for testing purposes. However, due to the lack of generics
in its method signature, the response is typed as `interface{}`, requiring
type assertions to handle specific response types.

#### Using request pipeline behavior

You can add pipeline behaviors to the request to introduce cross-cutting concerns such as **validation**,
**logging**, and **performance monitoring**.

Here an example with `mediator.LogRequestPipelineBehavior` to log information or error to each call.

```go
package main

import (
	"fmt"

	"github.com/Oleexo/mediator-go"
)

func main() {
    // Register with pipeline behaviors
    container := mediator.NewSendContainer(
        ..., // Registering request handler
		mediator.WithRequestPipelineBehavior(mediator.NewSlogRequestPipelineBehavior()), // Registering the logging behavior with slog
    )
    
    response, err := mediator.SendWithoutContext[MyRequest, MyResponse](container, request)
    if err != nil {
        // ❌ Handle errors properly
        panic(err)
    }
    
    fmt.Printf("🎉 Response: %s", response.Result)
}
```

---

## Publish

Publish is a mechanism for processing a Notification using multiple Handlers.
Each Handler executes specific actions in response to the Notification and does not return any value.
The Notification represents an event to be processed, while each Handler contains the logic to react to that event.

### Define a notification

A notification don't have a required implementation and can be anything.

```go
package mypackage

type MyNotification struct {
	Name string // A value
}
```

### Define a handler

A notification handler processes a specific notification type.
It should implement `mediator.NotificationHandler` interface
with the method `Handle(ctx context.Context, notification TNotification) error`

The `Handle` method is invoked when `mediator.Publish`  is called with the corresponding notification type.
Ensure that the method implementation is thread-safe and efficiently handles the incoming notification.

```go
package mypackage

import (
	"context"
	"fmt"
)

type MyNotificationHandler struct {
}

func NewMyNotificationHandler() MyNotificationHandler {
	return MyNotificationHandler{}
}
func (MyNotificationHandler) Handle(ctx context.Context, request MyNotification) error {
	fmt.Printf("Handler 1\n")
	return nil
}
```

### Publish container

Like the `SendContainer`, there is also a `PublishContainer`.
This struct maintans the mapping between each `Notification` type and its corresponding `RequestHandler`.
This is essential for routing the correct notification to handlers.

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
	// Create the notification handler
	handler := NewMyNotificationHandler()

	// Associate the handler with the notification
	requestDefinitions := []mediator.NotificationHandlerDefinition{
		mediator.NewNotificationHandlerDefinition[MyNotification](handler),
	}

	// Create the publish container with all handler definitions
	publisherContainer := mediator.NewPublishContainer(
		mediator.WithRequestDefinitionHandlers(requestDefinitions...),
	)
}
```

### Publish a notification

Now it's time to **publish** your request through the mediator!

There are currently two ways to send a request through the send container.

Using `mediator.Publish()`:

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

Using `Publisher`:

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
	// Create a new container with the notification definitions
	publishContainer := mediator.NewPublishContainer(...)
	publisher := mediator.NewPublisher(publishContainer)

	notification := MyNotification{}

	err := publisher.Publish(ctx, notification)
	if err != nil {
		// todo: handle error
		panic(err)
	}
}
```

The `Publisher` is more suitable for use in a testable context with mocking,
as it encapsulates the `PublisherContainer` logic, providing a cleaner and more
isolated interface for testing purposes.

#### Using a notification pipeline behavior

Like request, you can add a pipeline behavior to notification.
A notification behavior is executed for each handler.

Here an example with `mediator.LogNotificationPipelineBehavior` to log information or error to each call.

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
    // Register with pipeline behaviors
    container := mediator.NewPublishContainer
        ..., // Registering notification handler
		mediator.WithNotificationPipelineBehavior(mediator.NewSlogNotificationPipelineBehavior()), // Registering the logging behavior with slog
    )
    
    err := mediator.PublishWithoutContext(container, notification)
    if err != nil {
        // ❌ Handle errors properly
        panic(err)
    }
}
```

### Publish strategy

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
        mediator.WithParallelPublishStrategy(),
    )
}
```

#### Using a strategy pipeline behavior

Like the request or the notification, you can add a pipeline behavior to strategy.
The strategy pipeline behavior is executed before the strategy and the handlers.

Here an example with `mediator.LogStrategyPipelineBehavior`to log information or error before strategy.

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
    // Register with pipeline behaviors
    container := mediator.NewPublishContainer
        ..., // Registering notification handler
		mediator.WithStrategyPipelineBehavior(mediator.NewSlogNotificationPipelineBehavior()), // Registering the logging behavior with slog
    )
    
    err := mediator.PublishWithoutContext(container, notification)
    if err != nil {
        // ❌ Handle errors properly
        panic(err)
    }
}
```

---

# 📚 Modules

- [🔗 Fx integration](https://github.com/Oleexo/mediator-go-fx): Easily integrate mediator-go
  with [Fx](https://uber-go.github.io/fx/index.html) for dependency injection.
- [✅ Validation pipeline](https://github.com/Oleexo/mediator-go-valid): Add robust validation to your requests
  using [validator](https://github.com/go-playground/validator).

---

# 💡 Contributing

🤝 Pull requests are welcome! For significant changes, please create an issue first to discuss your proposal.

✅ **Best practices for contributing**:

1. Ensure high test coverage.
2. Write clear and concise documentation.
3. Follow Go idioms and naming conventions.