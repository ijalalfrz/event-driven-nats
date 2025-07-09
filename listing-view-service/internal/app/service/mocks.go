package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/model"
)

// Mock errors
var (
	ErrMockDB = errors.New("mock db error")
)

// MockUserRepository implements UserRepository interface
type MockUserRepository struct {
	users []model.User
	err   error
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (model.User, error) {
	if m.err != nil {
		return model.User{}, m.err
	}
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return model.User{}, sql.ErrNoRows
}

func (m *MockUserRepository) CreateTx(ctx context.Context, tx *sql.Tx, user *model.User) error {
	if m.err != nil {
		return m.err
	}
	m.users = append(m.users, *user)
	return nil
}

func (m *MockUserRepository) WithTransaction(ctx context.Context, txFunc func(context.Context, *sql.Tx) error) error {
	if m.err != nil {
		return m.err
	}
	return txFunc(ctx, &sql.Tx{})
}

// MockListingViewRepository implements ListingViewRepository interface
type MockListingViewRepository struct {
	listings []model.Listing
	err      error
}

func (m *MockListingViewRepository) GetAll(ctx context.Context, limit, offset int, userID *int64) ([]model.Listing, error) {
	if m.err != nil {
		return nil, m.err
	}

	if userID != nil {
		// Filter by userID
		filtered := make([]model.Listing, 0)
		for _, listing := range m.listings {
			if listing.User.ID == *userID {
				filtered = append(filtered, listing)
			}
		}
		return filtered, nil
	}

	return m.listings, nil
}

// MockListingRepository implements ListingRepository interface
type MockListingRepository struct {
	listings []model.Listing
	err      error
}

func (m *MockListingRepository) CreateTx(ctx context.Context, tx *sql.Tx, listing *model.Listing) error {
	if m.err != nil {
		return m.err
	}
	m.listings = append(m.listings, *listing)
	return nil
}

func (m *MockListingRepository) WithTransaction(ctx context.Context, txFunc func(context.Context, *sql.Tx) error) error {
	if m.err != nil {
		return m.err
	}
	return txFunc(ctx, &sql.Tx{})
}

// Test data
var mockUsers = []model.User{
	{
		ID:        1,
		Name:      "John Doe",
		CreatedAt: 1234567890,
		UpdatedAt: 1234567890,
	},
	{
		ID:        2,
		Name:      "Jane Doe",
		CreatedAt: 1234567890,
		UpdatedAt: 1234567890,
	},
}

// Test data
var mockListings = []model.Listing{
	{
		ID:          1,
		ListingType: "SALE",
		Price:       1000,
		CreatedAt:   1234567890,
		UpdatedAt:   1234567890,
		User: model.User{
			ID:        1,
			Name:      "John Doe",
			CreatedAt: 1234567890,
			UpdatedAt: 1234567890,
		},
	},
	{
		ID:          2,
		ListingType: "RENT",
		Price:       500,
		CreatedAt:   1234567890,
		UpdatedAt:   1234567890,
		User: model.User{
			ID:        2,
			Name:      "Jane Doe",
			CreatedAt: 1234567890,
			UpdatedAt: 1234567890,
		},
	},
}
