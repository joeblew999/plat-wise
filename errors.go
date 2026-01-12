package wise

import "fmt"

// APIError represents an error returned by the Wise API.
type APIError struct {
	StatusCode int              `json:"-"`
	Type       string           `json:"type,omitempty"`
	Message    string           `json:"message,omitempty"`
	Errors     []ValidationError `json:"errors,omitempty"`
}

// ValidationError represents a validation error from the API.
type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Path    string `json:"path,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("wise: API error (status %d): %s - %v", e.StatusCode, e.Message, e.Errors)
	}
	return fmt.Sprintf("wise: API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found error.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsForbidden returns true if the error is a 403 Forbidden error.
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == 403
}

// IsRateLimited returns true if the error is a 429 Too Many Requests error.
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == 429
}
