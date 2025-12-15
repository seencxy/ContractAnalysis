package errors

import "fmt"

// ErrorCode represents API error codes
type ErrorCode int

const (
	// Client errors (4xx)
	ErrBadRequest       ErrorCode = 400
	ErrUnauthorized     ErrorCode = 401
	ErrForbidden        ErrorCode = 403
	ErrNotFound         ErrorCode = 404
	ErrValidationFailed ErrorCode = 422

	// Server errors (5xx)
	ErrInternalServer ErrorCode = 500
	ErrDatabase       ErrorCode = 501
	ErrService        ErrorCode = 502
)

// APIError represents an API error
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Type    string    `json:"type"`
	Details []string  `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Type, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(code ErrorCode, message string, errorType string, details ...string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Type:    errorType,
		Details: details,
	}
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string, details ...string) *APIError {
	return NewAPIError(ErrBadRequest, message, "BadRequest", details...)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *APIError {
	return NewAPIError(ErrNotFound, message, "NotFound")
}

// NewValidationError creates a validation error
func NewValidationError(message string, details ...string) *APIError {
	return NewAPIError(ErrValidationFailed, message, "ValidationError", details...)
}

// NewInternalServerError creates an internal server error
func NewInternalServerError(message string) *APIError {
	return NewAPIError(ErrInternalServer, message, "InternalServerError")
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string) *APIError {
	return NewAPIError(ErrDatabase, message, "DatabaseError")
}
