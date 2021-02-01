package core

import (
	"context"
)

type PolicyHandler interface {
	Run(ctx context.Context, cmd interface{}) (context.Context, error)
}

type Handler interface {
	Run(ctx context.Context, cmd interface{}) (interface{}, error)
}

type Module struct {
	Policies map[interface{}][]PolicyHandler
	Handlers map[interface{}]Handler
}

func (m Module) Process(ctx context.Context, cmd interface{}) (interface{}, error) {
	ps, ok := m.Policies[cmd]
	var updateCtx = ctx
	if ok {
		for _, p := range ps {
			updateCtx, _ = p.Run(updateCtx, cmd)
		}
	}
	h, ok := m.Handlers[cmd]
	if ok {
		h.Run(updateCtx, cmd)
	}
	return nil, nil
}
