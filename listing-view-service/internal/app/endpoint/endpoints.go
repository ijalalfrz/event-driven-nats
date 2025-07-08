package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/pkg/exception"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/pkg/lang"
)

// ErrInvalidType invalid type of request.
var ErrInvalidType = exception.ApplicationError{
	Localizable: lang.Localizable{
		Message: "invalid type",
	},
	StatusCode: exception.CodeBadRequest,
}

type Listing struct {
	GetAll    endpoint.Endpoint
	OnCreated endpoint.Endpoint
}

type User struct {
	OnCreated endpoint.Endpoint
}

type Endpoint struct {
	Listing
	User
}
