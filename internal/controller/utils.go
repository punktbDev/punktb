package controller

import (
	"encoding/json"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
	"go.uber.org/zap"
	"net/http"
)

func SendResponse(status int, w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		zap.L().Error("Encode", zap.Error(err))
	}
}

func GetManager(r *http.Request) *dto.Manager {
	mng := r.Context().Value("manager")
	return mng.(*dto.Manager)
}
