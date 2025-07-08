package service

import (
	"context"
	"fmt"

	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/dto"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/model"
)

type ListingViewRepository interface {
	GetAll(ctx context.Context, limit, offset int, userID *int64) ([]model.Listing, error)
}

type ListingViewService struct {
	listingRepository ListingViewRepository
}

func NewListingViewService(listingRepository ListingViewRepository) *ListingViewService {
	return &ListingViewService{
		listingRepository: listingRepository,
	}
}

func (s *ListingViewService) GetAllListings(ctx context.Context, req dto.GetAllListingsRequest) (dto.GetAllListingsResponse, error) {
	limit := req.PageSize
	offset := (req.PageNum - 1) * req.PageSize
	userID := req.UserID

	result, err := s.listingRepository.GetAll(ctx, limit, offset, userID)
	if err != nil {
		return dto.GetAllListingsResponse{}, fmt.Errorf("failed to get all listings: %w", err)
	}

	listings := make([]dto.ListingResponse, len(result))
	for i, listing := range result {
		listings[i] = dto.ListingResponse{
			ID:          listing.ID,
			ListingType: listing.ListingType,
			Price:       listing.Price,
			CreatedAt:   listing.CreatedAt,
			UpdatedAt:   listing.UpdatedAt,
			User: dto.UserResponse{
				ID:        listing.User.ID,
				Name:      listing.User.Name,
				CreatedAt: listing.User.CreatedAt,
				UpdatedAt: listing.User.UpdatedAt,
			},
		}
	}

	return dto.GetAllListingsResponse{
		Result:   true,
		Listings: listings,
	}, nil
}
