package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/dto"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/model"
	"github.com/nats-io/nats.go/jetstream"
)

type UserRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]model.User, error)
	GetByID(ctx context.Context, id int64) (model.User, error)
	CreateTx(ctx context.Context, tx *sql.Tx, user *model.User) error
	WithTransaction(ctx context.Context,
		txFunc func(context.Context, *sql.Tx) error,
	) error
}

type Publisher interface {
	Publish(ctx context.Context, subject string, request interface{}) (*jetstream.PubAck, error)
}

type UserService struct {
	userRepository UserRepository
	publisher      Publisher
}

func NewUserService(userRepository UserRepository,
	publisher Publisher) *UserService {
	return &UserService{
		userRepository: userRepository,
		publisher:      publisher,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	user := model.User{
		Name:      req.Name,
		CreatedAt: time.Now().UnixMicro(),
		UpdatedAt: time.Now().UnixMicro(),
	}

	err := s.userRepository.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := s.userRepository.CreateTx(ctx, tx, &user)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// publish event
		_, err = s.publisher.Publish(ctx, model.UserCreatedEvent, user)
		if err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}

		return nil
	})
	if err != nil {
		return dto.CreateUserResponse{}, fmt.Errorf("failed to create user with transaction: %w", err)
	}

	return dto.CreateUserResponse{
		Result: true,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req dto.GetAllUsersRequest) (dto.GetAllUsersResponse, error) {
	limit := req.PageSize
	offset := (req.PageNumber - 1) * req.PageSize

	users, err := s.userRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return dto.GetAllUsersResponse{}, err
	}

	usersResponse := make([]dto.UserResponse, len(users))
	for i, user := range users {
		usersResponse[i] = dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return dto.GetAllUsersResponse{
		Result: true,
		Users:  usersResponse,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, req dto.GetUserByIDRequest) (dto.GetUserByIDResponse, error) {
	user, err := s.userRepository.GetByID(ctx, req.ID)
	if err != nil {
		return dto.GetUserByIDResponse{}, err
	}

	return dto.GetUserByIDResponse{
		Result: true,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}
