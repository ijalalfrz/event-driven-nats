package service

import (
	"context"
	"fmt"

	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/dto"
)

type PublicUserService struct {
	userServiceClient *UserServiceClient
}

func NewPublicUserService(
	userServiceClient *UserServiceClient,
) *PublicUserService {
	return &PublicUserService{userServiceClient: userServiceClient}
}

func (s *PublicUserService) CreateUser(ctx context.Context,
	request dto.CreateUserRequest,
) (dto.CreateUserResponse, error) {
	response, err := s.userServiceClient.CreateUser(ctx, request)
	if err != nil {
		return dto.CreateUserResponse{}, fmt.Errorf("create user: %w", err)
	}

	return response, nil
}
