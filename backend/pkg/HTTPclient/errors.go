package httpclient

import "fmt"

type APIError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("API error: status=%d, message=%s, error=%v",
			e.StatusCode, e.Message, e.Err)
	}
	return fmt.Sprintf("API error: status=%d, message=%s",
		e.StatusCode, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 404
	}
	return false
}

func IsInternalServerError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode >= 500
	}
	return false
}
