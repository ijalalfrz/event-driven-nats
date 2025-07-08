package dto

import (
	"net/http"

	"github.com/ijalalfrz/event-driven-nats/user-service/internal/pkg/exception"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/pkg/lang"
)

func NewInvalidRequestError(err error) exception.ApplicationError {
	return exception.ApplicationError{
		StatusCode: http.StatusBadRequest,
		Localizable: lang.Localizable{
			MessageID: "errors.invalid_request",
			MessageVars: map[string]interface{}{
				"message": err.Error(),
			},
		},
		Cause: err,
	}
}
