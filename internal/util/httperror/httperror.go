package httperror

import (
	"context"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Info    string `json:"info"`
	Message string `json:"message"`
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError: %d - %s", h.Code, h.Info)
}

type Option func(*HTTPError)

func WithMessage(message string) Option {
	return func(h *HTTPError) {
		h.Message = message
	}
}
func NewBadRequest(ctx context.Context, opts ...Option) *HTTPError {
	httpError := &HTTPError{
		Code:    http.StatusBadRequest,
		Info:    "BAD_REQUEST",
		Message: "",
	}
	for _, opt := range opts {
		opt(httpError)
	}
	return httpError
}
func NewDataNotFound(ctx context.Context, opts ...Option) *HTTPError {
	httpError := &HTTPError{
		Code:    http.StatusNotFound,
		Info:    "DATA_NOT_FOUND",
		Message: "",
	}
	for _, opt := range opts {
		opt(httpError)
	}
	return httpError
}
func NewInternalServer(ctx context.Context, opts ...Option) *HTTPError {
	httpError := &HTTPError{
		Code:    http.StatusInternalServerError,
		Info:    "INTERNAL_SERVER",
		Message: "",
	}
	for _, opt := range opts {
		opt(httpError)
	}
	return httpError
}
