package apperror

import (
	"fmt"
	"net/http"
)

type Error struct {
	errorType ErrorType
	message   string
	error     string
	thrownAt  string
}

// Error returns the error
func (e *Error) Error() string {
	return fmt.Sprintf("[type]: %s - [message]: %s - [error]: %s - [thrownAt]: %s", e.errorType.String(), e.message, e.error, e.thrownAt)
}

// NewNotFound creates a error configured to be a NotFound response
func NewNotFound(message string, error interface{}, thrownAt string) error {
	return NewError(
		NotFound,
		message,
		error,
		thrownAt,
	)
}

func NewUnauthorized(message string, error interface{}, thrownAt string) error {
	return NewError(
		Unauthorized,
		message,
		error,
		thrownAt,
	)
}

// NewUnexpectedDatabaseError corresponds to an unexpected error thrown by the database wrapper
func NewUnexpectedDatabaseError(message string, error interface{}, thrownAt string) error {
	return NewError(
		UnexpectedDatabaseError,
		message,
		error,
		thrownAt,
	)
}

// NewForbidden corresponds to a forbidden error
func NewForbidden(message string, error interface{}, thrownAt string) error {
	return NewError(
		Forbidden,
		message,
		error,
		thrownAt,
	)
}

// NewThrottling corresponds to an unexpected error thrown by the database wrapper
func NewThrottling(message string, error interface{}, thrownAt string) error {
	return NewError(
		TooManyRequests,
		message,
		error,
		thrownAt,
	)
}

// NewGone corresponds to an unexpected error thrown by the database wrapper or the message client
func NewGone(message string, error interface{}, thrownAt string) error {
	return NewError(
		Gone,
		message,
		error,
		thrownAt,
	)
}

// NewInternalServer corresponds to an unexpected error thrown by the database wrapper
func NewInternalServer(message string, error interface{}, thrownAt string) error {
	return NewError(
		InternalServer,
		message,
		error,
		thrownAt,
	)
}

// NewInvalidData corresponds to invalid data (probably from the incoming HTTP request's parameters)
func NewInvalidData(message string, error interface{}, thrownAt string) error {
	return NewError(
		InvalidData,
		message,
		error,
		thrownAt,
	)
}

// NewPayloadTooLarge corresponds to a payload too large to be processed
func NewPayloadTooLarge(message string, error interface{}, thrownAt string) error {
	return NewError(
		PayloadTooLarge,
		message,
		error,
		thrownAt,
	)
}

// NewError creates a new error with the given type, message, error and thrownAt
func NewError(errorType ErrorType, message string, error interface{}, thrownAt string) error {
	return &Error{
		errorType: errorType,
		message:   message,
		error:     errorString(error),
		thrownAt:  thrownAt,
	}
}

func errorString(err interface{}) string {
	switch newErr := err.(type) {
	case error:
		return newErr.Error()
	case string:
		return newErr
	default:
		return "unknown error variable"
	}
}

// StatusCodeAndMessage extract and return the status code and message of the given error
func StatusCodeAndMessage(e error) (int, string) {
	return StatusCode(e), Message(e)
}

// Message retrieves the message of the given error
func Message(e error) string {
	if err, ok := e.(*Error); ok {
		return fmt.Sprintf("%s: [error]: %s - [thrownAt]: %s", err.message, err.error, err.thrownAt)
	}

	return e.Error()
}

// StatusCode retrieves the http status code of the given error
func StatusCode(e error) int {
	defaultStatus := http.StatusInternalServerError

	errorTypeToStatus := map[ErrorType]int{
		NotFound:                http.StatusNotFound,
		UnexpectedDatabaseError: defaultStatus,
		Gone:                    http.StatusGone,
		InternalServer:          defaultStatus,
		Forbidden:               http.StatusForbidden,
		InvalidData:             http.StatusBadRequest,
		TooManyRequests:         http.StatusTooManyRequests,
		PayloadTooLarge:         http.StatusRequestEntityTooLarge,
	}

	if err, ok := e.(*Error); ok {
		if status, ok := errorTypeToStatus[err.errorType]; ok {
			return status
		}
	}
	return defaultStatus
}

// Is make assertion of error type of the given error
func Is(e error, t ErrorType) bool {
	if err, ok := e.(*Error); ok {
		return err.errorType == t
	}
	return false
}
