package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/dto"
)

type PublicUserService interface {
	CreateUser(ctx context.Context, request dto.CreateUserRequest) (dto.CreateUserResponse, error)
}

func NewPublicUserEndpoint(
	service PublicUserService,
) PublicUser {
	return PublicUser{
		Create: makeCreateUserEndpoint(service),
	}
}

func makeCreateUserEndpoint(service PublicUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*dto.CreateUserRequest)
		if !ok {
			return nil, ErrInvalidType
		}

		response, err := service.CreateUser(ctx, *req)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}
