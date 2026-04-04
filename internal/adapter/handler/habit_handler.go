package handler

import (
	"net/http"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type HabitHttpHandler struct {
	handler *domain.HabitRepository
}

func HabitHandler(w http.ResponseWriter, r *http.Request) {

}
