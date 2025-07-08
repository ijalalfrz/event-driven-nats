package nats

import (
	"context"

	"log/slog"

	"github.com/go-kit/kit/endpoint"
	"github.com/nats-io/nats.go/jetstream"
)

func AutoAckMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			msg := ctx.Value("nats-msg").(jetstream.Msg)

			defer func() {
				if err != nil {
					msg.Nak()
				} else {
					if ackErr := msg.Ack(); ackErr != nil {
						slog.ErrorContext(ctx, "failed to ack message", "error", ackErr)
					}
				}
			}()

			return next(ctx, request)
		}
	}
}

// run middlewares in reverse order
func Chain(middlewares ...endpoint.Middleware) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
