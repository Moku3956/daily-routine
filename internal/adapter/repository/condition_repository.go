package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type sqlConditionRepository struct {
	db *sql.DB
}

func NewSqlConditionRepository(db *sql.DB) *sqlConditionRepository {
	return &sqlConditionRepository{db: db}
}

func (r *sqlConditionRepository) Save(c *domain.Condition) error {
	query := `
		INSERT INTO conditions (user_id, p_condition, m_condition, date)
		VALUES ($1, $2, $3, $4)
		RETURNING condition_id
	`
	return r.db.QueryRow(query, c.UserId, c.PCondition, c.MCondition, c.Date).Scan(&c.ConditionId)
}

func (r *sqlConditionRepository) FindByUserIdAndDate(userId int, date time.Time) (*domain.Condition, error) {
	query := `
		SELECT condition_id, user_id, p_condition, m_condition, date
		FROM conditions
		WHERE user_id = $1 AND date = $2
	`
	c := &domain.Condition{}
	err := r.db.QueryRow(query, userId, date).Scan(
		&c.ConditionId, &c.UserId, &c.PCondition, &c.MCondition, &c.Date,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	c.Rank = domain.CalcConditionRank(c.PCondition, c.MCondition)
	return c, nil
}
