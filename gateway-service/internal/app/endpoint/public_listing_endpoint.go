package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/dto"
)

type PublicListingService interface {
	CreateListing(ctx context.Context, request dto.CreateListingRequest) (dto.CreateListingResponse, error)
	GetAllListings(ctx context.Context, request dto.GetAllListingsRequest) (dto.GetAllListingsResponse, error)
}

func NewPublicListingEndpoint(
	service PublicListingService,
) PublicListing {
	return PublicListing{
		Create: makeCreateListingEndpoint(service),
		GetAll: makeGetAllListingsEndpoint(service),
	}
}

func makeCreateListingEndpoint(service PublicListingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*dto.CreateListingRequest)
		if !ok {
			return nil, ErrInvalidType
		}

		response, err := service.CreateListing(ctx, *req)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}

func makeGetAllListingsEndpoint(service PublicListingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*dto.GetAllListingsRequest)
		if !ok {
			return nil, ErrInvalidType
		}

		response, err := service.GetAllListings(ctx, *req)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}
