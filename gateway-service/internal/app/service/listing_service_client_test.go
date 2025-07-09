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

func TestListingServiceClient_CreateListing_Positive(t *testing.T) {
	createListingRequest := func(request dto.CreateListingRequest, want dto.CreateListingResponse) func(t *testing.T) {
		return func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `{
					"listing": {
						"id": 1,
						"user_id": 1,
						"price": 1000,
						"listing_type": "rent",
						"created_at": 1234567890,
						"updated_at": 1234567890
					}
				}`)
			}))
			defer server.Close()

			subject := NewListingServiceClient(server.URL, WithMaxRetries(1))
			got, err := subject.CreateListing(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}

	t.Run("success_create_listing", createListingRequest(
		dto.CreateListingRequest{
			UserID:      1,
			Price:       1000,
			ListingType: "rent",
		},
		dto.CreateListingResponse{
			ListingResponse: dto.ListingResponse{
				ID:          1,
				UserID:      1,
				Price:       1000,
				ListingType: "rent",
				CreatedAt:   1234567890,
				UpdatedAt:   1234567890,
			},
		},
	))
}

func TestListingServiceClient_CreateListing_Negative(t *testing.T) {
	createListingRequest := func(request dto.CreateListingRequest, wantErr string) func(t *testing.T) {
		return func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `{
					"errors": ["invalid request"]
				}`)
			}))
			defer server.Close()

			subject := NewListingServiceClient(server.URL, WithMaxRetries(1))
			_, err := subject.CreateListing(context.Background(), request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), wantErr)
		}
	}

	t.Run("invalid_request", createListingRequest(
		dto.CreateListingRequest{
			UserID:      0,
			Price:       0,
			ListingType: "",
		},
		"invalid request",
	))
}
