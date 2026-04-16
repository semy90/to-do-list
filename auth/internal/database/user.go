package database

import (
	"auth/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{pool: pool}
}

type UserStorage struct {
	pool *pgxpool.Pool
}

func (s *UserStorage) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	row := s.pool.QueryRow(context.Background(), "SELECT * FROM users WHERE email = $1", email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.HashPassword)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserStorage) AddUser(name, email, hashPassword string) (int, error) {
	var id int
	err := s.pool.QueryRow(context.Background(), "INSERT INTO users (name, email, hash_pass) VALUES ($1, $2, $3) RETURNING id;", name, email, hashPassword).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
func (s *UserStorage) DelUser(id int) {
	tag, err := s.pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		// do something
	}
	if tag.RowsAffected() == 0 {
		//do something
	}
}
