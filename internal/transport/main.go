package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"to-do-list/internal/database"
	"to-do-list/internal/models"
)

type TaskCRUD struct {
	Storage database.Storage
}

func (t *TaskCRUD) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	task := t.Storage.GetTask(id)
	if task == nil {
		http.Error(w, "not found", http.StatusNotFound)
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
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	json.Unmarshal(body, &task)
	id := t.Storage.AddTask(task.Text)
	idStr := strconv.Itoa(id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(idStr + " task successfully added"))
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
