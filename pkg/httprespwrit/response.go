package httprespwrit

import "net/http"

type Response[T any] struct {
	Success bool         `json:"success"`
	Error   *ErrorDetail `json:"error,omitempty"`
	Data    T            `json:"data,omitempty"`
}

type ErrorDetail struct {
	Code    string  `json:"code"`    // format: xxyyyy (xx - service code, yyyy - error code)
	Message *string `json:"message"` // code like "{service_code}.internal" - can be used for i18n message key
	Detail  *string `json:"detail"`  // human readable message for debugging purpose
}

func (res *Response[T]) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func ResponseJSON[Payload any](statusCode int, payload Payload) *Response[Payload] {
	return &Response[Payload]{
		Success: true,
		Data:    payload,
		Error:   nil,
	}
}

func ResponseError(err error) *Response[any] {
	return nil
}
