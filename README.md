# Mediator for Go ✨

Inspired by [MediatR](https://github.com/jbogard/MediatR)

Simple mediator implementation for Go with no dependencies.

## 📜 Table of Contents

1. [✨ Summary](#-table-of-contents)
2. [🚀 Installation](#-installation)
3. [📖 Usage](#-usage)
    - [📦 Requests](#-requests)
        - [Step 1: Define a request and its response](#step-1-define-a-request-and-its-response)
        - [Step 2: Define a request handler](#step-2-define-a-request-handler)
        - [⚙️ Using `mediator.Send[]()`](#️-using-mediatorsend)
        - [⚡ Using `sender.Send()`](#️-using-sendersend)
    - [🔗 Pipeline Behavior](#-pipeline-behavior)
4. [📢 Notifications](#-notifications)

---

## 🚀 Installation

Use Go modules to install mediator-go in your application:

```bash
go get github.com/Oleexo/mediator-go
```

---

## 📖 Usage

There are two main concepts in mediator-go: **request** and **notification**.

- **Request**: A message processed by one handler, returning a response.
- **Notification**: A message processed by multiple handlers, returning nothing.

### 📦 Requests

All requests should implement the `mediator.Request` interface, and all request handlers should implement the
`mediator.RequestHandler` interface. Requests can pass through `PipelineBehavior` if some are registered.

✅ **Best practice:** Always ensure your request-response types are clearly documented and tested for edge cases.

#### Step 1: Define a request and its response

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

// MyResponse is an example of a response.
type MyResponse struct {
	Result string
}
```

#### Step 2: Define a request handler

```go
package mypackage

import (
	"context"
)

// MyRequestHandler is an example of a request handler.
// All request handlers should implement mediator.RequestHandler interface.
type MyRequestHandler struct {
}

// Constructor function for MyRequestHandler
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

---

Now it's time to call your request through the mediator!

### ⚙️ Using `mediator.Send[]()`

This method uses minimal reflection to send the request to the handler. Use this method for performance-critical
scenarios.

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
	// 🌟 Create the request handler
	requestHandler := NewMyRequestHandler()

	// Associate the handler with the request and response
	requestDefinitions := []mediator.RequestHandlerDefinition{
		mediator.NewRequestHandlerDefinition[MyRequest, MyResponse](requestHandler),
	}

	// 🚀 Create the send container with all handler definitions
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
	// 🧠 Register the send container
	sendContainer := mediator.NewSendContainer(...)
	ctx := context.Background()

	request := MyRequest{Name: "John"}

	// ✨ Send and process the request
	response, err := mediator.Send[MyRequest, MyResponse](ctx, sendContainer, request)
	if err != nil {
		// ❌ Handle errors properly
		panic(err)
	}

	fmt.Printf("🎉 Response: %s", response.Result)
}
```

---

### ⚡ Using `sender.Send()`

This method uses reflection to send the request to the handler. It is more flexible and easier to inject.

```go
package main

import (
	"github.com/Oleexo/mediator-go"
)

func main() {
	// 🌟 Create the request handler
	requestHandler := NewMyRequestHandler()

	// Associate the handler with the request and response
	requestDefinitions := []mediator.RequestHandlerDefinition{
		mediator.NewRequestHandlerDefinition[MyRequest, MyResponse](requestHandler),
	}

	// 🚀 Create the send container
	sendContainer := mediator.NewSendContainer(
		mediator.WithRequestDefinitionHandlers(requestDefinitions...),
	)

	// 🌟 Create the sender
	sender := mediator.NewSender(sendContainer)
}
```

The second step is sending the request through the sender:

```go
package main

import (
	"context"
	"fmt"

	"github.com/Oleexo/mediator-go"
)

func main() {
	// 🧠 Register the sender
	sender := mediator.NewSender(sendContainer)
	ctx := context.Background()

	request := MyRequest{Name: "Jane"}

	// ✨ Send and process the request
	r, err := sender.Send(ctx, request)
	if err != nil {
		// ❌ Handle errors properly
		panic(err)
	}

	response := r.(MyResponse)

	fmt.Printf("🎉 Response: %s", response.Result)
}
```

---

### 🔗 Pipeline behavior

You can add pipeline behaviors to the request to introduce cross-cutting concerns such as **validation**, **logging**,
and **performance monitoring**.

```go
package main

import (
	"fmt"

	"github.com/Oleexo/mediator-go"
	"github.com/Oleexo/mediator-go/pipelines"
)

func main() {
	// 🚀 Register with pipeline behaviors
	container := mediator.NewSendContainer(
	    ..., // Registering request handler
	    pipelines.WithStructValidation(), // Example: Validation pipeline
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

### 📢 Notifications

Notifications work differently—they are processed by multiple handlers and do not return results. Use notifications for
**event-driven systems** or **pub-sub designs**.

✅ **Best practice**: Keep handler logic short and idempotent for notifications.

---

### 📚 Modules

- [🔗 Fx integration](https://github.com/Oleexo/mediator-go-fx): Easily integrate mediator-go
  with [Fx](https://uber-go.github.io/fx/index.html) for dependency injection.
- [✅ Validation pipeline](https://github.com/Oleexo/mediator-go-valid): Add robust validation to your requests
  using [validator](https://github.com/go-playground/validator).

## 💡 Contributing

🤝 Pull requests are welcome! For significant changes, please create an issue first to discuss your proposal.

✅ **Best practices for contributing**:

1. Ensure high test coverage.
2. Write clear and concise documentation.
3. Follow Go idioms and naming conventions.