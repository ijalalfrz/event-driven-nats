package endpoint

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error)
	GetAllUsers(ctx context.Context, req dto.GetAllUsersRequest) (dto.GetAllUsersResponse, error)
	GetUserByID(ctx context.Context, req dto.GetUserByIDRequest) (dto.GetUserByIDResponse, error)
}

func NewUserEndpoint(userService UserService) User {
	return User{
		CreateUser:  makeCreateUserEndpoint(userService),
		GetAllUsers: makeGetAllUsersEndpoint(userService),
		GetUserByID: makeGetUserByIDEndpoint(userService),
	}
}

func makeCreateUserEndpoint(userService UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*dto.CreateUserRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request type: %w", ErrInvalidType)
		}

		res, err := userService.CreateUser(ctx, *req)
		if err != nil {
			return nil, fmt.Errorf("user service: %w", err)
		}

		return res, nil
	}
}

func makeGetAllUsersEndpoint(userService UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*dto.GetAllUsersRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request type: %w", ErrInvalidType)
		}

		res, err := userService.GetAllUsers(ctx, *req)
		if err != nil {
			return nil, fmt.Errorf("user service: %w", err)
		}

		return res, nil
	}
}

func makeGetUserByIDEndpoint(userService UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*dto.GetUserByIDRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request type: %w", ErrInvalidType)
		}

		res, err := userService.GetUserByID(ctx, *req)
		if err != nil {
			return nil, fmt.Errorf("user service: %w", err)
		}

		return res, nil
	}
}
