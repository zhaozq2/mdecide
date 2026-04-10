package errorx

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidParams  = errors.New("invalid params")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInternalServer = errors.New("internal server error")
)

type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewCodeError(code int, message string) *CodeError {
	return &CodeError{Code: code, Message: message}
}

func Wrap(err error, code int, message string) *CodeError {
	return &CodeError{Code: code, Message: message + ": " + err.Error()}
}
