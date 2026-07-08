package httprespwrit

import (
	"net/http"

	"github.com/go-chi/render"
)

type JSONResponse struct {
	StatusCode int  `json:"-"`
	Success    bool `json:"success"`
	Data       any  `json:"data,omitempty"`
}

type ErrorResponse struct {
	JSONResponse
	Err   error        `json:"-"`
	Error *ErrorDetail `json:"error,omitempty"`
}

type AppErr interface {
	StatusCode() int
	Error() string
	Code() string
	Key() string
	Detail() string
}

type ErrorDetail struct {
	Code   string  `json:"code"`   // format: service.domain.error (e.g. "12.01.01")
	Key    *string `json:"key"`    // key like "{service_code}.internal" - can be used for i18n message key
	Detail *string `json:"detail"` // human readable message for debugging purpose
}

// Render sets the HTTP status code; chi render marshals the struct to JSON after this returns.
func (res *JSONResponse) Render(w http.ResponseWriter, req *http.Request) error {
	if res.StatusCode == 0 {
		res.StatusCode = http.StatusOK
	}
	res.Success = res.StatusCode >= http.StatusContinue && res.StatusCode <= http.StatusIMUsed
	render.Status(req, res.StatusCode)
	return nil
}

// Render populates Error from Err, sets the HTTP status code, and lets chi marshal the struct.
func (res *ErrorResponse) Render(w http.ResponseWriter, req *http.Request) error {
	if res.StatusCode == 0 {
		res.StatusCode = http.StatusInternalServerError
	}

	res.Success = false

	if res.Err != nil {
		appErr, ok := res.Err.(AppErr)
		if ok {
			res.StatusCode = appErr.StatusCode()
			key := appErr.Key()
			detail := appErr.Detail()
			res.Error = &ErrorDetail{
				Code:   appErr.Code(),
				Key:    &key,
				Detail: &detail,
			}
		} else {
			detail := res.Err.Error()
			res.Error = &ErrorDetail{
				Detail: &detail,
			}
		}
	}

	render.Status(req, res.StatusCode)
	return nil
}
