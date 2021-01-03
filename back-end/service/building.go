package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"strconv"
)

type IBuildingService interface {
	AddBuilding(building model.Building) error
	DeleteBuilding(id int) error
	UpdateBuilding(building model.Building) error
	GetBuildingByID(id int) (model.Building, error)
	GetBuildingsByPage(page, onePageCount, campusID int) ([]model.Building, error)
	GetBuildingsByKeyword(page, onePageCount int, keyword string) ([]model.Building, error)
	GetBuildingCountByCampus(campusID int) (int, error)
	GetBuildingCountByKeyword(buildingKeyword string) (int, error)
	GetAllBuildingsByCampus(campusID int) ([]model.Building, error)
	GetBuildingLayer(campusID int) (int, error)
	IsBuildingExists(buildingID int) (bool, error)
}

func NewBuildingService(buildingRepo repository.IBuildingRepository) IBuildingService {
	return &BuildingService{buildingRepo}
}

type BuildingService struct {
	BuildingRepository repository.IBuildingRepository
}

func (bs *BuildingService) AddBuilding(building model.Building) error {
	return bs.BuildingRepository.InsertBuilding(building)
}

func (bs *BuildingService) DeleteBuilding(id int) error {
	return bs.BuildingRepository.DeleteBuilding(id)
}

func (bs *BuildingService) UpdateBuilding(building model.Building) error {
	//// 确保campus_id存在
	//isExists, err := bs.CampusRepository.IsExists(building.CampusID)
	//if err != nil {
	//	return err
	//}
	//if !isExists {
	//	return errors.New("the campus no exists")
	//}
	return bs.BuildingRepository.UpdateBuilding(building)
}

func (bs *BuildingService) GetBuildingsByPage(page, onePageCount, campusID int) ([]model.Building, error) {
	return bs.BuildingRepository.SelectBuildingsByPage(page, onePageCount, campusID)
}

func (bs *BuildingService) GetBuildingsByKeyword(page, onePageCount int, keyword string) ([]model.Building, error) {
	return bs.BuildingRepository.SearchBuildingsByKeyword(page, onePageCount, keyword)
}

func (bs *BuildingService) GetBuildingByID(id int) (model.Building, error) {
	return bs.BuildingRepository.SelectBuildingByID(id)
}

func (bs *BuildingService) GetBuildingCountByCampus(campusID int) (int, error) {
	return bs.BuildingRepository.SelectBuildingCount("campus_id", strconv.Itoa(campusID), true)
}

func (bs *BuildingService) GetBuildingCountByKeyword(buildingKeyword string) (int, error) {
	return bs.BuildingRepository.SelectBuildingCount("building_name", buildingKeyword, false)
}

func (bs *BuildingService) GetAllBuildingsByCampus(campusID int) ([]model.Building, error) {
	return bs.BuildingRepository.SelectAllBuildingsByCampus(campusID)
}

func (bs *BuildingService) IsBuildingExists(buildingID int) (bool, error) {
	return bs.BuildingRepository.IsBuildingExists(buildingID)
}

func (bs *BuildingService) GetBuildingLayer(campusID int) (int, error) {
	return bs.BuildingRepository.SelectBuildingLayer(campusID)
}