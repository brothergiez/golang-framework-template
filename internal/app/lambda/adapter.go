package lambda

import (
	"context"
	"encoding/json"
	"fmt"
)

type Handler interface {
	Handle(ctx context.Context, event any) (any, error)
}

type Adapter struct {
	handler Handler
}

func NewAdapter(handler Handler) *Adapter {
	return &Adapter{handler: handler}
}

func (a *Adapter) Invoke(ctx context.Context, raw []byte) ([]byte, error) {
	var event any
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &event); err != nil {
			return nil, fmt.Errorf("invalid event json: %w", err)
		}
	}
	result, err := a.handler.Handle(ctx, event)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result)
}
