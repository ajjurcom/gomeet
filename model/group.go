package model

type Group struct {
	ID        int    `json:"id" db:"id"`
	Creator   int    `json:"creator" db:"creator"`
	GroupName string `json:"group_name" db:"group_name"`
	Members   string `json:"members" db:"members"`
}
