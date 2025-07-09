//go:build unit

package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/dto"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetAllUsers(t *testing.T) {
	getAllUsersRequest := func(name string, req dto.GetAllUsersRequest, mockRepo *MockUserRepository, want dto.GetAllUsersResponse) func(t *testing.T) {
		return func(t *testing.T) {
			svc := NewUserService(mockRepo, &MockPublisher{})
			got, err := svc.GetAllUsers(context.Background(), req)
			if mockRepo.err != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}

	t.Run("success", getAllUsersRequest(
		"success",
		dto.GetAllUsersRequest{PageSize: 10, PageNumber: 1},
		&MockUserRepository{users: mockUsers},
		dto.GetAllUsersResponse{
			Result: true,
			Users: []dto.UserResponse{
				{
					ID:        mockUsers[0].ID,
					Name:      mockUsers[0].Name,
					CreatedAt: mockUsers[0].CreatedAt,
					UpdatedAt: mockUsers[0].UpdatedAt,
				},
				{
					ID:        mockUsers[1].ID,
					Name:      mockUsers[1].Name,
					CreatedAt: mockUsers[1].CreatedAt,
					UpdatedAt: mockUsers[1].UpdatedAt,
				},
			},
		},
	))

	t.Run("db_error", getAllUsersRequest(
		"db_error",
		dto.GetAllUsersRequest{PageSize: 10, PageNumber: 1},
		&MockUserRepository{err: ErrMockDB},
		dto.GetAllUsersResponse{},
	))
}

func TestUserService_GetUserByID(t *testing.T) {
	getUserByIDRequest := func(name string, req dto.GetUserByIDRequest, mockRepo *MockUserRepository, want dto.GetUserByIDResponse, wantErr error) func(t *testing.T) {
		return func(t *testing.T) {
			svc := NewUserService(mockRepo, &MockPublisher{})
			got, err := svc.GetUserByID(context.Background(), req)
			if wantErr != nil {
				assert.ErrorIs(t, err, wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}

	t.Run("success", getUserByIDRequest(
		"success",
		dto.GetUserByIDRequest{ID: 1},
		&MockUserRepository{users: mockUsers},
		dto.GetUserByIDResponse{
			Result: true,
			User: dto.UserResponse{
				ID:        mockUsers[0].ID,
				Name:      mockUsers[0].Name,
				CreatedAt: mockUsers[0].CreatedAt,
				UpdatedAt: mockUsers[0].UpdatedAt,
			},
		},
		nil,
	))

	t.Run("not_found", getUserByIDRequest(
		"not_found",
		dto.GetUserByIDRequest{ID: 999},
		&MockUserRepository{users: mockUsers},
		dto.GetUserByIDResponse{},
		sql.ErrNoRows,
	))

	t.Run("db_error", getUserByIDRequest(
		"db_error",
		dto.GetUserByIDRequest{ID: 1},
		&MockUserRepository{err: ErrMockDB},
		dto.GetUserByIDResponse{},
		ErrMockDB,
	))
}

func TestUserService_CreateUser(t *testing.T) {
	createUserRequest := func(name string, req dto.CreateUserRequest, mockRepo *MockUserRepository, mockPub *MockPublisher, want dto.CreateUserResponse) func(t *testing.T) {
		return func(t *testing.T) {
			svc := NewUserService(mockRepo, mockPub)
			got, err := svc.CreateUser(context.Background(), req)
			if mockRepo.err != nil || mockPub.err != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			// Since CreatedAt and UpdatedAt are set at runtime, we only compare ID and Name
			assert.Equal(t, want.Result, got.Result)
			assert.Equal(t, want.User.Name, got.User.Name)
			assert.NotZero(t, got.User.ID)
			assert.NotZero(t, got.User.CreatedAt)
			assert.NotZero(t, got.User.UpdatedAt)
		}
	}

	t.Run("success", createUserRequest(
		"success",
		dto.CreateUserRequest{Name: "Test User"},
		&MockUserRepository{users: mockUsers},
		&MockPublisher{},
		dto.CreateUserResponse{
			Result: true,
			User: dto.UserResponse{
				Name: "Test User",
			},
		},
	))

	t.Run("db_error", createUserRequest(
		"db_error",
		dto.CreateUserRequest{Name: "Test User"},
		&MockUserRepository{err: ErrMockDB},
		&MockPublisher{},
		dto.CreateUserResponse{},
	))

	t.Run("publish_error", createUserRequest(
		"publish_error",
		dto.CreateUserRequest{Name: "Test User"},
		&MockUserRepository{users: mockUsers},
		&MockPublisher{err: ErrMockPublish},
		dto.CreateUserResponse{},
	))
}
