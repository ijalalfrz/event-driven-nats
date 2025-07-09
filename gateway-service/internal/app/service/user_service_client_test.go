//go:build unit

package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/dto"
	"github.com/stretchr/testify/assert"
)

func TestUserServiceClient_CreateUser_Positive(t *testing.T) {
	createUserRequest := func(request dto.CreateUserRequest, want dto.CreateUserResponse) func(t *testing.T) {
		return func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `{
					"user": {
						"id": 1,
						"name": "John Doe",
						"created_at": 1234567890,
						"updated_at": 1234567890
					}
				}`)
			}))
			defer server.Close()

			subject := NewUserServiceClient(server.URL, WithMaxRetries(1))
			got, err := subject.CreateUser(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}

	t.Run("success_create_user", createUserRequest(
		dto.CreateUserRequest{
			Name: "John Doe",
		},
		dto.CreateUserResponse{
			UserResponse: dto.UserResponse{
				ID:        1,
				Name:      "John Doe",
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
	))
}

func TestUserServiceClient_CreateUser_Negative(t *testing.T) {
	createUserRequest := func(request dto.CreateUserRequest, wantErr string) func(t *testing.T) {
		return func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `{
					"error": "invalid request"
				}`)
			}))
			defer server.Close()

			subject := NewUserServiceClient(server.URL, WithMaxRetries(1))
			_, err := subject.CreateUser(context.Background(), request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), wantErr)
		}
	}

	t.Run("invalid_request", createUserRequest(
		dto.CreateUserRequest{
			Name: "",
		},
		"invalid request",
	))
}

func TestUserServiceClient_GetUserByID_Positive(t *testing.T) {
	getUserByIDRequest := func(userID int64, want dto.UserResponse) func(t *testing.T) {
		return func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `{
					"id": 1,
					"name": "John Doe",
					"created_at": 1234567890,
					"updated_at": 1234567890
				}`)
			}))
			defer server.Close()

			subject := NewUserServiceClient(server.URL, WithMaxRetries(1))
			got, err := subject.GetUserByID(context.Background(), userID)

			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}

	t.Run("success_get_user", getUserByIDRequest(
		1,
		dto.UserResponse{
			ID:        1,
			Name:      "John Doe",
			CreatedAt: 1234567890,
			UpdatedAt: 1234567890,
		},
	))
}

func TestUserServiceClient_GetUserByID_Negative(t *testing.T) {
	getUserByIDRequest := func(userID int64, wantErr string) func(t *testing.T) {
		return func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, `{
					"error": "user not found"
				}`)
			}))
			defer server.Close()

			subject := NewUserServiceClient(server.URL, WithMaxRetries(1))
			_, err := subject.GetUserByID(context.Background(), userID)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), wantErr)
		}
	}

	t.Run("user_not_found", getUserByIDRequest(
		999,
		"user not found",
	))
}
