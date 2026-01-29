package evatr

import "fmt"

// ErrorResponse represents the JSON error response from the API.
type ErrorResponse struct {
	// eVATR status code (e.g., "evatr-0002")
	Status string `json:"status"`

	// Error message in German
	Message string `json:"meldung"`
}

// Error represents an eVATR API error.
type Error struct {
	// HTTP status code
	StatusCode int

	// eVATR status code (e.g., "evatr-0002")
	Status string

	// Human-readable error message
	Message string
}

func (e *Error) Error() string {
	if e.Status != "" {
		return fmt.Sprintf("evatr: %s (HTTP %d): %s", e.Status, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("evatr: HTTP %d: %s", e.StatusCode, e.Message)
}

// IsEvatrErr returns whether the error is an eVATR error.
func IsEvatrErr(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*Error)
	return ok
}

// Common error constructors based on status codes from the API

// NewBadRequestError returns a 400 Bad Request error.
func NewBadRequestError(status, message string) *Error {
	return &Error{
		StatusCode: 400,
		Status:     status,
		Message:    message,
	}
}

// NewForbiddenError returns a 403 Forbidden error.
func NewForbiddenError(status, message string) *Error {
	return &Error{
		StatusCode: 403,
		Status:     status,
		Message:    message,
	}
}

// NewNotFoundError returns a 404 Not Found error.
func NewNotFoundError(status, message string) *Error {
	return &Error{
		StatusCode: 404,
		Status:     status,
		Message:    message,
	}
}

// NewInternalServerError returns a 500 Internal Server Error.
func NewInternalServerError(status, message string) *Error {
	return &Error{
		StatusCode: 500,
		Status:     status,
		Message:    message,
	}
}

// NewServiceUnavailableError returns a 503 Service Unavailable error.
func NewServiceUnavailableError(status, message string) *Error {
	return &Error{
		StatusCode: 503,
		Status:     status,
		Message:    message,
	}
}

// Known error status codes

// Bad Request (400) errors
const (
	StatusMissingRequiredField        = "evatr-0002" // At least one required field is missing
	StatusInvalidRequestingVATID      = "evatr-0004" // Requesting DE VAT ID is syntactically incorrect
	StatusInvalidRequestedVATID       = "evatr-0005" // Requested VAT ID is syntactically incorrect
	StatusMaxQualifiedRequestsReached = "evatr-0008" // Maximum qualified requests for session reached
	StatusInvalidVATIDFormat          = "evatr-0012" // Requested VAT ID doesn't match format
	StatusInvalidCountryCode          = "evatr-2003" // Country code is not valid
)

// Forbidden (403) errors
const (
	StatusNotAuthorizedDE = "evatr-0006" // Requesting DE VAT ID not authorized to query DE VAT IDs
	StatusInvalidCall     = "evatr-0007" // Invalid call
)

// Not Found (404) errors
const (
	StatusVATIDNotAssigned        = "evatr-2001" // VAT ID not assigned at request time
	StatusRequestingVATIDNotValid = "evatr-2005" // Requesting DE VAT ID not valid at request time
)

// Success (200) with special meanings
const (
	StatusValid                = "evatr-0000" // VAT ID is valid at request time
	StatusNotYetValid          = "evatr-2002" // VAT ID not yet valid, see gueltigAb
	StatusNoLongerValid        = "evatr-2006" // VAT ID no longer valid, see gueltigAb and gueltigBis
	StatusValidWithSpecialCase = "evatr-2008" // VAT ID valid but special case, contact BZSt
)

// Internal Server Error (500) errors
const (
	StatusProcessingError1 = "evatr-2004" // Processing temporarily not possible
	StatusProcessingError2 = "evatr-2011" // Processing temporarily not possible
	StatusProcessingError3 = "evatr-3011" // Processing temporarily not possible
)

// Service Unavailable (503) errors
const (
	StatusServiceUnavailable1 = "evatr-0011" // Service temporarily unavailable
	StatusServiceUnavailable2 = "evatr-1001" // Service temporarily unavailable
	StatusServiceUnavailable3 = "evatr-1002" // Service temporarily unavailable
	StatusServiceUnavailable4 = "evatr-1003" // Service temporarily unavailable
	StatusServiceUnavailable5 = "evatr-1004" // Service temporarily unavailable
)
