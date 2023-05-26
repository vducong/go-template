package failure

import (
	"fmt"
	"net/http"
)

type ErrCode int

type AppErr struct {
	Code        ErrCode
	OriginalErr error
}

func (e *AppErr) Error() string {
	if val, ok := errMsgMap[e.Code]; ok {
		return val
	}
	return fmt.Sprintf("Lỗi hệ thống: %d\nVui lòng thử lại sau", e.Code)
}

func (e *AppErr) HTTPCode() int {
	if val, ok := errCodeMap[e.Code]; ok {
		return val
	}
	return http.StatusBadRequest
}

var errMsgMap = map[ErrCode]string{
	ErrReusableCodeNotFound: "Mã quà tặng không tồn tại",
}

var errCodeMap = map[ErrCode]int{
	ErrReusableCodeFailed: http.StatusInternalServerError,
}

const (
	ErrReusableCodeGetByCodeBinding ErrCode = 991001
	ErrReusableCodeNotFound         ErrCode = 991002
	ErrReusableCodeFailed           ErrCode = 991003
)
