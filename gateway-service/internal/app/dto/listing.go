package dto

import (
	"fmt"
	"net/http"
	"strconv"
)

type CreateListingRequest struct {
	UserID      int64  `json:"user_id" validate:"required"`
	ListingType string `json:"listing_type" validate:"required,oneof=rent sale"`
	Price       int64  `json:"price" validate:"required,min=1"`
}

func (r *CreateListingRequest) Bind(req *http.Request) error {
	if err := validate.Struct(r); err != nil {
		return NewInvalidRequestError(fmt.Errorf("invalid request: %w", err))
	}

	return nil
}

type CreateListingResponse struct {
	ListingResponse `json:"listing"`
}

type ListingResponse struct {
	ID          int64         `json:"id"`
	UserID      int64         `json:"user_id,omitempty"`
	ListingType string        `json:"listing_type"`
	Price       int64         `json:"price"`
	CreatedAt   int64         `json:"created_at"`
	UpdatedAt   int64         `json:"updated_at"`
	User        *UserResponse `json:"user,omitempty"`
}

type GetAllListingsRequest struct {
	PageNumber int    `json:"page_number" validate:"required,min=1"`
	PageSize   int    `json:"page_size" validate:"required,min=1"`
	UserID     *int64 `json:"user_id" validate:"required"`
}

func (r *GetAllListingsRequest) Bind(req *http.Request) error {
	var err error

	pageNumbersStr := req.URL.Query().Get("page_num")
	pageSizesStr := req.URL.Query().Get("page_size")
	userIDStr := req.URL.Query().Get("user_id")

	if pageNumbersStr == "" {
		pageNumbersStr = "1"
	}
	if pageSizesStr == "" {
		pageSizesStr = "10"
	}

	pageNumber, err := strconv.Atoi(pageNumbersStr)
	if err != nil {
		return NewInvalidRequestError(fmt.Errorf("invalid page number: %w", err))
	}

	pageSize, err := strconv.Atoi(pageSizesStr)
	if err != nil {
		return NewInvalidRequestError(fmt.Errorf("invalid page size: %w", err))
	}

	r.PageNumber = pageNumber
	r.PageSize = pageSize

	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	if userID != 0 {
		r.UserID = &userID
	}

	return nil
}

type GetAllListingsResponse struct {
	Result   bool              `json:"result"`
	Listings []ListingResponse `json:"listings"`
}
