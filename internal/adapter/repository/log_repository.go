package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type sqlLogRepository struct {
	db *sql.DB
}

func NewSqlLogRepository(db *sql.DB) *sqlLogRepository {
	return &sqlLogRepository{db: db}
}

func (r *sqlLogRepository) Save(l *domain.HabitLog) error {
	query := `INSERT INTO logs (habit_id, is_completed, date) VALUES ($1, $2, $3) RETURNING log_id`
	return r.db.QueryRow(query, l.HabitId, l.IsCompleted, l.Date).Scan(&l.LogId)
}

func (r *sqlLogRepository) FindByUserId(userId int) ([]*domain.HabitLog, error) {
	query := `
		SELECT l.log_id, l.habit_id, l.is_completed, l.date
		FROM logs l
		JOIN habits h ON l.habit_id = h.habit_id
		WHERE h.user_id = $1
		ORDER BY l.date DESC
	`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.HabitLog
	for rows.Next() {
		l := &domain.HabitLog{}
		if err := rows.Scan(&l.LogId, &l.HabitId, &l.IsCompleted, &l.Date); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func (r *sqlLogRepository) ExistsByHabitIdAndDate(habitId int, date time.Time) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM logs WHERE habit_id = $1 AND date = $2`, habitId, date,
	).Scan(&count)
	return count > 0, err
}

// CountConsecutiveDays は今日から遡って全習慣が完了した連続日数を返す
func (r *sqlLogRepository) CountConsecutiveDays(userId int, today time.Time) (int, error) {
	streak := 0
	date := today

	for i := 0; i < 365; i++ {
		query := `
			SELECT
				COUNT(*)                                                  AS total,
				COALESCE(SUM(CASE WHEN l.is_completed THEN 1 ELSE 0 END), 0) AS completed
			FROM logs l
			JOIN habits h ON l.habit_id = h.habit_id
			WHERE h.user_id = $1 AND l.date = $2
		`
		var total, completed int
		if err := r.db.QueryRow(query, userId, date).Scan(&total, &completed); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
			return streak, err
		}
		if total == 0 || completed < total {
			break
		}
		streak++
		date = date.AddDate(0, 0, -1)
	}
	return streak, nil
}
