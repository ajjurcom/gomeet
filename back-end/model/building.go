package model

type Building struct {
	ID           int    `json:"id" db:"id"`
	BuildingName string `json:"building_name" db:"building_name" binding:"required"`
	CampusID     int    `json:"campus_id" db:"campus_id" binding:"required"`
	CampusName   string `json:"campus_name" db:"campus_name"`
	Layer        int    `json:"layer" db:"layer" binding:"required"`
	Count        int    `json:"count" db:"count"`
}
