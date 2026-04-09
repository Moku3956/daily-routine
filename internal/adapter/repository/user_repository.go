package repository

import (
	"database/sql"
	"errors"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type sqlUserRepository struct {
	db *sql.DB
}

func NewSqlUserRepository(db *sql.DB) *sqlUserRepository {
	return &sqlUserRepository{db: db}
}

func (r *sqlUserRepository) Save(user *domain.User) error {
	query := `INSERT INTO users (user_name, password) VALUES ($1, $2) RETURNING user_id`
	return r.db.QueryRow(query, user.UserName, user.Password).Scan(&user.UserId)
}

func (r *sqlUserRepository) FindByUserName(userName string) (*domain.User, error) {
	query := `SELECT user_id, user_name, password FROM users WHERE user_name = $1`
	u := &domain.User{}
	err := r.db.QueryRow(query, userName).Scan(&u.UserId, &u.UserName, &u.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}
