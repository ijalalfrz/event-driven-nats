package nats

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/go-kit/kit/endpoint"
	"github.com/nats-io/nats.go/jetstream"
)

func NewSubscriber[T any](
	ctx context.Context,
	js jetstream.Stream,
	subject string,
	ep endpoint.Endpoint,
	dec Decoder[T],
	mw []endpoint.Middleware,
) (*Consumer[T], error) {
	c := &Consumer[T]{
		js:      js,
		subject: subject,
		ep:      ep,
		dec:     dec,
	}

	if err := c.createConsumer(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

type Consumer[T any] struct {
	js           jetstream.Stream
	subject      string
	ep           endpoint.Endpoint
	dec          Decoder[T]
	mw           []endpoint.Middleware
	consumerName string
	consumer     jetstream.Consumer
	consumerCtx  jetstream.ConsumeContext
}

func (c *Consumer[T]) createConsumer(ctx context.Context) error {

	cons, err := c.js.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          c.consumerName,
		FilterSubject: c.subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    1,
	})
	if err != nil {
		return fmt.Errorf("failed to create or update consumer: %w", err)
	}

	c.consumer = cons

	return nil
}

func (c *Consumer[T]) Start(ctx context.Context) error {
	slog.Info("starting nats consumer", "subject", c.subject)
	var err error

	c.consumerCtx, err = c.consumer.Consume(func(msg jetstream.Msg) {
		request, err := c.dec(ctx, msg)
		if err != nil {
			slog.Error("failed to decode message", "error", err)
			msg.Nak()

			return
		}

		_, err = c.ep(ctx, request)
		if err != nil {
			slog.Error("failed to execute endpoint", "error", err)
			msg.Nak()

			return
		}

		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	return nil
}

func (c *Consumer[T]) Stop() {
	slog.Info("stopping nats consumer", "subject", c.subject)
	c.consumerCtx.Drain()
}
