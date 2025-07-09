//go:build unit

package dto

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserRequest_Bind(t *testing.T) {
	bindRequest := func(name string, req *CreateUserRequest, httpReq *http.Request, wantErr bool) func(t *testing.T) {
		return func(t *testing.T) {
			err := req.Bind(httpReq)
			if wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		}
	}

	t.Run("success", bindRequest(
		"success",
		&CreateUserRequest{Name: "John Doe"},
		&http.Request{},
		false,
	))

	t.Run("empty_name", bindRequest(
		"empty_name",
		&CreateUserRequest{Name: ""},
		&http.Request{},
		true,
	))
}

func TestGetAllUsersRequest_Bind(t *testing.T) {
	bindRequest := func(name string, req *GetAllUsersRequest, queryParams url.Values, wantErr bool, wantPageNum, wantPageSize int) func(t *testing.T) {
		return func(t *testing.T) {
			httpReq := &http.Request{URL: &url.URL{RawQuery: queryParams.Encode()}}
			err := req.Bind(httpReq)
			if wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, wantPageNum, req.PageNumber)
			assert.Equal(t, wantPageSize, req.PageSize)
		}
	}

	t.Run("success_with_defaults", bindRequest(
		"success_with_defaults",
		&GetAllUsersRequest{},
		url.Values{},
		false,
		1,  // default page number
		10, // default page size
	))

	t.Run("success_with_valid_params", bindRequest(
		"success_with_valid_params",
		&GetAllUsersRequest{},
		url.Values{
			"page_num":  []string{"2"},
			"page_size": []string{"20"},
		},
		false,
		2,
		20,
	))

	t.Run("invalid_page_number", bindRequest(
		"invalid_page_number",
		&GetAllUsersRequest{},
		url.Values{
			"page_num": []string{"invalid"},
		},
		true,
		0,
		0,
	))

	t.Run("invalid_page_size", bindRequest(
		"invalid_page_size",
		&GetAllUsersRequest{},
		url.Values{
			"page_size": []string{"invalid"},
		},
		true,
		0,
		0,
	))

	t.Run("zero_page_number", bindRequest(
		"zero_page_number",
		&GetAllUsersRequest{},
		url.Values{
			"page_num": []string{"0"},
		},
		true,
		0,
		0,
	))

	t.Run("zero_page_size", bindRequest(
		"zero_page_size",
		&GetAllUsersRequest{},
		url.Values{
			"page_size": []string{"0"},
		},
		true,
		0,
		0,
	))
}

func TestGetUserByIDRequest_Bind(t *testing.T) {
	bindRequest := func(name string, req *GetUserByIDRequest, urlParam string, wantErr bool, wantID int64) func(t *testing.T) {
		return func(t *testing.T) {
			// Setup chi router context with URL parameter
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", urlParam)

			httpReq := &http.Request{}
			httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))

			err := req.Bind(httpReq)
			if wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, wantID, req.ID)
		}
	}

	t.Run("success", bindRequest(
		"success",
		&GetUserByIDRequest{},
		"123",
		false,
		123,
	))

	t.Run("invalid_id", bindRequest(
		"invalid_id",
		&GetUserByIDRequest{},
		"invalid",
		true,
		0,
	))

	t.Run("empty_id", bindRequest(
		"empty_id",
		&GetUserByIDRequest{},
		"",
		true,
		0,
	))
}
