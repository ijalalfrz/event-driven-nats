package dto

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreateUserRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

func (r *CreateUserRequest) Bind(req *http.Request) error {
	if err := validate.Struct(r); err != nil {
		return NewInvalidRequestError(err)
	}

	return nil
}

type GetAllUsersRequest struct {
	PageNumber int `json:"page_number" validate:"required,min=1"`
	PageSize   int `json:"page_size" validate:"required,min=1"`
}

func (r *GetAllUsersRequest) Bind(req *http.Request) error {
	var err error

	pageNumberStr := req.URL.Query().Get("page_num")
	pageNumber := 1
	if pageNumberStr != "" {
		pageNumber, err = strconv.Atoi(pageNumberStr)
		if err != nil {
			return NewInvalidRequestError(fmt.Errorf("invalid page number: %w", err))
		}
	}

	r.PageNumber = pageNumber

	pageSizeStr := req.URL.Query().Get("page_size")
	pageSize := 10
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			return NewInvalidRequestError(fmt.Errorf("invalid page size: %w", err))
		}
	}

	r.PageSize = pageSize

	if err := validate.Struct(r); err != nil {
		return NewInvalidRequestError(err)
	}

	return nil
}

type GetUserByIDRequest struct {
	ID int64 `json:"id" validate:"required"`
}

func (r *GetUserByIDRequest) Bind(req *http.Request) error {
	id, err := strconv.ParseInt(chi.URLParam(req, "id"), 10, 64)
	if err != nil {
		return NewInvalidRequestError(err)
	}

	r.ID = id

	if err := validate.Struct(r); err != nil {
		return NewInvalidRequestError(err)
	}

	return nil
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type CreateUserResponse struct {
	Result bool         `json:"result"`
	User   UserResponse `json:"user"`
}

type GetUserByIDResponse struct {
	Result bool         `json:"result"`
	User   UserResponse `json:"user"`
}

type GetAllUsersResponse struct {
	Result bool           `json:"result"`
	Users  []UserResponse `json:"users"`
}
