package router

import (
	"context"

	natstransport "github.com/ijalalfrz/event-driven-nats/user-service/internal/pkg/transport/nats"
	"github.com/nats-io/nats.go/jetstream"
)

func MakeNATSHandler(
	ctx context.Context,
	js jetstream.JetStream,
) (*natstransport.Consumer, error) {
	return nil, nil
}
