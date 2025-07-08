package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/dto"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/model"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (model.User, error)
	CreateTx(ctx context.Context, tx *sql.Tx, user *model.User) error
	WithTransaction(ctx context.Context,
		txFunc func(context.Context, *sql.Tx) error,
	) error
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) OnCreatedUser(ctx context.Context, req dto.UserCreated) error {
	err := s.userRepo.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		user := &model.User{
			ID:        req.ID,
			Name:      req.Name,
			CreatedAt: req.CreatedAt,
			UpdatedAt: req.UpdatedAt,
		}

		err := s.userRepo.CreateTx(ctx, tx, user)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("on created user: %w", err)
	}

	return nil
}
