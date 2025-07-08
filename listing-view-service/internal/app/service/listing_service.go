package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/dto"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/model"
)

type ListingRepository interface {
	CreateTx(ctx context.Context, tx *sql.Tx, listing *model.Listing) error
	WithTransaction(ctx context.Context,
		txFunc func(context.Context, *sql.Tx) error,
	) error
}

type ListingService struct {
	listingsRepo ListingRepository
	userRepo     UserRepository
}

func NewListingService(listingsRepo ListingRepository, userRepo UserRepository) *ListingService {
	return &ListingService{
		listingsRepo: listingsRepo,
		userRepo:     userRepo,
	}
}

func (s *ListingService) OnCreatedListing(ctx context.Context, req dto.ListingCreated) error {
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	err = s.listingsRepo.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		listing := &model.Listing{
			ID:          req.ID,
			UserID:      req.UserID,
			ListingType: req.ListingType,
			Price:       req.Price,
			CreatedAt:   req.CreatedAt,
			UpdatedAt:   req.UpdatedAt,
			User:        user,
		}

		err := s.listingsRepo.CreateTx(ctx, tx, listing)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("on created listing: %w", err)
	}

	return nil
}
