//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/dto"
	"github.com/stretchr/testify/assert"
)

func TestUserService_OnCreatedUser(t *testing.T) {
	onCreatedUser := func(name string, req dto.UserCreated, mockRepo *MockUserRepository, wantErr error) func(t *testing.T) {
		return func(t *testing.T) {
			svc := NewUserService(mockRepo)
			err := svc.OnCreatedUser(context.Background(), req)
			if wantErr != nil {
				assert.ErrorIs(t, err, wantErr)
				return
			}
			assert.NoError(t, err)

			// Verify user was created
			user, err := mockRepo.GetByID(context.Background(), req.ID)
			assert.NoError(t, err)
			assert.Equal(t, req.ID, user.ID)
			assert.Equal(t, req.Name, user.Name)
			assert.Equal(t, req.CreatedAt, user.CreatedAt)
			assert.Equal(t, req.UpdatedAt, user.UpdatedAt)
		}
	}

	t.Run("success", onCreatedUser(
		"success",
		dto.UserCreated{
			ID:        3,
			Name:      "Test User",
			CreatedAt: 1234567890,
			UpdatedAt: 1234567890,
		},
		&MockUserRepository{users: mockUsers},
		nil,
	))

	t.Run("db_error", onCreatedUser(
		"db_error",
		dto.UserCreated{
			ID:        3,
			Name:      "Test User",
			CreatedAt: 1234567890,
			UpdatedAt: 1234567890,
		},
		&MockUserRepository{err: ErrMockDB},
		ErrMockDB,
	))
}
