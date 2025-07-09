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

func TestListingViewServiceClient_GetAllListings_Positive(t *testing.T) {
	getAllListingsRequest := func(request dto.GetAllListingsRequest, want dto.GetAllListingsResponse) func(t *testing.T) {
		return func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `{
					"result": true,
					"listings": [
						{
							"id": 1,
							"user_id": 1,
							"price": 1000,
							"listing_type": "rent",
							"created_at": 1234567890,
							"updated_at": 1234567890
						}
					]
				}`)
			}))
			defer server.Close()

			subject := NewListingViewServiceClient(server.URL, WithMaxRetries(1))
			got, err := subject.GetAllListings(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}

	userID := int64(1)
	t.Run("success_get_all_listings", getAllListingsRequest(
		dto.GetAllListingsRequest{
			PageNumber: 1,
			PageSize:   10,
			UserID:     &userID,
		},
		dto.GetAllListingsResponse{
			Result: true,
			Listings: []dto.ListingResponse{
				{
					ID:          1,
					UserID:      1,
					Price:       1000,
					ListingType: "rent",
					CreatedAt:   1234567890,
					UpdatedAt:   1234567890,
				},
			},
		},
	))
}

func TestListingViewServiceClient_GetAllListings_Negative(t *testing.T) {
	getAllListingsRequest := func(request dto.GetAllListingsRequest, wantErr string) func(t *testing.T) {
		return func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `{
					"error": "invalid request"
				}`)
			}))
			defer server.Close()

			subject := NewListingViewServiceClient(server.URL, WithMaxRetries(1))
			_, err := subject.GetAllListings(context.Background(), request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), wantErr)
		}
	}

	t.Run("invalid_request", getAllListingsRequest(
		dto.GetAllListingsRequest{
			PageNumber: 0,
			PageSize:   0,
		},
		"invalid request",
	))
}
