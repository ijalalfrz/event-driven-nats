package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/pkg/exception"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/pkg/lang"
)

// ErrInvalidType invalid type of request.
var ErrInvalidType = exception.ApplicationError{
	Localizable: lang.Localizable{
		Message: "invalid type",
	},
	StatusCode: exception.CodeBadRequest,
}

type User struct {
	CreateUser  endpoint.Endpoint
	GetAllUsers endpoint.Endpoint
	GetUserByID endpoint.Endpoint
}

type Endpoint struct {
	User
}
