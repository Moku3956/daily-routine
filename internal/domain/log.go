package domain

import "time"

type HabitLog struct {
	LogId       int
	HabitId     int
	IsCompleted bool
	Date        time.Time
}

type LogRepository interface {
	Save(log *HabitLog) error
	FindByUserId(userId int) ([]*HabitLog, error)
	ExistsByHabitIdAndDate(habitId int, date time.Time) (bool, error)
	CountConsecutiveDays(userId int, today time.Time) (int, error)
}
