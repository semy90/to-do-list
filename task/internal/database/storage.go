package database

import (
	"context"
	"task/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTaskStorage(pool *pgxpool.Pool) *TaskStorage {
	// pool.Exec(context.Background())
	return &TaskStorage{pool: pool}
}

type TaskStorage struct {
	pool *pgxpool.Pool
}

func (s *TaskStorage) GetTask(id int) (*models.Task, error) {
	task := &models.Task{}
	row := s.pool.QueryRow(context.Background(), "SELECT * FROM tasks WHERE id = $1", id)
	err := row.Scan(&task.Id, &task.Text, &task.UserId)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskStorage) AddTask(text string, user_id int) (int, error) {
	var id int
	err := s.pool.QueryRow(context.Background(), "INSERT INTO tasks (text, user_id) VALUES ($1, $2) RETURNING id;", text, user_id).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
func (s *TaskStorage) DelTask(id int) {
	tag, err := s.pool.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		// do something
	}
	if tag.RowsAffected() == 0 {
		//do something
	}
}
