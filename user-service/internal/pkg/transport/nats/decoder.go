package nats

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type Decoder func(ctx context.Context, msg *jetstream.Msg) (interface{}, error)
