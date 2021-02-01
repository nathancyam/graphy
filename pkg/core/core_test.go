package core

import (
	"context"
	"testing"
)

type CommandExample struct {
	Example string
}

type PolicyExample struct {
}

func (p PolicyExample) Run(ctx context.Context, _ interface{}) (context.Context, error) {
	return context.WithValue(ctx, PolicyExample{}, "really"), nil
}

type PolicyExample1 struct {
}

func (p PolicyExample1) Run(ctx context.Context, _ interface{}) (context.Context, error) {
	return ctx, nil
}

type CommandHandler struct {
}

func (h CommandHandler) Run(ctx context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

func TestModuleInit(t *testing.T) {
	mod := &Module{
		Policies: map[interface{}][]PolicyHandler{
			CommandExample{}: {
				PolicyExample{},
				PolicyExample1{},
			}},
		Handlers: map[interface{}]Handler{
			CommandExample{}: CommandHandler{},
		},
	}

	_, _ = mod.Process(context.Background(), &CommandExample{Example: ""})
}
