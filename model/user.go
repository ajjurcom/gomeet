package model

import "github.com/golang/protobuf/ptypes/timestamp"

type User struct {
	ID int `json:"id" db:"id"`
	Sno string	`json:"sno" db:"sno" binding:"required"`
	Phone string `json:"phone" db:"phone" binding:"required"`
	Password string `json:"password" db:"password" binding:"required"`
	IsAdmin int `json:"is_admin" db:"is_admin"`
	State string `json:"state" db:"state"`
	Ban timestamp.Timestamp `json:"ban" db:"ban"`
	CampusID int `json:"campus_id" db:"campus_id" binding:"required"`
	Username string `json:"username" db:"username" binding:"required"`
	Email string `json:"email" db:"email" binding:"required,email"`
}

const (
	StateVerifyUser = "verify_user"
	StateVerifyAdmin = "verify_admin"
	StateRefuse = "refuse"
	StateNormal = "normal"
	StateBlacklist = "blacklist"
)
