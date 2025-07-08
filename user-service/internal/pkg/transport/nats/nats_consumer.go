package nats

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/go-kit/kit/endpoint"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func NewSubscriber(
	ctx context.Context,
	js jetstream.JetStream,
	subject string,
	ep endpoint.Endpoint,
	dec Decoder,
	enc Encoder,
	mw []endpoint.Middleware,
	opts ...nats.SubOpt,
) (*Consumer, error) {
	c := &Consumer{
		js:      js,
		subject: subject,
		ep:      ep,
		dec:     dec,
		enc:     enc,
	}

	if err := c.createConsumer(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

type Consumer struct {
	js           jetstream.JetStream
	subject      string
	ep           endpoint.Endpoint
	dec          Decoder
	enc          Encoder
	mw           []endpoint.Middleware
	opts         []nats.SubOpt
	consumerName string
	streamName   string
	consumer     jetstream.Consumer
	consumerCtx  jetstream.ConsumeContext
}

func (c *Consumer) createConsumer(ctx context.Context) error {

	cons, err := c.js.CreateOrUpdateConsumer(ctx, c.streamName, jetstream.ConsumerConfig{
		Name:          c.consumerName,
		FilterSubject: c.subject,
	})
	if err != nil {
		return fmt.Errorf("failed to create or update consumer: %w", err)
	}

	c.consumer = cons

	return nil
}

func (c *Consumer) Start(ctx context.Context) error {
	slog.Info("starting nats consumer", "subject", c.subject)
	var err error

	epWithMiddleware := Chain(c.mw...)(c.ep)

	c.consumerCtx, err = c.consumer.Consume(func(msg jetstream.Msg) {
		request, err := c.dec(ctx, &msg)
		if err != nil {
			slog.Error("failed to decode message", "error", err)
			msg.Nak()

			return
		}

		_, err = epWithMiddleware(ctx, request)
		if err != nil {
			slog.Error("failed to execute endpoint", "error", err)
			msg.Nak()

			return
		}

	})
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	return nil
}

func (c *Consumer) Stop() {
	slog.Info("stopping nats consumer", "subject", c.subject)
	c.consumerCtx.Drain()
}
