package endpoint

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/dto"
)

type UserService interface {
	OnCreatedUser(ctx context.Context, req dto.UserCreated) error
}

func NewUserEndpoint(svc UserService) User {
	return User{
		OnCreated: MakeOnCreatedUserEndpoint(svc),
	}
}

func MakeOnCreatedUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*dto.UserCreated)
		if !ok {
			return nil, fmt.Errorf("user service: %w", ErrInvalidType)
		}

		err := svc.OnCreatedUser(ctx, *req)
		if err != nil {
			return nil, fmt.Errorf("user service: %w", err)
		}

		return nil, nil
	}
}
