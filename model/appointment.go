package model

type Appointment struct {
	ID          int    `json:"id" db:"id"`
	CreatorID   int    `json:"creator_id" db:"creator_id"`
	CreatorName string `json:"creator_name" db:"creator_name"`
	MeetingID   int    `json:"meeting_id" db:"meeting_id"`
	Day         string `json:"day" db:"day"`
	StartTime   string `json:"start_time" db:"start_time"`
	EndTime     string `json:"end_time" db:"end_time"`
	State       string `json:"state" db:"state"`
	Theme       string `json:"theme" db:"theme"`
	Content     string `json:"content" db:"content"`
	Members     string `json:"members" db:"members"`
	Locate		string `json:"locate" db:"locate"`
}

const (
	AppointmentVerify = "verify"
	AppointmentAdopt = "adopt"
)

func AppointmentStates() []string {
	return []string{AppointmentVerify, AppointmentAdopt}
}

func IsAppointmentState(state string) bool {
	if state == AppointmentAdopt || state == AppointmentVerify {
		return true
	}
	return false
}
