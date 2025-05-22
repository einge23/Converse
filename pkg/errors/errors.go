package errors

import "net/http"

type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
    return e.Message
}

func NewConflictError(message string) *AppError {
    return &AppError{
        Code:    http.StatusConflict,
        Message: message,
    }
}

func NewBadRequestError(message string, details string) *AppError {
    return &AppError{
        Code:    http.StatusBadRequest,
        Message: message,
        Details: details,
    }
}