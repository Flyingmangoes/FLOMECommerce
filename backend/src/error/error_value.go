package error_service

import (
	"net/http"
)

func ErrBadRequest(msg string) *AppError {
    return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func ErrConflict(msg string) *AppError {
    return &AppError{Code: http.StatusConflict, Message: msg}
}

func ErrInternal(msg string) *AppError {
    return &AppError{Code: http.StatusInternalServerError, Message: msg}
}

func ErrUnauthorized(msg string) *AppError {
    return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func ErrNotFound(msg string) *AppError {
    return &AppError{Code: http.StatusNotFound, Message: msg}
}

func ErrPreconditionRequired(msg string) *AppError {
    return &AppError{Code: http.StatusPreconditionRequired, Message: msg}
}