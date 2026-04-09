package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Moku3956/daily-routine/internal/adapter/middleware"
	"github.com/Moku3956/daily-routine/internal/domain"
	"github.com/Moku3956/daily-routine/internal/usecase"
)

type logUsecase interface {
	RecordLog(userId, habitId int, date time.Time, isCompleted bool) error
	GetLogs(userId int) ([]*domain.HabitLog, error)
}

type LogHttpHandler struct {
	usecase logUsecase
}

func NewLogHttpHandler(u logUsecase) *LogHttpHandler {
	return &LogHttpHandler{usecase: u}
}

type habitLogRequest struct {
	HabitId     int    `json:"habit_id"`
	Date        string `json:"date"`
	IsCompleted bool   `json:"is_completed"`
}

func (h *LogHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return
	}

	var req habitLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストの形式が不正です")
		return
	}
	if req.HabitId == 0 || req.Date == "" {
		writeError(w, http.StatusBadRequest, "habit_idとdateは必須です")
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "dateはYYYY-MM-DD形式で入力してください")
		return
	}

	if err := h.usecase.RecordLog(userId, req.HabitId, date, req.IsCompleted); err != nil {
		if errors.Is(err, usecase.ErrHabitNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, usecase.ErrHabitForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, usecase.ErrLogAlreadyExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "ログ記録に失敗しました")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *LogHttpHandler) List(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return
	}

	logs, err := h.usecase.GetLogs(userId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ログ取得に失敗しました")
		return
	}

	type logResponse struct {
		HabitId int    `json:"habit_id"`
		Date    string `json:"date"`
	}
	resp := make([]logResponse, 0, len(logs))
	for _, l := range logs {
		resp = append(resp, logResponse{
			HabitId: l.HabitId,
			Date:    l.Date.Format("2006-01-02"),
		})
	}
	writeJSON(w, http.StatusOK, resp)
}
