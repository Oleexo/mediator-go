package mediator

import "reflect"

type RequestHandlerDefinition interface {
	RequestType() reflect.Type
	Handler() interface{}
}

func NewRequestHandlerDefinition[TRequest Request[TResponse], TResponse interface{}](handler RequestHandler[TRequest, TResponse]) RequestHandlerDefinition {
	var request TRequest
	requestType := reflect.TypeOf(request)

	return &TypedRequestHandlerDefinition[TRequest, TResponse]{
		requestType: requestType,
		handler:     handler,
	}
}

type TypedRequestHandlerDefinition[TRequest Request[TResponse], TResponse interface{}] struct {
	requestType reflect.Type
	handler     RequestHandler[TRequest, TResponse]
}

func (t TypedRequestHandlerDefinition[TRequest, TResponse]) Handler() interface{} {
	return t.handler
}

func (t TypedRequestHandlerDefinition[TRequest, TResponse]) RequestType() reflect.Type {
	return t.requestType
}
