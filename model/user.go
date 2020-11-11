package model

import "github.com/golang/protobuf/ptypes/timestamp"

type User struct {
	ID int `json:"id" db:"id"`
	Sno string	`json:"sno" db:"sno" binding:"required"`
	Phone string `json:"phone" db:"phone" binding:"required"`
	Password string `json:"password" db:"password" binding:"required"`
	State string `json:"state" db:"state"`
	Ban timestamp.Timestamp `json:"ban" db:"ban"`
	Username string `json:"username" db:"username" binding:"required"`
	Email string `json:"email" db:"email" binding:"required,email"`
}

const (
	VerifyUser = "verify_user"
	NormalUser = "normal_user"
	RefuseUser = "refuse_user"
	BlackList = "blacklist"
	VerifyAdmin = "verify_admin"
	NormalAdmin = "normal_admin"
	Root = "root"
)
