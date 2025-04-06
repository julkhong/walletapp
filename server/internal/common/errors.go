package common

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    int    `json:"code"`    // internal error code
	Message string `json:"message"` // human-readable message
}

// Internal error codes
const (
	ErrInvalidRequest      = 1000
	ErrWalletNotFound      = 1001
	ErrInsufficientBalance = 1002
	ErrDatabase            = 1003
	ErrUnknown             = 1099
)

// WriteError responds with internal error format
func WriteError(w http.ResponseWriter, httpStatus int, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}
