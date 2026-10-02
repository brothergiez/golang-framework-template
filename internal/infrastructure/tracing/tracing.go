package tracing

import "context"

// Ready stubs for OpenTelemetry integration without making it mandatory.
type Tracer interface {
	Start(ctx context.Context, name string) (context.Context, func())
}

type NoopTracer struct{}

func NewNoop() *NoopTracer { return &NoopTracer{} }

func (n *NoopTracer) Start(ctx context.Context, _ string) (context.Context, func()) {
	return ctx, func() {}
}
