package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
)

type IRecordService interface {
	TransferAppointment(model.Appointment, string) error
	GetAppointmentByDay(string, string) ([]model.Appointment, error)
}

func NewRecordService(repository repository.IRecordRepository) IRecordService {
	return &RecordService{repository}
}

type RecordService struct {
	RecordRepository repository.IRecordRepository
}

func (rs *RecordService) TransferAppointment(appointment model.Appointment, members string) error {
	return rs.RecordRepository.TransferAppointment(appointment, members)
}

func (rs *RecordService) GetAppointmentByDay(startDay string, endDay string) (appointments []model.Appointment, err error) {
	return rs.RecordRepository.SelectAppointmentByDay(startDay, endDay)
}