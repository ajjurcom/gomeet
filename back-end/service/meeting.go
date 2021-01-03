package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"strconv"
)

type IMeetingService interface {
	AddMeeting(meeting model.Meeting) error
	DeleteMeeting(id int) error
	UpdateMeeting(meeting model.Meeting) error
	GetMeetingByID(id int) (model.Meeting, error)
	GetMeetingCountByBuilding(buildingID int) (int, error)
	GetMeetingsByPage(page, onePageCount, buildingID int) ([]model.Meeting, error)
	GetMeetingsByID(ids string) ([]model.Meeting, error)
	GetMeetingCountByKeyword(keyword string) (int, error)
	GetMeetingsByKeyword(page, onePageCount int, keyword string)  ([]model.Meeting, error)
	GetAllMeetingTypes() []string
	GetAllScaleTypes() []string
	GetAllMeetingByBuilding(buildingID int) ([]model.Meeting, error)
	GetAllMeetingsByParams(buildingID int, layer int, meetingType []string, scales []string) ([]model.Meeting, error)
	GetMeetingByInfo(meetingsID, campusID, meetingType, meetingScale string) (model.Meeting, error)
	GetAllMeetingsID() ([]int, error)
}

func NewMeetingService(meetingRepo repository.IMeetingRepository) IMeetingService {
	return &MeetingService{meetingRepo}
}

type MeetingService struct {
	MeetingRepository repository.IMeetingRepository
}

// AddMeeting 添加会议室
func (ms *MeetingService) AddMeeting(meeting model.Meeting) error {
	return ms.MeetingRepository.InsertMeeting(meeting)
}

func (ms *MeetingService) DeleteMeeting(id int) error {
	return ms.MeetingRepository.DeleteMeeting(id)
}

func (ms *MeetingService) UpdateMeeting(meeting model.Meeting) error {
	//// 确保建筑存在
	//exists := false
	//var err error
	//if exists, err = ms.BuildingRepository.IsBuildingExists(meeting.BuildingID); err != nil {
	//	return err
	//}
	//if !exists {
	//	return errors.New("the building no exists")
	//}
	//// 确保会议室类型符合
	//if !model.IsMeetingType(meeting.MeetingType) {
	//	return errors.New("the meeting type no exists")
	//}
	//// 确保会议室容量符合
	//if !model.IsScaleType(meeting.Scale) {
	//	return errors.New("the meeting scale no exists")
	//}
	return ms.MeetingRepository.UpdateMeeting(meeting)
}

func (ms *MeetingService) GetMeetingByID(id int) (model.Meeting, error) {
	return ms.MeetingRepository.SelectMeetingByID(id)
}

func (ms *MeetingService) GetMeetingCountByBuilding(buildingID int) (int, error) {
	return ms.MeetingRepository.SelectMeetingCount("building_id", strconv.Itoa(buildingID), true)
}

func (ms *MeetingService) GetMeetingCountByKeyword(keyword string) (int, error) {
	return ms.MeetingRepository.SelectMeetingCount("meeting_name", keyword, false)
}

func (ms *MeetingService) GetMeetingsByKeyword(page, onePageCount int, keyword string)  ([]model.Meeting, error) {
	return ms.MeetingRepository.SearchMeetingsByKeyword(page, onePageCount, keyword)
}

func (ms *MeetingService) GetMeetingsByPage(page, onePageCount, buildingID int) ([]model.Meeting, error) {
	return ms.MeetingRepository.SelectMeetingsByBuilding(buildingID, page, onePageCount)
}

func (ms *MeetingService) GetMeetingsByID(ids string) ([]model.Meeting, error) {
	return ms.MeetingRepository.SelectMeetingsByID(ids)
}

func (ms *MeetingService) GetAllMeetingTypes() []string {
	return ms.MeetingRepository.SelectAllMeetingTypes()
}

func (ms *MeetingService) GetAllScaleTypes() []string {
	return ms.MeetingRepository.SelectAllScaleTypes()
}

func (ms *MeetingService) GetAllMeetingByBuilding(buildingID int) ([]model.Meeting, error) {
	return ms.MeetingRepository.SelectMeetingsByBuilding(buildingID)
}

func (ms *MeetingService) GetAllMeetingsByParams(buildingID int, layer int, meetingType []string, scales []string) ([]model.Meeting, error) {
	return ms.MeetingRepository.SelectAllMeetingsByParams(buildingID, layer, meetingType, scales)
}

func (ms *MeetingService) GetMeetingByInfo(meetingsID, campusID, meetingType, meetingScale string) (model.Meeting, error) {
	return ms.MeetingRepository.SelectMeetingByInfo(meetingsID, campusID, meetingType, meetingScale)
}

func (ms *MeetingService) GetAllMeetingsID() ([]int, error) {
	meetings, err := ms.MeetingRepository.SelectAllMeetingsID()
	if err != nil {
		return []int{}, err
	}
	ids := make([]int, len(meetings))
	for i, v := range meetings {
		ids[i] = v.ID
	}
	return ids, nil
}

