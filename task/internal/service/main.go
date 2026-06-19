package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"task/internal/database"
	"task/internal/models"
	"task/jwt_utils"
	"task/utils"

	"go.uber.org/zap"
)

type TaskCRUD struct {
	Storage database.TaskStorage
	Ctx     context.Context
}

func (t *TaskCRUD) GetTask(w http.ResponseWriter, r *http.Request) {
	logger, _ := t.Ctx.Value(("logger")).(*zap.Logger)
	userId, _ := strconv.Atoi(r.Header.Get("X-User-Id"))
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		logger.Info("id in path is not a number", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	task, err := t.Storage.GetTask(t.Ctx, id, userId)
	if err != nil {
		logger.Info("task not found", zap.Int("id", id), zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	taskByte, err := json.Marshal(task)
	if err != nil {
		logger.Info("json marshal error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(taskByte)
}

func (t *TaskCRUD) PostTask(w http.ResponseWriter, r *http.Request) {
	logger, _ := t.Ctx.Value(("logger")).(*zap.Logger)
	userId, _ := strconv.Atoi(r.Header.Get("X-User-Id"))
	var task models.Task
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Info("cant read the body", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	json.Unmarshal(body, &task)
	id, err := t.Storage.AddTask(t.Ctx, task.Text, userId)
	if err != nil {
		logger.Info("cant add task", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	idStr := strconv.Itoa(id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(idStr))
}
func (t *TaskCRUD) DelTask(w http.ResponseWriter, r *http.Request) {
	logger, _ := t.Ctx.Value(("logger")).(*zap.Logger)
	userId, _ := strconv.Atoi(r.Header.Get("X-User-Id"))
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Info("id in path is not a number", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	t.Storage.DelTask(t.Ctx, id, userId)
	w.WriteHeader(http.StatusOK)
}

func (t *TaskCRUD) GetTaskFromTo(w http.ResponseWriter, r *http.Request) {
	logger, _ := t.Ctx.Value(("logger")).(*zap.Logger)
	userId, _ := strconv.Atoi(r.Header.Get("X-User-Id"))
	from, to, err := utils.ParseQuery(r)
	if err != nil {
		logger.Info("query error", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	offset := from - 1
	limit := to - from
	tasks, err := t.Storage.GetTasksByLimitAndOffset(t.Ctx, userId, limit, offset)
	if err != nil {
		logger.Info("cant get user tasks", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	tasksByte, err := json.Marshal(tasks)
	if err != nil {
		logger.Info("json marshal err", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(tasksByte)
}

func (t *TaskCRUD) UpdateTask(w http.ResponseWriter, r *http.Request) {
	logger, _ := t.Ctx.Value(("logger")).(*zap.Logger)
	userId, _ := strconv.Atoi(r.Header.Get("X-User-Id"))
	var task models.Task
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Info("cant read the body", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	json.Unmarshal(body, &task)
	t.Storage.UpdateTask(context.Background(), task.Text, task.Id, userId)
	w.WriteHeader(http.StatusOK)
}

// переписать так чтобы не трогать боди а просто записывать в хедер юз айди
func (t *TaskCRUD) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, _ := t.Ctx.Value(("logger")).(*zap.Logger)
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			logger.Info("authorization header not found", zap.Error(err))
			http.Error(w, "go to authorization", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(cookie.Value, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Info("wrong token", zap.Any("parts", parts))
			http.Error(w, "wrong token", http.StatusUnauthorized)
			return
		}
		userId, err := jwt_utils.GetIdFromToken(parts[1])
		if err != nil {
			logger.Info("wrong token", zap.Error(err))
			http.Error(w, "wrong token", http.StatusUnauthorized)
			return
		}
		r.Header.Add("X-User-Id", strconv.Itoa(userId))
		next.ServeHTTP(w, r)
	})
}
