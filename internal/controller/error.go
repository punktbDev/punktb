package controller

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type (
	errorResponse struct {
		Error   bool   `json:"error"`
		Message string `json:"message,omitempty"`
	}
	ErrorResponse interface {
		BadRequest(w http.ResponseWriter, err error, code int)
	}
)

func NewError() ErrorResponse {
	return &errorResponse{}
}

func (e *errorResponse) BadRequest(w http.ResponseWriter, err error, code int) {
	out := &errorResponse{
		Error:   true,
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err = json.NewEncoder(w).Encode(out); err != nil {
		zap.L().Error("errorResponse marshaling", zap.Error(err))
	}
}
