package dto

import (
	"context"
	"net/http"
)

type RequestContext struct {
	Language string `mapstructure:"language"`
}

type contextKey string

// requestContextKey is the context.Context key to store the request context.
var requestContextKey = contextKey("request_context")

func RequestWithContext(req *http.Request) (*http.Request, error) {
	var reqContext RequestContext

	reqContext.Language = getLanguage(req)

	ctx := context.WithValue(req.Context(), requestContextKey, reqContext)

	return req.WithContext(ctx), nil
}

func RequestFromContext(ctx context.Context) (RequestContext, bool) {
	reqContext, ok := ctx.Value(requestContextKey).(RequestContext)

	return reqContext, ok
}

func getLanguage(req *http.Request) string {
	return req.Header.Get("Accept-Language")
}
