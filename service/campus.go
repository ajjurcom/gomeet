package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
)

type ICampusService interface {
	AddCampus(campus model.Campus) error
	DeleteCampus(id int) error
	UpdateCampus(campus model.Campus) error
	GetCampusByPage(page, onePageCount int) ([]model.Campus, error)
	GetCampusCount() (int, error)
	GetAllCampus() ([]model.Campus, error)
	IsCampusExists(id int) (bool, error)
}

func NewCampusService(repository repository.ICampusRepository) ICampusService {
	return &CampusService{repository}
}

type CampusService struct {
	CampusRepository repository.ICampusRepository
}

func (cs *CampusService) AddCampus(campus model.Campus) error {
	return cs.CampusRepository.InsertCampus(campus)
}

func (cs *CampusService) DeleteCampus(id int) error {
	return cs.CampusRepository.DeleteCampus(id)
}

func (cs *CampusService) UpdateCampus(campus model.Campus) error {
	return cs.CampusRepository.UpdateCampus(campus)
}

func (cs *CampusService) GetCampusByPage(page, onePageCount int) ([]model.Campus, error) {
	return cs.CampusRepository.SelectCampusByPage(page, onePageCount)
}

func (cs *CampusService) GetCampusCount() (int, error) {
	return cs.CampusRepository.SelectCount()
}

func (cs *CampusService) GetAllCampus() ([]model.Campus, error) {
	return cs.CampusRepository.SelectAllCampus()
}

func (cs *CampusService)  IsCampusExists(id int) (bool, error) {
	return cs.CampusRepository.IsCampusExists(id)
}