package models

type Task struct {
	Id     int    `json:"id"  db:"id"`
	Text   string `json:"text" db:"text"`
	UserId int    `json:"user_id" db:"user_id"`
}
