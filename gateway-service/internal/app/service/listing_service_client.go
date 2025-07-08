package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/dto"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/pkg/exception"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/pkg/lang"
)

type ListingErrorResponse struct {
	Errors []string `json:"errors"`
}

type ListingServiceClient struct {
	httpClient HTTPClient
}

func NewListingServiceClient(
	serviceURL string,
	opts ...ClientOption,
) *ListingServiceClient {
	client := &ListingServiceClient{
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

func (c *ListingServiceClient) CreateListing(ctx context.Context,
	request dto.CreateListingRequest,
) (dto.CreateListingResponse, error) {
	var response dto.CreateListingResponse

	path := "/listings"

	headerFunc := func(req *http.Request) {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	formData := url.Values{}
	formData.Add("price", fmt.Sprintf("%d", request.Price))
	formData.Add("user_id", fmt.Sprintf("%d", request.UserID))
	formData.Add("listing_type", request.ListingType)

	resp, err := c.httpClient.doRequestWithResponse(ctx, http.MethodPost, path, headerFunc,
		formData.Encode(), listingErrorResponseFunc)
	if err != nil {
		return dto.CreateListingResponse{}, fmt.Errorf("create listing request failed: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return dto.CreateListingResponse{}, fmt.Errorf("decode response: %w", err)
	}

	return response, nil
}

func listingErrorResponseFunc(resp *http.Response) error { //nolint:unused
	var errorResp ListingErrorResponse

	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		return fmt.Errorf("decode error response: %w", err)
	}

	// internal server errors should not be forwarded to the client
	if resp.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("returned status code: %d", resp.StatusCode)
	}

	host := strings.Split(resp.Request.URL.Host, ":")[0]
	errJoined := strings.Join(errorResp.Errors, ", ")

	return exception.ApplicationError{
		StatusCode: resp.StatusCode,
		Localizable: lang.Localizable{
			MessageID: "errors.bad_request_from_service",
			MessageVars: map[string]interface{}{
				"service": host,
				"error":   errJoined,
			},
		},
		Cause: fmt.Errorf("error response: %s", errJoined),
	}
}
