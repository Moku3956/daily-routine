package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/Moku3956/daily-routine/internal/domain"
)

var ErrLogAlreadyExists = errors.New("この習慣はすでに本日チェック済みです")

type LogUsecase struct {
	logRepo   domain.LogRepository
	habitRepo domain.HabitRepository
}

func NewLogUsecase(logRepo domain.LogRepository, habitRepo domain.HabitRepository) *LogUsecase {
	return &LogUsecase{logRepo: logRepo, habitRepo: habitRepo}
}

func (u *LogUsecase) RecordLog(userId, habitId int, date time.Time, isCompleted bool) error {
	habit, err := u.habitRepo.FindById(habitId)
	if err != nil {
		return fmt.Errorf("習慣取得に失敗: %w", err)
	}
	if habit == nil {
		return ErrHabitNotFound
	}
	if habit.UserId != userId {
		return ErrHabitForbidden
	}

	exists, err := u.logRepo.ExistsByHabitIdAndDate(habitId, date)
	if err != nil {
		return fmt.Errorf("ログ確認に失敗: %w", err)
	}
	if exists {
		return ErrLogAlreadyExists
	}

	l := &domain.HabitLog{HabitId: habitId, IsCompleted: isCompleted, Date: date}
	if err := u.logRepo.Save(l); err != nil {
		return fmt.Errorf("ログ保存に失敗: %w", err)
	}
	return nil
}

func (u *LogUsecase) GetLogs(userId int) ([]*domain.HabitLog, error) {
	logs, err := u.logRepo.FindByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("ログ取得に失敗: %w", err)
	}
	return logs, nil
}
