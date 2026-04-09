package repository

import (
	"database/sql"
	"errors"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type sqlSessionRepository struct {
	db *sql.DB
}

func NewSqlSessionRepository(db *sql.DB) *sqlSessionRepository {
	return &sqlSessionRepository{db: db}
}

func (r *sqlSessionRepository) Save(s *domain.Session) error {
	query := `INSERT INTO sessions (user_id, token) VALUES ($1, $2) RETURNING session_id`
	return r.db.QueryRow(query, s.UserId, s.Token).Scan(&s.SessionId)
}

func (r *sqlSessionRepository) FindByToken(token string) (*domain.Session, error) {
	query := `SELECT session_id, user_id, token FROM sessions WHERE token = $1`
	s := &domain.Session{}
	err := r.db.QueryRow(query, token).Scan(&s.SessionId, &s.UserId, &s.Token)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *sqlSessionRepository) DeleteByToken(token string) error {
	_, err := r.db.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}
