# Mediator for Go

Inspire from [MediatR](https://github.com/jbogard/MediatR)

Simple mediator implementation for Go with no dependencies.

## Installation

Use Go modules to install mediator-go in your application

```bash
go get github.com/Oleexo/mediator-go
```

## Usage

### Registering in container

The first step is to register request handler and notification handler in container to be able to send and publish.

For example a request with the handler

```go
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

func (h MyRequestHandler) Handle(_ context.Context, cmd MyRequest) (MyResponse, error) {
    // todo: your request code
	
    return MyResponse{
        Result: "Hello " + cmd.Name,
    }, nil
}

```

And for example a notification with the handler

```go
package main

import (
	"context"
	"fmt"

	"github.com/Oleexo/mediator-go"
)

type MyNotification struct {
    Name string
}

type MyNotificationHandler struct {
}

func NewMyNotificationHandler() *MyNotificationHandler {
    return &MyNotificationHandler{}
}

func (h MyNotificationHandler) Handle(_ context.Context, request MyNotification) error {
    fmt.Printf("Handler 1: %s\n", request.Name)
    return nil
}
```

Create container:

```go
package main

import (
    "github.com/Oleexo/mediator-go"
)

func main() {
    requestHandler := NewMyRequestHandler()
    notificationHandler = NewMyNotificationHandler()

    requestDefinitions := []mediator.RequestHandlerDefinition{
        mediator.NewRequestHandlerDefinition[MyRequest, MyResponse](requestHandler),
    }

    notificationDefinitions := []mediator.NotificationHandlerDefinition{
        mediator.NewNotificationHandlerDefinition[MyNotification](handnotificationHandlerler1),
    }
    sendContainer := mediator.NewSendContainer(
		mediator.WithRequestDefinitionHandlers(requestDefinitions...),
	)
	publishContainer := mediator.NewPublishContainer(
		mediator.WithNotificationDefinitionHandlers(notificationDefinitions...),
	)
}
```

### Send

With container, you can now send a request.

```go
package main

import (
    "fmt"

    "github.com/Oleexo/mediator-go"
)

func main() {
    // registering 
    container := mediator.NewSendContainer(...)

    response, err := mediator.SendWithoutContext[MyRequest, MyResponse](container, request)
    if err != nil {
        // todo: handle error
        panic(err)
    }

    fmt.Printf("Response: %s", response.Result)
}
```

### Publish

With container, you can now publish notification.

```go
package main

import (
    "fmt"

    "github.com/Oleexo/mediator-go"
)

func main() {
    // registering 
    container := mediator.NewPublishContainer(...)

    notification := MyNotification{}

    err := mediator.PublishWithoutContext(container, notification)
    if err != nil {
        // todo: handle error
        panic(err)
    }
}
```
