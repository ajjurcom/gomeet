package model

type Email struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Content string `json:"content" db:"content"`
}

