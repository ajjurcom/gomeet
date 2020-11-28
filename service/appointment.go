package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"strconv"
)

type IAppointmentService interface {
	CreateAppointment(appointment model.Appointment, members string) error
	DeleteAppointment(id int, members string) error
	IsAppointmentConflict(appointment model.Appointment) (bool, error)
	GetAllMembersAndCreatorIDByID(id int) (string, int, error)
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

func (as *AppointmentService) IsAppointmentConflict(appointment model.Appointment) (bool, error) {
	return as.AppointmentRepository.IsAppointmentExists(appointment)
}

func (as *AppointmentService) GetAllMembersAndCreatorIDByID(id int) (string, int, error) {
	appointment, err := as.AppointmentRepository.SelectOneByCondition("id", strconv.Itoa(id), "creator_id, all_members")
	return appointment.AllMembers, appointment.CreatorID, err
}
