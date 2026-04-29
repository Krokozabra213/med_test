package domain

import "fmt"

type Code string

const (
	CodeNotFound   Code = "NOT_FOUND"
	CodeInternal   Code = "INTERNAL_ERROR"
	CodeBadRequest Code = "BAD_REQUEST"
)

const (
	MessageInternal = "internal error"
)

type AppError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return string(e.Code)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewNotFound(entity string) *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: fmt.Sprintf("%s not found", entity),
	}
}

func NewInternal(msg string, err error) *AppError {
	return &AppError{
		Code:    CodeInternal,
		Message: msg,
		Err:     err,
	}
}

func NewBadRequest(msg string, err error) *AppError {
	return &AppError{
		Code:    CodeBadRequest,
		Message: msg,
		Err:     err,
	}
}
