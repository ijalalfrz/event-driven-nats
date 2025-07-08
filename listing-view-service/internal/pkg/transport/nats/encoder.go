package nats

import (
	"context"
	"encoding/json"
)

type Encoder func(context.Context, interface{}) ([]byte, error)

func JSONEncoder(ctx context.Context, request interface{}) ([]byte, error) {
	return json.Marshal(request)
}
