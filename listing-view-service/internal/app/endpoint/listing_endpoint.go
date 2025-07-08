package endpoint

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/dto"
)

type ListingViewService interface {
	GetAllListings(ctx context.Context, req dto.GetAllListingsRequest) (dto.GetAllListingsResponse, error)
}

type ListingService interface {
	OnCreatedListing(ctx context.Context, req dto.ListingCreated) error
}

func NewListingEndpoint(svc ListingViewService, listingSvc ListingService) Listing {
	return Listing{
		GetAll:    MakeGetAllListingsEndpoint(svc),
		OnCreated: MakeOnCreatedListingEndpoint(listingSvc),
	}
}

func MakeOnCreatedListingEndpoint(svc ListingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*dto.ListingCreated)
		if !ok {
			return nil, fmt.Errorf("listing service: %w", ErrInvalidType)
		}

		err := svc.OnCreatedListing(ctx, *req)
		if err != nil {
			return nil, fmt.Errorf("listing service: %w", err)
		}

		return nil, nil
	}
}

func MakeGetAllListingsEndpoint(svc ListingViewService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*dto.GetAllListingsRequest)
		if !ok {
			return nil, fmt.Errorf("listing view service: %w", ErrInvalidType)
		}

		res, err := svc.GetAllListings(ctx, *req)
		if err != nil {
			return nil, fmt.Errorf("listing view service: %w", err)
		}

		return res, nil
	}
}
