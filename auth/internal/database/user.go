package database

import (
	"auth/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{pool: pool}
}

type UserStorage struct {
	pool *pgxpool.Pool
}

func (s *UserStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "UserStorage.GetUserByEmail"
	logger, _ := ctx.Value(("logger")).(*zap.Logger)

	user := &models.User{}
	row := s.pool.QueryRow(context.Background(), "SELECT * FROM users WHERE email = $1", email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.HashPassword)
	if err != nil {
		logger.Info("user not found", zap.String("email", email), zap.String("path", op))
		return nil, err
	}
	logger.Info("user found", zap.String("email", email), zap.String("path", op))
	return user, nil
}

func (s *UserStorage) AddUser(ctx context.Context, name, email, hashPassword string) (int, error) {
	var id int
	const op = "UserStorage.AddUser"
	logger, _ := ctx.Value(("logger")).(*zap.Logger)

	err := s.pool.QueryRow(context.Background(), "INSERT INTO users (name, email, hash_pass) VALUES ($1, $2, $3) RETURNING id;", name, email, hashPassword).Scan(&id)
	if err != nil {
		logger.Info("add user err", zap.String("email", email), zap.String("path", op), zap.Error(err))
		return -1, err
	}
	logger.Info("add user", zap.String("email", email), zap.String("path", op))
	return id, nil
}
func (s *UserStorage) DelUserById(ctx context.Context, id int) {
	const op = "UserStorage.DelUserById"
	logger, _ := ctx.Value(("logger")).(*zap.Logger)
	tag, err := s.pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		logger.Error("del user err", zap.Int("id", id), zap.String("path", op), zap.Error(err))
	}
	if tag.RowsAffected() == 0 {
	}
	logger.Info("del user", zap.Int("id", id), zap.String("path", op))

}
