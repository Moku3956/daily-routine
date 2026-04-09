package usecase

import (
	"fmt"
	"time"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type StatsUsecase struct {
	logRepo domain.LogRepository
}

func NewStatsUsecase(logRepo domain.LogRepository) *StatsUsecase {
	return &StatsUsecase{logRepo: logRepo}
}

func (u *StatsUsecase) GetStreakDays(userId int) (int, error) {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	count, err := u.logRepo.CountConsecutiveDays(userId, today)
	if err != nil {
		return 0, fmt.Errorf("連続日数取得に失敗: %w", err)
	}
	return count, nil
}
