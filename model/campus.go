package model

type Campus struct {
	ID         int    `json:"id" db:"id"`
	CampusName string `json:"campus_name" db:"campus_name" binding:"required"`
	Count      string `json:"count" db:"count"`
}
