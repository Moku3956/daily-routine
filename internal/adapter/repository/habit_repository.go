package repository

import (
	"database/sql"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type sqlHabitRepository struct {
	db *sql.DB
}

func (r *sqlHabitRepository) Save(h *domain.Habit) error {
	query := `
		INSERT INTO habits
			(user_id, habit_name, category, value, unit, p_cost, m_cost, must)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING habit_id
	`
	err := r.db.QueryRow(
		query,
		h.UserId,
		h.HabitName,
		h.Category,
		h.Value,
		h.Unit,
		h.PCost,
		h.MCost,
		h.Must,
	).Scan(&h.HabitId)
	if err != nil {
		return err
	}
	return nil
}
