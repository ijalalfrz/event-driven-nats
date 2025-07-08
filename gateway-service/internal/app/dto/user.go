package dto

import (
	"fmt"
	"net/http"
)

type CreateUserRequest struct {
	Name string `json:"name" validate:"required"`
}

func (r *CreateUserRequest) Bind(req *http.Request) error {
	if err := validate.Struct(r); err != nil {
		return NewInvalidRequestError(fmt.Errorf("invalid request: %w", err))
	}

	return nil
}

type CreateUserResponse struct {
	UserResponse `json:"user"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
