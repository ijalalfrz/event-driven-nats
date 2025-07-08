package dto

import (
	"fmt"
	"net/http"
	"strconv"
)

type ListingResponse struct {
	ID          int64        `json:"id"`
	ListingType string       `json:"listing_type"`
	Price       int64        `json:"price"`
	CreatedAt   int64        `json:"created_at"`
	UpdatedAt   int64        `json:"updated_at"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type GetAllListingsRequest struct {
	PageNum  int    `json:"page_num" default:"1"`
	PageSize int    `json:"page_size" default:"10"`
	UserID   *int64 `json:"user_id"`
}

func (r *GetAllListingsRequest) Bind(req *http.Request) error {
	pageNum := req.URL.Query().Get("page_num")
	pageSize := req.URL.Query().Get("page_size")
	userID := req.URL.Query().Get("user_id")

	if pageNum == "" {
		pageNum = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		return NewInvalidRequestError(fmt.Errorf("invalid page_num: %w", err))
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return NewInvalidRequestError(fmt.Errorf("invalid page_size: %w", err))
	}

	r.PageNum = pageNumInt
	r.PageSize = pageSizeInt

	var userIDInt int64
	if userID == "" {
		r.UserID = nil
		return nil
	}

	userIDInt, err = strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return NewInvalidRequestError(fmt.Errorf("invalid user_id: %w", err))
	}

	r.UserID = &userIDInt

	return nil
}

type GetAllListingsResponse struct {
	Result   bool              `json:"result"`
	Listings []ListingResponse `json:"listings"`
}

type ListingCreated struct {
	ID          int64  `json:"id"`
	ListingType string `json:"listing_type"`
	Price       int64  `json:"price"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	UserID      int64  `json:"user_id"`
}
