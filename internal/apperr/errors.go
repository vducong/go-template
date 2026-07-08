package apperr

import (
	"fmt"
	"net/http"
)

// Error represents a standardized API error.
// Treat as immutable. Use WithDetail() to create a variant with a custom detail.
type Error struct {
	ErrCode   Code
	Err       error
	ErrDetail string
}

func (e *Error) StatusCode() int {
	if status, ok := errorCodeToStatusCode[e.ErrCode]; ok {
		return status
	}
	return http.StatusInternalServerError
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Detail()
}

// Code returns the wire error code as service.domain.error (e.g. "12.01.01").
func (e *Error) Code() string {
	return fmt.Sprintf("%02d.%02d.%02d", codePrefix, e.ErrCode.Domain, e.ErrCode.Number)
}

func (e *Error) Key() string {
	key, ok := errorCodeToKey[e.ErrCode]
	if !ok {
		key = "generic"
	}
	return messagePrefix + "." + key
}

func (e *Error) Detail() string {
	if e.ErrDetail != "" {
		return e.ErrDetail
	}
	if detail, ok := messageForKey(e.Key()); ok {
		return detail
	}
	return "internal server error"
}

// WithDetail returns a NEW Error with the same code but a custom detail.
// The original sentinel is never mutated.
func (e *Error) WithDetail(detail string) *Error {
	newErr := *e
	newErr.ErrDetail = detail
	return &newErr
}
