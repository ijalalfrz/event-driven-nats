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

// CreateListing godoc
// @Summary      Create Listing
// @Description  Create a Listing
// @Tags         Listing
// @ID           createListing
// @Produce      json
// @Param        req body create listing	body		dto.CreateListingRequest	true	"Listing"
// @Success      200  {object}  dto.CreateListingResponse	"Created"
// @Failure      400  {object}  dto.ErrorResponse	"Bad Request"
// @Failure      500  {object}  dto.ErrorResponse	"Internal Server Error"
// @Router       /public/listings [post].
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

// GetAllListings godoc
// @Summary      Get All Listings
// @Description  Get All Listings
// @Tags         Listing
// @ID           getAllListings
// @Produce      json
// @Param        req body get all listings	body		dto.GetAllListingsRequest	true	"Listing"
// @Success      200  {object}  dto.GetAllListingsResponse	"Listings"
// @Failure      400  {object}  dto.ErrorResponse	"Bad Request"
// @Failure      500  {object}  dto.ErrorResponse	"Internal Server Error"
// @Router       /public/listings [get].
func (s *PublicListingService) GetAllListings(ctx context.Context,
	request dto.GetAllListingsRequest,
) (dto.GetAllListingsResponse, error) {
	response, err := s.listingViewServiceClient.GetAllListings(ctx, request)
	if err != nil {
		return dto.GetAllListingsResponse{}, fmt.Errorf("get all listings: %w", err)
	}

	return response, nil
}
