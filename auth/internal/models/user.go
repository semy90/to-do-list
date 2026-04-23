package models

type User struct {
	Id           int    `json:"id"  db:"id"`
	Name         string `json:"name"  db:"name"`
	Email        string `json:"email" db:"email"`
	HashPassword string `json:"pass" db:"hash_pass"`
}
