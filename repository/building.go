package repository

import (
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"errors"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type IBuildingRepository interface {
	Conn() error
	InsertBuilding(building model.Building) error
	DeleteBuilding(id int) error
	UpdateBuilding(building model.Building) error
	SelectBuildingByID(id int) (model.Building, error)
	SelectBuildingsByPage(page, onePageCount, campusID int) ([]model.Building, error)
	SelectBuildingCountByCampus(campusID int) (int, error)
	IsBuildingExists(id int) (bool, error)
	SelectAllBuildingsByCampus(campusID int) ([]model.Building, error)
	SelectBuildingLayer(campusID int) (int, error)
}

func NewBuildingRepository(table string) IBuildingRepository {
	return &BuildingManagerRepository{table, database.MysqlDB}
}

type BuildingManagerRepository struct {
	buildingTable string
	mysqlConn     *sqlx.DB
}

func (bmr *BuildingManagerRepository) Conn() error {
	if bmr.mysqlConn == nil {
		if err := database.InitMysql(); err != nil {
			return err
		}
		bmr.mysqlConn = database.MysqlDB
	}
	if bmr.buildingTable == "" {
		bmr.buildingTable = "building"
	}
	return nil
}

func (bmr *BuildingManagerRepository) InsertBuilding(building model.Building) error {
	if err := bmr.Conn(); err != nil {
		return err
	}

	sqlStr := "insert into " + bmr.buildingTable + "(campus_id, building_name, layer) values(?, ?, ?)"

	stmt, err := bmr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(building.CampusID, building.BuildingName, building.Layer)
	return err
}

func (bmr *BuildingManagerRepository) DeleteBuilding(id int) error {
	if err := bmr.Conn(); err != nil {
		return err
	}

	sqlStr := "delete from " + bmr.buildingTable + " where id = ?"

	stmt, err := bmr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	num, _ := result.RowsAffected()
	if num == 0 {
		return errors.New("no exists")
	}
	return nil
}

func (bmr *BuildingManagerRepository) UpdateBuilding(building model.Building) error {
	if err := bmr.Conn(); err != nil {
		return err
	}

	sqlStr := "update " + bmr.buildingTable + " set campus_id=?, building_name=?,layer=? where id = ?"

	stmt, err := bmr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(building.CampusID, building.BuildingName, building.Layer, building.ID)
	if err != nil {
		return err
	}

	num, _ := result.RowsAffected()
	if num == 0 {
		return errors.New("the building no exists")
	}
	return nil
}

func (bmr *BuildingManagerRepository) SelectBuildingsByPage(page, onePageCount, campusID int) (buildings []model.Building, err error) {
	if err = bmr.Conn(); err != nil {
		return
	}

	page -= 1
	startIndex := strconv.Itoa(page * onePageCount)
	sqlStr := "select id, campus_id, building_name, layer, count from " + bmr.buildingTable + " where campus_id = ? order by id desc limit " + startIndex + ", " + strconv.Itoa(onePageCount)
	err = bmr.mysqlConn.Select(&buildings, sqlStr, campusID)
	return
}

func (bmr *BuildingManagerRepository) SelectBuildingByID(id int) (building model.Building, err error) {
	if err = bmr.Conn(); err != nil {
		return
	}

	sqlStr := "select * from " + bmr.buildingTable + " where id = ? limit 1"

	err = bmr.mysqlConn.Get(&building, sqlStr, id)
	return
}

func (bmr *BuildingManagerRepository) SelectBuildingCountByCampus(campusID int) (count int, err error) {
	if err = bmr.Conn(); err != nil {
		return
	}

	sqlStr := "select count(*) from " + bmr.buildingTable + " where campus_id = ?"
	err = bmr.mysqlConn.QueryRow(sqlStr, campusID).Scan(&count)
	return
}

func (bmr *BuildingManagerRepository) SelectBuildingLayer(campusID int) (count int, err error) {
	if err = bmr.Conn(); err != nil {
		return
	}

	sqlStr := "select layer from " + bmr.buildingTable + " where id = ?"

	err = bmr.mysqlConn.QueryRow(sqlStr, campusID).Scan(&count)
	return
}

func (bmr *BuildingManagerRepository) IsBuildingExists(id int) (bool, error) {
	if err := bmr.Conn(); err != nil {
		return false, err
	}
	sqlStr := "select 1 from " + bmr.buildingTable + " where id = ? limit 1"

	num := 0
	bmr.mysqlConn.QueryRow(sqlStr, id).Scan(&num)
	if num == 1 {
		return true, nil
	}
	return false, nil
}

func (bmr *BuildingManagerRepository) SelectAllBuildingsByCampus(campusID int) (buildings []model.Building, err error) {
	if err = bmr.Conn(); err != nil {
		return
	}

	sqlStr := "select id, building_name, layer from " + bmr.buildingTable + " where campus_id = ?"
	err = bmr.mysqlConn.Select(&buildings, sqlStr, campusID)
	return
}
