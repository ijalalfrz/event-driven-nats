package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/dto"
)

type ListingViewServiceClient struct {
	httpClient HTTPClient
}

func NewListingViewServiceClient(
	serviceURL string,
	opts ...ClientOption,
) *ListingViewServiceClient {
	client := &ListingViewServiceClient{
		httpClient: HTTPClient{
			client: http.DefaultClient,
			url:    serviceURL,
		},
	}

	for _, opt := range opts {
		opt(&client.httpClient)
	}

	return client
}

func (c *ListingViewServiceClient) GetAllListings(ctx context.Context,
	request dto.GetAllListingsRequest,
) (dto.GetAllListingsResponse, error) {
	var response dto.GetAllListingsResponse

	values := url.Values{}
	values.Add("page_num", fmt.Sprintf("%d", request.PageNumber))
	values.Add("page_size", fmt.Sprintf("%d", request.PageSize))
	if request.UserID != nil {
		values.Add("user_id", fmt.Sprintf("%d", *request.UserID))
	}
	path := fmt.Sprintf("/listings?%s", values.Encode())

	headerFunc := func(req *http.Request) {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.httpClient.doRequestWithResponse(ctx, http.MethodGet, path, headerFunc,
		nil, defaultErrorResponseFunc)
	if err != nil {
		return dto.GetAllListingsResponse{}, fmt.Errorf("get listings request failed: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return dto.GetAllListingsResponse{}, fmt.Errorf("decode response: %w", err)
	}

	return response, nil
}
