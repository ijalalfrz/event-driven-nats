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

// CreateUser godoc
// @Summary      Create User
// @Description  Create a User
// @Tags         User
// @ID           createUser
// @Produce      json
// @Param        req body create user	body		dto.CreateUserRequest	true	"User"
// @Success      200  {object}  dto.CreateUserResponse	"User"
// @Failure      400  {object}  dto.ErrorResponse	"Bad Request"
// @Failure      500  {object}  dto.ErrorResponse	"Internal Server Error"
// @Router       /public/users [post].
func (s *PublicUserService) CreateUser(ctx context.Context,
	request dto.CreateUserRequest,
) (dto.CreateUserResponse, error) {
	response, err := s.userServiceClient.CreateUser(ctx, request)
	if err != nil {
		return dto.CreateUserResponse{}, fmt.Errorf("create user: %w", err)
	}

	return response, nil
}
