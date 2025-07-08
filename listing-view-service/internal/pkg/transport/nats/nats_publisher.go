package nats

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type Publisher struct {
	js  jetstream.JetStream
	enc Encoder
}

func NewPublisher(js jetstream.JetStream, enc Encoder) *Publisher {
	return &Publisher{
		js:  js,
		enc: enc,
	}
}

// Publish encodes and sends a message to JetStream.
func (p *Publisher) Publish(ctx context.Context, subject string, request interface{}) (*jetstream.PubAck, error) {
	data, err := p.enc(ctx, request)
	if err != nil {
		return nil, err
	}

	return p.js.Publish(ctx, subject, data)
}
