package usecase

import (
	"fmt"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type HabitUsecase struct {
	repo domain.HabitRepository
}

func (u *HabitUsecase) AddHabit(h *domain.Habit) error {
	err := u.repo.Save(h)
	if err != nil {
		return fmt.Errorf("習慣登録に失敗: %w", err)
	}
	return nil
}
