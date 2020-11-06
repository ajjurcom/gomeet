package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"errors"
)

type IBuildingService interface {
	AddBuilding(building model.Building) error
	DeleteBuilding(id int) error
	UpdateBuilding(building model.Building) error
	GetBuildingByCampusAndPage(page, onePageCount, campusID int) ([]model.Building, error)
	GetBuildingCountByCampus(campusID int) (int, error)
	GetBuildingByID(id int) (model.Building, error)
}

func NewBuildingService(buildingRepo repository.IBuildingRepository, campusRepo repository.ICampusRepository ) IBuildingService {
	return &BuildingService{buildingRepo, campusRepo}
}

type BuildingService struct {
	BuildingRepository repository.IBuildingRepository
	CampusRepository repository.ICampusRepository
}

func (bs *BuildingService) AddBuilding(building model.Building) error {
	// 确保campus_id存在
	isExists, err := bs.CampusRepository.IsExists(building.CampusID)
	if err != nil {
		return err
	}
	if !isExists {
		return errors.New("the campus is no exists")
	}
	return bs.BuildingRepository.InsertBuilding(building)
}

func (bs *BuildingService) DeleteBuilding(id int) error {
	return bs.BuildingRepository.DeleteBuilding(id)
}

func (bs *BuildingService) UpdateBuilding(building model.Building) error {
	// 确保campus_id存在
	isExists, err := bs.CampusRepository.IsExists(building.CampusID)
	if err != nil {
		return err
	}
	if !isExists {
		return errors.New("the campus no exists")
	}
	return bs.BuildingRepository.UpdateBuilding(building)
}

func (bs *BuildingService) GetBuildingByCampusAndPage(page, onePageCount, campusID int) ([]model.Building, error) {
	return bs.BuildingRepository.SelectBuildingByCampusAndPage(page, onePageCount, campusID)
}

func (bs *BuildingService) GetBuildingByID(id int) (model.Building, error) {
	return bs.BuildingRepository.SelectBuildingByID(id)
}

func (bs *BuildingService) GetBuildingCountByCampus(campusID int) (int, error) {
	return bs.BuildingRepository.SelectBuildingCountByCampus(campusID)
}
