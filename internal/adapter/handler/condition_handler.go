package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Moku3956/daily-routine/internal/adapter/middleware"
	"github.com/Moku3956/daily-routine/internal/domain"
	"github.com/Moku3956/daily-routine/internal/usecase"
)

type conditionUsecase interface {
	InputCondition(userId, pCondition, mCondition int) (*domain.Condition, error)
	GetTodayCondition(userId int) (*domain.Condition, error)
}

type ConditionHttpHandler struct {
	usecase conditionUsecase
}

func NewConditionHttpHandler(u conditionUsecase) *ConditionHttpHandler {
	return &ConditionHttpHandler{usecase: u}
}

type conditionRequest struct {
	PCondition int `json:"p_condition"`
	MCondition int `json:"m_condition"`
}

func (h *ConditionHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return
	}

	var req conditionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストの形式が不正です")
		return
	}
	if req.PCondition < 1 || req.PCondition > 5 || req.MCondition < 1 || req.MCondition > 5 {
		writeError(w, http.StatusUnprocessableEntity, "コンディションは1〜5の値で入力してください")
		return
	}

	_, err := h.usecase.InputCondition(userId, req.PCondition, req.MCondition)
	if err != nil {
		if errors.Is(err, usecase.ErrConditionAlreadyExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "コンディション登録に失敗しました")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ConditionHttpHandler) Get(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return
	}

	c, err := h.usecase.GetTodayCondition(userId)
	if err != nil {
		if errors.Is(err, usecase.ErrConditionNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "コンディション取得に失敗しました")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"p_condition": c.PCondition,
		"m_condition": c.MCondition,
		"rank":        string(c.Rank),
		"date":        c.Date.Format("2006-01-02"),
	})
}
