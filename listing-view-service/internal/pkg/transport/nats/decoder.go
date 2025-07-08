package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type Decoder[T any] func(ctx context.Context, msg jetstream.Msg) (*T, error)

func NewDecoder[T any]() Decoder[T] {
	return func(ctx context.Context, msg jetstream.Msg) (*T, error) {
		var data T
		if err := json.Unmarshal(msg.Data(), &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}
		return &data, nil
	}
}
