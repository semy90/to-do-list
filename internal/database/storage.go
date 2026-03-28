package database

import (
	"crypto/rand"
	"math/big"
	"to-do-list/internal/models"
)

func New() *Storage {
	tasks := map[int]models.Task{}
	return &Storage{tasks: tasks}
}

func newTaskID() int {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(2*1e9))
	n := int(nBig.Int64())
	return n
}

type Storage struct {
	tasks map[int]models.Task
}

func (s *Storage) GetTask(id int) *models.Task {
	task, ok := s.tasks[id]
	if !ok {
		return nil
	}
	return &task
}
func (s *Storage) AddTask(text string) int {
	for {
		id := newTaskID()
		_, ok := s.tasks[id]
		if !ok {
			s.tasks[id] = models.Task{Id: id, Text: text}
			return id
		}
	}
}
func (s *Storage) DelTask(id int) {
	_, ok := s.tasks[id]
	if !ok {
		return
	}
	delete(s.tasks, id)

}
