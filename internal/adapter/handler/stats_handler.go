package handler

import (
	"net/http"

	"github.com/Moku3956/daily-routine/internal/adapter/middleware"
)

type statsUsecase interface {
	GetStreakDays(userId int) (int, error)
}

type StatsHttpHandler struct {
	usecase statsUsecase
}

func NewStatsHttpHandler(u statsUsecase) *StatsHttpHandler {
	return &StatsHttpHandler{usecase: u}
}

func (h *StatsHttpHandler) Get(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return
	}

	days, err := h.usecase.GetStreakDays(userId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "統計取得に失敗しました")
		return
	}

	writeJSON(w, http.StatusOK, days)
}
