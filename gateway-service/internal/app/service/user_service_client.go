package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/dto"
)

type UserServiceClient struct {
	httpClient HTTPClient
}

func NewUserServiceClient(
	serviceURL string,
	opts ...ClientOption,
) *UserServiceClient {
	client := &UserServiceClient{
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

func (c *UserServiceClient) CreateUser(ctx context.Context,
	request dto.CreateUserRequest,
) (dto.CreateUserResponse, error) {
	var response dto.CreateUserResponse

	path := "/users"

	headerFunc := func(req *http.Request) {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	// Convert request to form values
	formData := url.Values{}
	formData.Add("name", request.Name)

	resp, err := c.httpClient.doRequestWithResponse(ctx, http.MethodPost, path, headerFunc,
		formData.Encode(), defaultErrorResponseFunc)
	if err != nil {
		return dto.CreateUserResponse{}, fmt.Errorf("create user request failed: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return dto.CreateUserResponse{}, fmt.Errorf("decode response: %w", err)
	}

	return response, nil
}

func (c *UserServiceClient) GetUserByID(ctx context.Context, userID int64) (dto.UserResponse, error) {
	var response dto.UserResponse

	path := fmt.Sprintf("/users/%d", userID)

	headerFunc := func(req *http.Request) {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.httpClient.doRequestWithResponse(ctx, http.MethodGet, path, headerFunc,
		"", defaultErrorResponseFunc)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("get user request failed: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("decode response: %w", err)
	}

	return response, nil
}
