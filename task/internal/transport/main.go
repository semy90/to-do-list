package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"task/internal/database"
	"task/internal/models"
)

type TaskCRUD struct {
	Storage database.TaskStorage
}

func (t *TaskCRUD) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	task, err := t.Storage.GetTask(id)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	taskByte, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(taskByte)
}

func (t *TaskCRUD) PostTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	json.Unmarshal(body, &task)

	id, err := t.Storage.AddTask(task.Text, task.UserId)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	idStr := strconv.Itoa(id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(idStr))
}
func (t *TaskCRUD) DelTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	t.Storage.DelTask(id)
}
