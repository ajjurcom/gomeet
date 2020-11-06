package model

import "github.com/dgrijalva/jwt-go"

type Session struct {
	Sno      string `json:"sno"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required"`
	IsAdmin  int    `json:"is_admin"`
	IsRoot   int   `json: "is_root"`
	State    string `json:"state"`
	jwt.StandardClaims
}
