package apperr

import "net/http"

// Code identifies an application error as domain + number within that domain.
// Wire form is service.domain.error (e.g. "12.01.01").
type Code struct {
	Domain int
	Number int
}

const (
	domainGeneric = 99
	domainUser    = 1
)

var (
	// 99.xx — generic errors
	CodeGeneric        = Code{Domain: domainGeneric, Number: 99}
	CodeInternalServer = Code{Domain: domainGeneric, Number: 1}
	CodeInvalidArg     = Code{Domain: domainGeneric, Number: 2}
	CodeInvalidBinding = Code{Domain: domainGeneric, Number: 3}
	CodeUnauthorized   = Code{Domain: domainGeneric, Number: 4}

	// 01.xx — user domain errors
	CodeUserNotFound = Code{Domain: domainUser, Number: 1}
)

var errorCodeToStatusCode = map[Code]int{
	CodeGeneric:        http.StatusInternalServerError,
	CodeInternalServer: http.StatusInternalServerError,
	CodeInvalidArg:     http.StatusBadRequest,
	CodeInvalidBinding: http.StatusBadRequest,
	CodeUnauthorized:   http.StatusUnauthorized,

	CodeUserNotFound: http.StatusNotFound,
}

var errorCodeToKey = map[Code]string{
	CodeGeneric:        "generic",
	CodeInternalServer: "error.internal.server",
	CodeInvalidArg:     "invalid.argument",
	CodeInvalidBinding: "invalid.binding",
	CodeUnauthorized:   "error.unauthorized",

	CodeUserNotFound: "user.not_found",
}

func orderedErrorKeys() []string {
	keys := make([]string, 0, len(errorCodeToKey))
	for _, key := range errorCodeToKey {
		keys = append(keys, key)
	}
	return keys
}

// Sentinels below cover generic, service-wide errors only.
// Domain-specific errors should be constructed directly at call sites.
var (
	Generic         = &Error{ErrCode: CodeGeneric}
	InternalServer  = &Error{ErrCode: CodeInternalServer}
	InvalidArgument = &Error{ErrCode: CodeInvalidArg}
	InvalidBinding  = &Error{ErrCode: CodeInvalidBinding}
	Unauthorized    = &Error{ErrCode: CodeUnauthorized}
)
