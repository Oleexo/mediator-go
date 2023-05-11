package mediator

import "reflect"

type Container interface {
	resolve(request interface{}) (interface{}, bool)
}

type container struct {
	requestHandlers      map[reflect.Type]interface{}
	notificationHandlers map[reflect.Type]interface{}
}

func (c container) resolve(request interface{}) (interface{}, bool) {
	handler, ok := c.requestHandlers[reflect.TypeOf(request)]
	return handler, ok
}

func New(requestHandlers []RequestHandlerDefinition,
	notificationHandlers []NotificationHandlerDefinition) Container {
	reqHandlers := make(map[reflect.Type]interface{},
		len(requestHandlers))
	for _, requestHandler := range requestHandlers {
		reqHandlers[requestHandler.RequestType()] = requestHandler.Handler()
	}
	notifHandlers := make(map[reflect.Type]interface{},
		len(notificationHandlers))
	for _, notificationHandler := range notificationHandlers {
		notifHandlers[notificationHandler.NotificationType()] = notificationHandler.Handler()
	}
	return &container{
		requestHandlers:      reqHandlers,
		notificationHandlers: notifHandlers,
	}
}
