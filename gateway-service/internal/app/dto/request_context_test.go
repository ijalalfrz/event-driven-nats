//go:build unit

package dto

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestContext(t *testing.T) {
	var (
		language = "en"
	)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/foo", nil)
	assert.NoError(t, err)

	req.Header.Add("Accept-Language", language)

	out, err := RequestWithContext(req)
	assert.NoError(t, err)

	reqContext, ok := RequestFromContext(out.Context())
	assert.True(t, ok)

	assert.Equal(t, language, reqContext.Language)
}
