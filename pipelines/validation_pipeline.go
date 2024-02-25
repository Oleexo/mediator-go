package pipeline

import (
	"context"
	"github.com/Oleexo/mediator-go"
	"github.com/go-playground/validator/v10"
)

type ValidationPipeline struct {
	validate *validator.Validate
}

func NewValidationPipeline(optFns ...func(options *Options)) *ValidationPipeline {
	validate := validator.New()
	options := Options{
		Validator: validate,
	}
	for _, optFn := range optFns {
		optFn(&options)
	}
	return &ValidationPipeline{
		validate: validate,
	}
}

type Options struct {
	Validator *validator.Validate
}

func WithStructValidation(fn validator.StructLevelFunc, types ...interface{}) func(options *Options) {
	return func(options *Options) {
		options.Validator.RegisterStructValidation(fn, types...)
	}
}

func WithCustomValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) func(options *Options) {
	return func(options *Options) {
		err := options.Validator.RegisterValidation(tag, fn, callValidationEvenIfNull...)
		if err != nil {
			panic(err)
		}
	}
}

func (l ValidationPipeline) Handle(_ context.Context,
	request mediator.BaseRequest,
	next mediator.RequestHandlerFunc) (interface{}, error) {

	if err := l.validate.Struct(request); err != nil {
		return nil, err
	}

	return next()
}
