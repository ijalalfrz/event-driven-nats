//go:build unit

package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/dto"
	"github.com/stretchr/testify/assert"
)

func TestListingService_OnCreatedListing(t *testing.T) {
	onCreatedListing := func(name string, req dto.ListingCreated, mockListingRepo *MockListingRepository, mockUserRepo *MockUserRepository, wantErr error) func(t *testing.T) {
		return func(t *testing.T) {
			svc := NewListingService(mockListingRepo, mockUserRepo)
			err := svc.OnCreatedListing(context.Background(), req)
			if wantErr != nil {
				assert.ErrorIs(t, err, wantErr)
				return
			}
			assert.NoError(t, err)

			// Verify listing was created with correct data
			assert.Len(t, mockListingRepo.listings, 1)
			listing := mockListingRepo.listings[0]
			assert.Equal(t, req.ID, listing.ID)
			assert.Equal(t, req.UserID, listing.UserID)
			assert.Equal(t, req.ListingType, listing.ListingType)
			assert.Equal(t, req.Price, listing.Price)
			assert.Equal(t, req.CreatedAt, listing.CreatedAt)
			assert.Equal(t, req.UpdatedAt, listing.UpdatedAt)

			// Verify user data was included
			assert.Equal(t, mockUsers[0].ID, listing.User.ID)
			assert.Equal(t, mockUsers[0].Name, listing.User.Name)
			assert.Equal(t, mockUsers[0].CreatedAt, listing.User.CreatedAt)
			assert.Equal(t, mockUsers[0].UpdatedAt, listing.User.UpdatedAt)
		}
	}

	t.Run("success", onCreatedListing(
		"success",
		dto.ListingCreated{
			ID:          1,
			UserID:      1,
			ListingType: "SALE",
			Price:       1000,
			CreatedAt:   1234567890,
			UpdatedAt:   1234567890,
		},
		&MockListingRepository{},
		&MockUserRepository{users: mockUsers},
		nil,
	))

	t.Run("user_not_found", onCreatedListing(
		"user_not_found",
		dto.ListingCreated{
			ID:          1,
			UserID:      999, // non-existent user
			ListingType: "SALE",
			Price:       1000,
			CreatedAt:   1234567890,
			UpdatedAt:   1234567890,
		},
		&MockListingRepository{},
		&MockUserRepository{users: mockUsers},
		sql.ErrNoRows,
	))

	t.Run("db_error_on_create", onCreatedListing(
		"db_error_on_create",
		dto.ListingCreated{
			ID:          1,
			UserID:      1,
			ListingType: "SALE",
			Price:       1000,
			CreatedAt:   1234567890,
			UpdatedAt:   1234567890,
		},
		&MockListingRepository{err: ErrMockDB},
		&MockUserRepository{users: mockUsers},
		ErrMockDB,
	))

	t.Run("db_error_on_get_user", onCreatedListing(
		"db_error_on_get_user",
		dto.ListingCreated{
			ID:          1,
			UserID:      1,
			ListingType: "SALE",
			Price:       1000,
			CreatedAt:   1234567890,
			UpdatedAt:   1234567890,
		},
		&MockListingRepository{},
		&MockUserRepository{err: ErrMockDB},
		ErrMockDB,
	))
}
