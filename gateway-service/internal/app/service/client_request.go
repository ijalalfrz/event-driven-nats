package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/pkg/exception"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/pkg/lang"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ClientOption func(*HTTPClient)

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *HTTPClient) {
		c.client.Timeout = timeout
	}
}

func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *HTTPClient) {
		c.maxRetries = maxRetries
	}
}

type HTTPClient struct {
	client     *http.Client
	url        string
	maxRetries int
}

func (hc *HTTPClient) doRequestWithResponse(
	ctx context.Context, method, path string, headerFunc func(req *http.Request), req interface{},
	errorResponseFunc func(resp *http.Response) error,
) (*http.Response, error) {
	var reqBody *bytes.Buffer

	// Check if request is a string (form-urlencoded data)
	if formData, ok := req.(string); ok {
		reqBody = bytes.NewBufferString(formData)
	} else {
		// Handle JSON encoding for other types
		reqBody = &bytes.Buffer{}
		err := json.NewEncoder(reqBody).Encode(req)
		if err != nil {
			return nil, fmt.Errorf("encode request: %w", err)
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, hc.url+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create HTTP request: %w", err)
	}

	headerFunc(httpReq)

	backOffTime := 100
	maxRetries := hc.maxRetries
	backoff := time.Duration(backOffTime) * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := hc.client.Do(httpReq)
		if err != nil {
			// Check if context was canceled or deadline exceeded
			if ctx.Err() != nil {
				return nil, fmt.Errorf("context error: %w", ctx.Err())
			}

			// Wait with exponential backoff before retrying
			if attempt < maxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2

				continue
			}

			return nil, fmt.Errorf("do request: %w", err)
		}

		if resp.StatusCode >= http.StatusBadRequest {
			return nil, errorResponseFunc(resp)
		}

		return resp, nil
	}

	return nil, fmt.Errorf("max retries exceeded")
}

func defaultErrorResponseFunc(resp *http.Response) error { //nolint:unused
	var errorResp ErrorResponse

	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		return fmt.Errorf("decode error response: %w", err)
	}

	// internal server errors should not be forwarded to the client
	if resp.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("returned status code: %d", resp.StatusCode)
	}

	host := strings.Split(resp.Request.URL.Host, ":")[0]

	return exception.ApplicationError{
		StatusCode: resp.StatusCode,
		Localizable: lang.Localizable{
			MessageID: "errors.bad_request_from_service",
			MessageVars: map[string]interface{}{
				"service": host,
				"error":   errorResp.Error,
			},
		},
		Cause: fmt.Errorf("error response: %s", errorResp.Error),
	}
}
