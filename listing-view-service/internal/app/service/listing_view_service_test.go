//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/dto"
	"github.com/stretchr/testify/assert"
)

func TestListingViewService_GetAllListings(t *testing.T) {
	getAllListings := func(name string, req dto.GetAllListingsRequest, mockRepo *MockListingViewRepository, want dto.GetAllListingsResponse, wantErr error) func(t *testing.T) {
		return func(t *testing.T) {
			svc := NewListingViewService(mockRepo)
			got, err := svc.GetAllListings(context.Background(), req)
			if wantErr != nil {
				assert.ErrorIs(t, err, wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}

	t.Run("success_all_listings", getAllListings(
		"success_all_listings",
		dto.GetAllListingsRequest{
			PageNum:  1,
			PageSize: 10,
		},
		&MockListingViewRepository{listings: mockListings},
		dto.GetAllListingsResponse{
			Result: true,
			Listings: []dto.ListingResponse{
				{
					ID:          mockListings[0].ID,
					ListingType: mockListings[0].ListingType,
					Price:       mockListings[0].Price,
					CreatedAt:   mockListings[0].CreatedAt,
					UpdatedAt:   mockListings[0].UpdatedAt,
					User: dto.UserResponse{
						ID:        mockListings[0].User.ID,
						Name:      mockListings[0].User.Name,
						CreatedAt: mockListings[0].User.CreatedAt,
						UpdatedAt: mockListings[0].User.UpdatedAt,
					},
				},
				{
					ID:          mockListings[1].ID,
					ListingType: mockListings[1].ListingType,
					Price:       mockListings[1].Price,
					CreatedAt:   mockListings[1].CreatedAt,
					UpdatedAt:   mockListings[1].UpdatedAt,
					User: dto.UserResponse{
						ID:        mockListings[1].User.ID,
						Name:      mockListings[1].User.Name,
						CreatedAt: mockListings[1].User.CreatedAt,
						UpdatedAt: mockListings[1].User.UpdatedAt,
					},
				},
			},
		},
		nil,
	))

	userID := int64(1)
	t.Run("success_filter_by_user", getAllListings(
		"success_filter_by_user",
		dto.GetAllListingsRequest{
			PageNum:  1,
			PageSize: 10,
			UserID:   &userID,
		},
		&MockListingViewRepository{listings: mockListings},
		dto.GetAllListingsResponse{
			Result: true,
			Listings: []dto.ListingResponse{
				{
					ID:          mockListings[0].ID,
					ListingType: mockListings[0].ListingType,
					Price:       mockListings[0].Price,
					CreatedAt:   mockListings[0].CreatedAt,
					UpdatedAt:   mockListings[0].UpdatedAt,
					User: dto.UserResponse{
						ID:        mockListings[0].User.ID,
						Name:      mockListings[0].User.Name,
						CreatedAt: mockListings[0].User.CreatedAt,
						UpdatedAt: mockListings[0].User.UpdatedAt,
					},
				},
			},
		},
		nil,
	))

	t.Run("db_error", getAllListings(
		"db_error",
		dto.GetAllListingsRequest{
			PageNum:  1,
			PageSize: 10,
		},
		&MockListingViewRepository{err: ErrMockDB},
		dto.GetAllListingsResponse{},
		ErrMockDB,
	))
}
