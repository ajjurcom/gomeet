package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"strconv"
)

type IAppointmentService interface {
	CreateAppointment(appointment model.Appointment, members string) error
	DeleteAppointment(id int, members string) error
	PutAppointment(appointment model.Appointment, addMembers, deleteMembers string) error
	IsAppointmentConflict(appointment model.Appointment, way string) (bool, error)
	GetAllMembersAndCreatorIDByID(id int) (string, int, error)
	GetAllReserve(day, startTime, meetingID string) ([]model.Appointment, error)
}

func NewAppointmentService(repository repository.IAppointmentRepository) IAppointmentService {
	return &AppointmentService{repository}
}

type AppointmentService struct {
	AppointmentRepository repository.IAppointmentRepository
}

func (as *AppointmentService) CreateAppointment(appointment model.Appointment, members string) error {
	return as.AppointmentRepository.Add(appointment, members)
}

func (as *AppointmentService) DeleteAppointment(id int, members string) error {
	return as.AppointmentRepository.Delete(id, members)
}

func (as *AppointmentService) PutAppointment(appointment model.Appointment, addMembers, deleteMembers string) error {
	return as.AppointmentRepository.Put(appointment, addMembers, deleteMembers)
}

// way: post 添加会议/ put 修改会议
func (as *AppointmentService) IsAppointmentConflict(appointment model.Appointment, way string) (isConflict bool, err error) {
	isConflict = true

	limit := 1
	if way == "put" {
		limit = 2
	}

	var appointments []model.Appointment
	if appointments, err = as.AppointmentRepository.SelectConflictAppointments(appointment, limit); err != nil {
		return
	}
	num := len(appointments)

	// way = "post" || "put"
	if num == 0 {
		isConflict = false
		return
	}

	// way = "post"
	if way == "post" {
		return
	}

	// way = "put"
	if num == 1 && appointments[0].ID == appointment.ID {
		isConflict = false
		return
	}
	return
}

func (as *AppointmentService) GetAllMembersAndCreatorIDByID(id int) (string, int, error) {
	appointment, err := as.AppointmentRepository.SelectOneByCondition("id", strconv.Itoa(id), "creator_id, all_members")
	return appointment.AllMembers, appointment.CreatorID, err
}

func (as *AppointmentService) GetAllReserve(day, startTime, meetingID string) ([]model.Appointment, error) {
	return as.AppointmentRepository.SelectCreator(day, startTime, meetingID)
}