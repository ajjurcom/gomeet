package model

type Meeting struct {
	ID            int    `json:"id" db:"id"`
	MeetingName   string `json:"meeting_name" db:"meeting_name" binding:"required"`
	CampusID      int    `json:"campus_id" db:"campus_id"`
	BuildingID    int    `json:"building_id" db:"building_id" binding:"required"`
	BuildingName  string `json:"building_name" db:"building_name"`
	BuildingLayer int    `json:"building_layer" db:"building_layer"`
	Layer         int    `json:"layer" db:"layer" binding:"required"`
	MeetingType   string `json:"meeting_type" db:"meeting_type" binding:"required"`
	Scale         string `json:"scale" db:"scale" binding:"required"`
	RoomNumber    string `json:"room_number" db:"room_number"`
	ReverseCount  int    `json:"reverse_count"`
}

var meetingTypeDict = map[string]bool{
	"普通":  true,
	"多媒体": true,
	"课室类": true,
}

var scaleTypeDict = map[string]bool{
	"微型会议室(最多容纳10人)": true,
	"小型会议室(最多容纳20人)": true,
	"中型会议室(最多容纳30人)": true,
	"大型会议室(容纳30人以上)": true,
}

func IsMeetingType(typeStr string) bool {
	return meetingTypeDict[typeStr]
}

func GetMeetingTypeList() []string {
	var result []string
	for k, _ := range meetingTypeDict {
		result = append(result, k)
	}
	return result
}

func IsScaleType(typeStr string) bool {
	return scaleTypeDict[typeStr]
}

func GetScaleTypeList() []string {
	var result []string
	for k, _ := range scaleTypeDict {
		result = append(result, k)
	}
	return result
}
