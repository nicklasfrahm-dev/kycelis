package response

import (
	"errors"
	"net/http"
)

// Status represents the status of a response.
type Status struct {
	// Code is the status code.
	Code int `json:"code"`
	// Title is the status title.
	Title string `json:"title"`
	// Message provides additional information about the status.
	Message string `json:"message"`
}

// MapErrorToStatus maps an error to a status.
//
//nolint:gochecknoglobals // This is a map of errors to status codes.
var MapErrorToStatusCode = map[error]int{
	ErrServiceHealthy:  http.StatusOK,
	ErrUnexpectedError: http.StatusInternalServerError,
	ErrUnknownEndpoint: http.StatusNotFound,
}

var (
	// ErrServiceHealthy is returned when the service is healthy.
	ErrServiceHealthy = errors.New("ServiceHealthy")
	// ErrUnexpectedError is returned when an unexpected error occurs.
	ErrUnexpectedError = errors.New("UnexpectedError")
	// ErrUnknownEndpoint is returned when an unknown endpoint is requested.
	ErrUnknownEndpoint = errors.New("UnknownEndpoint")
)

// NewStatus creates a new Status instance.
func NewStatus(code int, message string) *Status {
	return &Status{
		Code:    code,
		Title:   http.StatusText(code),
		Message: message,
	}
}

// NewStatusFromError creates a new Status instance from an error.
func NewStatusFromError(err error) *Status {
	code, ok := MapErrorToStatusCode[err]
	if !ok {
		return NewStatusFromError(ErrUnexpectedError)
	}

	return NewStatus(code, err.Error())
}
