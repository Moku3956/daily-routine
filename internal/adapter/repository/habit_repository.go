package repository

import (
	"database/sql"
	"errors"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type sqlHabitRepository struct {
	db *sql.DB
}

func NewSqlHabitRepository(db *sql.DB) *sqlHabitRepository {
	return &sqlHabitRepository{db: db}
}

func (r *sqlHabitRepository) Save(h *domain.Habit) error {
	query := `
		INSERT INTO habits
			(user_id, habit_name, category, target_value, unit, p_cost, m_cost, must)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING habit_id
	`
	return r.db.QueryRow(query,
		h.UserId, h.HabitName, h.Category, h.TargetValue,
		h.Unit, h.PCost, h.MCost, h.Must,
	).Scan(&h.HabitId)
}

func (r *sqlHabitRepository) Update(h *domain.Habit) error {
	query := `
		UPDATE habits SET
			habit_name   = $1,
			category     = $2,
			target_value = $3,
			unit         = $4,
			p_cost       = $5,
			m_cost       = $6,
			must         = $7
		WHERE habit_id = $8
	`
	_, err := r.db.Exec(query,
		h.HabitName, h.Category, h.TargetValue,
		h.Unit, h.PCost, h.MCost, h.Must, h.HabitId,
	)
	return err
}

func (r *sqlHabitRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM habits WHERE habit_id = $1`, id)
	return err
}

func (r *sqlHabitRepository) FindAll(userId int) ([]*domain.Habit, error) {
	query := `
		SELECT habit_id, user_id, habit_name, category, target_value, unit, p_cost, m_cost, must
		FROM habits
		WHERE user_id = $1
		ORDER BY habit_id
	`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []*domain.Habit
	for rows.Next() {
		h := &domain.Habit{}
		if err := rows.Scan(&h.HabitId, &h.UserId, &h.HabitName, &h.Category,
			&h.TargetValue, &h.Unit, &h.PCost, &h.MCost, &h.Must); err != nil {
			return nil, err
		}
		habits = append(habits, h)
	}
	return habits, rows.Err()
}

func (r *sqlHabitRepository) FindById(id int) (*domain.Habit, error) {
	query := `
		SELECT habit_id, user_id, habit_name, category, target_value, unit, p_cost, m_cost, must
		FROM habits
		WHERE habit_id = $1
	`
	h := &domain.Habit{}
	err := r.db.QueryRow(query, id).Scan(&h.HabitId, &h.UserId, &h.HabitName, &h.Category,
		&h.TargetValue, &h.Unit, &h.PCost, &h.MCost, &h.Must)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return h, nil
}
