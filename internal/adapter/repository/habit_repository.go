package repository

import (
	"database/sql"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type sqlHabitRepository struct {
	db *sql.DB
}

// 通信が確立されたDBの器を受け取り、DBの実体を返す
func NewSqlHabitRepository(habitDB *sql.DB) *sqlHabitRepository {
	return &sqlHabitRepository{db: habitDB}
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

// func (r *sqlHabitRepository) Update(*domain.Habit) {
// 	// query := `
// 	// 	UPDATE habits
// 	// 	SET
// 	// `
// }

func (r *sqlHabitRepository) Update(h *domain.Habit) error {
	return nil // とりあえず中身は空でも、メソッドが存在すれば「合格」になる
}

func (r *sqlHabitRepository) Delete(id int) error {
	return nil
}
