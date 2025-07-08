package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/pkg/exception"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/pkg/lang"
)

// ErrInvalidType invalid type of request.
var ErrInvalidType = exception.ApplicationError{
	Localizable: lang.Localizable{
		Message: "invalid type",
	},
	StatusCode: exception.CodeBadRequest,
}

type PublicListing struct {
	Create endpoint.Endpoint
	GetAll endpoint.Endpoint
}

type PublicUser struct {
	Create endpoint.Endpoint
}

type Endpoint struct {
	PublicListing
	PublicUser
}
