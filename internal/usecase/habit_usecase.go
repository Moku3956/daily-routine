package usecase

import (
	"fmt"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type HabitUsecase struct {
	repo domain.HabitRepository
}

// handlerからのリクエストを受け取る
func (u *HabitUsecase) AddHabit(h *domain.Habit) error {
	if err := u.repo.Save(h); err != nil {
		return fmt.Errorf("習慣登録に失敗: %w", err)
	}
	return nil
}

func NewHabitUsecase(r domain.HabitRepository) *HabitUsecase {
	return &HabitUsecase{repo: r}
}
