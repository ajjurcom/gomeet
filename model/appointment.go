package model

type Appointment struct {
	ID         int    `json:"id" db:"id"`
	CreatorID  int    `json:"creator_id" db:"creator_id"`
	CreateName string `json:"creator_name" db:"create_name"`
	MeetingID  int    `json:"meeting_id" db:"meeting_id"`
	Day        string `json:"day" db:"day"`
	StartTime  string `json:"start_time" db:"start_time"`
	EndTime    string `json:"end_time" db:"end_time"`
	Theme      string `json:"theme" db:"theme"`
	Content    string `json:"content" db:"content"`
	Groups     string `json:"groups" db:"groups_list"`
	Members    string `json:"members" db:"members"`
	AllMembers string `json:"all_members" db:"all_members"`
}
