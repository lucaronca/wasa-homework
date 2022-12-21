package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
)

var (
	// ErrTypeAssertionError is thrown when type an interface does not match the asserted type
	ErrTypeAssertionError = errors.New("unable to assert type")
)

// ParsingError indicates that an error has occurred when parsing request parameters
type ParsingError struct {
	Err error
}

func (e *ParsingError) Unwrap() error {
	return e.Err
}

func (e *ParsingError) Error() string {
	return e.Err.Error()
}

// RequiredError indicates that an error has occurred when parsing required request parameters
type RequiredError struct {
	Field string
}

func (e *RequiredError) Error() string {
	return fmt.Sprintf("required field '%s' is zero value.", e.Field)
}

// NotFoundError indicates that an error has occurred when trying to access a not existing entity
type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Entity)
}

// Forbidden indicates that an error has occurred when trying to access a forbidden resource
type ForbiddenError struct {
	Err error
}

func (e *ForbiddenError) Unwrap() error {
	return e.Err
}

func (e *ForbiddenError) Error() string {
	return e.Err.Error()
}

// ErrorHandler defines the required method for handling error.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error, ctx reqcontext.RequestContext)

// errorHandler defines the default logic on how to handle errors from the controller. Any errors from parsing
// request params will return a StatusBadRequest
func errorHandler(w http.ResponseWriter, r *http.Request, err error, ctx reqcontext.RequestContext) {
	var pe *ParsingError
	var re *RequiredError
	var nfe *NotFoundError
	var fe *ForbiddenError

	switch {
	case
		// Handle parsing errors
		errors.As(err, &pe),
		// Handle missing required errors
		errors.As(err, &re):
		encodeTextResponse(err.Error(), http.StatusBadRequest, w, ctx)
	// Handle missing entity errors
	case errors.As(err, &nfe):
		encodeTextResponse(err.Error(), http.StatusNotFound, w, ctx)
	// Handle forbidden entity errors
	case errors.As(err, &fe):
		encodeTextResponse(err.Error(), http.StatusForbidden, w, ctx)
	// Handle all other errors
	default:
		ctx.Logger.WithError(err).Error("Internal server error")
		encodeTextResponse("Internal server error", http.StatusInternalServerError, w, ctx)
	}
}
