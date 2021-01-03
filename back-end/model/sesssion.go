package model

import "github.com/dgrijalva/jwt-go"

type Session struct {
	ID          int    `json:"id"`
	Sno         string `json:"sno"`
	Phone       string `json:"phone"`
	Password    string `json:"password" binding:"required"`
	OldPassword string `json:"old_password"`
	IsAdmin     bool   `json:"is_admin"`
	IsRoot      bool   `json:"is_root"`
	State       string `json:"state"`
	jwt.StandardClaims
}
