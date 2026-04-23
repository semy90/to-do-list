package database

import (
	"context"
	"task/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewTaskStorage(pool *pgxpool.Pool) *TaskStorage {
	return &TaskStorage{pool: pool}
}

type TaskStorage struct {
	pool *pgxpool.Pool
}

func (s *TaskStorage) GetTasksByLimitAndOffset(ctx context.Context, userId, limit, offset int) ([]*models.Task, error) {
	rows, err := s.pool.Query(context.Background(), "SELECT * FROM tasks WHERE user_id = $1 LIMIT $2 OFFSET $3", userId, limit, offset)
	if err != nil {
		return nil, err
	}
	tasks := []*models.Task{}
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.Id, &task.Text, &task.UserId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskStorage) GetTask(ctx context.Context, taskId int) (*models.Task, error) {
	// const op = "TaskStorage.GetTask"
	// logger, _ := ctx.Value(("logger")).(*zap.Logger)

	task := &models.Task{}
	row := s.pool.QueryRow(context.Background(), "SELECT * FROM tasks WHERE id=$1", taskId)
	err := row.Scan(&task.Id, &task.Text, &task.UserId)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskStorage) AddTask(ctx context.Context, text string, userId int) (int, error) {
	const op = "TaskStorage.GetTask"
	logger, _ := ctx.Value(("logger")).(*zap.Logger)
	var id int
	err := s.pool.QueryRow(context.Background(), "INSERT INTO tasks (text, user_id) VALUES ($1, $2) RETURNING id;", text, userId).Scan(&id)
	if err != nil {
		logger.Info("add task err", zap.Int("user_id", userId), zap.String("text", text), zap.Error(err))
		return -1, err
	}
	logger.Info("task added successfuly", zap.Int("task_id", id))
	return id, nil
}
func (s *TaskStorage) EditTask(ctx context.Context, text string, id int) {
	tag, err := s.pool.Exec(context.Background(), "UPDATE tasks SET text = '$1' WHERE id=$2", text, id)
	if err != nil {
		// do something
	}
	if tag.RowsAffected() == 0 {
		//do something
	}
}

func (s *TaskStorage) DelTask(ctx context.Context, id int) {
	// const op = "TaskStorage.GetTask"
	// logger, _ := ctx.Value(("logger")).(*zap.Logger)
	tag, err := s.pool.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		// do something
	}
	if tag.RowsAffected() == 0 {
		//do something
	}
}
