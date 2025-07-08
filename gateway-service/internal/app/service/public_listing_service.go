package service

import (
	"context"
	"fmt"

	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/dto"
)

type PublicListingService struct {
	listingViewServiceClient *ListingViewServiceClient
	listingServiceClient     *ListingServiceClient
	userServiceClient        *UserServiceClient
}

func NewPublicListingService(
	listingViewServiceClient *ListingViewServiceClient,
	listingServiceClient *ListingServiceClient,
	userServiceClient *UserServiceClient,
) *PublicListingService {
	return &PublicListingService{
		listingViewServiceClient: listingViewServiceClient,
		listingServiceClient:     listingServiceClient,
		userServiceClient:        userServiceClient,
	}
}

func (s *PublicListingService) CreateListing(ctx context.Context,
	request dto.CreateListingRequest,
) (dto.CreateListingResponse, error) {
	_, err := s.userServiceClient.GetUserByID(ctx, request.UserID)
	if err != nil {
		return dto.CreateListingResponse{}, fmt.Errorf("get user: %w", err)
	}

	response, err := s.listingServiceClient.CreateListing(ctx, request)
	if err != nil {
		return dto.CreateListingResponse{}, fmt.Errorf("create listing: %w", err)
	}

	return response, nil
}

func (s *PublicListingService) GetAllListings(ctx context.Context,
	request dto.GetAllListingsRequest,
) (dto.GetAllListingsResponse, error) {
	response, err := s.listingViewServiceClient.GetAllListings(ctx, request)
	if err != nil {
		return dto.GetAllListingsResponse{}, fmt.Errorf("get all listings: %w", err)
	}

	return response, nil
}
