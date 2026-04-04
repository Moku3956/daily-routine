package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type habitUsecase interface {
	AddHabit(habit *domain.Habit) error
}

type createHabitRequest struct {
	UserId    int    `json:"user_id"`
	HabitName string `json:"habit_name"`
	Category  string `json:"category"`
	Value     int    `json:"value"`
	Unit      string `json:"unit"`
	PCost     int    `json:"p_cost"`
	MCost     int    `json:"m_cost"`
	Must      bool   `json:"must"`
}

type HabitHttpHandler struct {
	usecase habitUsecase
}

func NewHabitHttpHandler(u habitUsecase) *HabitHttpHandler {
	return &HabitHttpHandler{usecase: u}
}

func (h *HabitHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req createHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	habit := &domain.Habit{
		UserId:    req.UserId,
		HabitName: req.HabitName,
		Category:  req.Category,
		Value:     req.Value,
		Unit:      req.Unit,
		PCost:     req.PCost,
		MCost:     req.MCost,
		Must:      req.Must,
	}
	if err := h.usecase.AddHabit(habit); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
