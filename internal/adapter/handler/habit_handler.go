package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Moku3956/daily-routine/internal/adapter/middleware"
	"github.com/Moku3956/daily-routine/internal/domain"
	"github.com/Moku3956/daily-routine/internal/usecase"
)

type habitUsecase interface {
	AddHabit(habit *domain.Habit) error
	GetHabits(userId int) ([]*domain.Habit, error)
	UpdateHabit(userId int, habit *domain.Habit) (*domain.Habit, error)
	DeleteHabit(userId, habitId int) error
	GetOrderedHabits(userId int) ([]*usecase.OrderedHabit, error)
}

type HabitHttpHandler struct {
	usecase habitUsecase
}

func NewHabitHttpHandler(u habitUsecase) *HabitHttpHandler {
	return &HabitHttpHandler{usecase: u}
}

type habitRequest struct {
	HabitName   string `json:"habit_name"`
	Category    string `json:"category"`
	TargetValue int    `json:"target_value"`
	Unit        string `json:"unit"`
	PCost       int    `json:"p_cost"`
	MCost       int    `json:"m_cost"`
	Must        bool   `json:"must"`
}

type habitResponse struct {
	HabitId     int    `json:"habit_id"`
	HabitName   string `json:"habit_name"`
	Category    string `json:"category"`
	TargetValue int    `json:"target_value"`
	Unit        string `json:"unit"`
	PCost       int    `json:"p_cost"`
	MCost       int    `json:"m_cost"`
	Must        bool   `json:"must"`
}

func toHabitResponse(h *domain.Habit) habitResponse {
	return habitResponse{
		HabitId:     h.HabitId,
		HabitName:   h.HabitName,
		Category:    h.Category,
		TargetValue: h.TargetValue,
		Unit:        h.Unit,
		PCost:       h.PCost,
		MCost:       h.MCost,
		Must:        h.Must,
	}
}

func (h *HabitHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return
	}

	var req habitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストの形式が不正です")
		return
	}
	if req.HabitName == "" {
		writeError(w, http.StatusBadRequest, "習慣名は必須です")
		return
	}
	if req.PCost < 1 || req.PCost > 10 || req.MCost < 1 || req.MCost > 10 {
		writeError(w, http.StatusUnprocessableEntity, "コストは1〜10の値で入力してください")
		return
	}

	habit := &domain.Habit{
		UserId:      userId,
		HabitName:   req.HabitName,
		Category:    req.Category,
		TargetValue: req.TargetValue,
		Unit:        req.Unit,
		PCost:       req.PCost,
		MCost:       req.MCost,
		Must:        req.Must,
	}
	if err := h.usecase.AddHabit(habit); err != nil {
		writeError(w, http.StatusInternalServerError, "習慣登録に失敗しました")
		return
	}
	writeJSON(w, http.StatusCreated, toHabitResponse(habit))
}

func (h *HabitHttpHandler) List(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return
	}

	habits, err := h.usecase.GetHabits(userId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "習慣一覧取得に失敗しました")
		return
	}

	resp := make([]habitResponse, 0, len(habits))
	for _, habit := range habits {
		resp = append(resp, toHabitResponse(habit))
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *HabitHttpHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return
	}

	habitId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "習慣IDが不正です")
		return
	}

	var req habitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストの形式が不正です")
		return
	}
	if req.PCost < 1 || req.PCost > 10 || req.MCost < 1 || req.MCost > 10 {
		writeError(w, http.StatusUnprocessableEntity, "コストは1〜10の値で入力してください")
		return
	}

	habit := &domain.Habit{
		HabitId:     habitId,
		HabitName:   req.HabitName,
		Category:    req.Category,
		TargetValue: req.TargetValue,
		Unit:        req.Unit,
		PCost:       req.PCost,
		MCost:       req.MCost,
		Must:        req.Must,
	}
	updated, err := h.usecase.UpdateHabit(userId, habit)
	if err != nil {
		if errors.Is(err, usecase.ErrHabitNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, usecase.ErrHabitForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "習慣更新に失敗しました")
		return
	}
	writeJSON(w, http.StatusOK, toHabitResponse(updated))
}

func (h *HabitHttpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return
	}

	habitId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "習慣IDが不正です")
		return
	}

	if err := h.usecase.DeleteHabit(userId, habitId); err != nil {
		if errors.Is(err, usecase.ErrHabitNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, usecase.ErrHabitForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "習慣削除に失敗しました")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HabitHttpHandler) GetOrdered(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return
	}

	ordered, err := h.usecase.GetOrderedHabits(userId)
	if err != nil {
		if errors.Is(err, usecase.ErrConditionNotFound) {
			writeError(w, http.StatusBadRequest, "本日のコンディションが未入力です")
			return
		}
		writeError(w, http.StatusInternalServerError, "習慣リスト取得に失敗しました")
		return
	}

	type orderedResponse struct {
		HabitId     int    `json:"habit_id"`
		HabitName   string `json:"habit_name"`
		Category    string `json:"category"`
		TargetValue int    `json:"target_value"`
		Unit        string `json:"unit"`
		PCost       int    `json:"p_cost"`
		MCost       int    `json:"m_cost"`
		Must        bool   `json:"must"`
	}
	resp := make([]orderedResponse, 0, len(ordered))
	for _, o := range ordered {
		resp = append(resp, orderedResponse{
			HabitId:     o.HabitId,
			HabitName:   o.HabitName,
			Category:    o.Category,
			TargetValue: o.TargetValue,
			Unit:        o.Unit,
			PCost:       o.PCost,
			MCost:       o.MCost,
			Must:        o.Must,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}
