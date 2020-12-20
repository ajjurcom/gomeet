package repository

import (
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"errors"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type ICampusRepository interface {
	Conn() error
	InsertCampus(campus model.Campus) error
	DeleteCampus(id int) error
	UpdateCampus(campus model.Campus) error
	SelectCount() (int, error)
	SelectAllCampus() ([]model.Campus, error)
	SelectCampusByPage(page, onePageCount int) ([]model.Campus, error)
	SelectCampusByID(id int) (model.Campus, error)
	IsCampusExists(id int) (bool, error)
}

func NewCampusRepository(table string) ICampusRepository {
	return &CampusManagerRepository{table, database.MysqlDB}
}

type CampusManagerRepository struct {
	table string
	mysqlConn *sqlx.DB
}

func (cmr *CampusManagerRepository) Conn() error {
	if cmr.mysqlConn == nil {
		if err := database.InitMysql(); err != nil {
			return err
		}
		cmr.mysqlConn = database.MysqlDB
	}
	if cmr.table == "" {
		cmr.table = "campus"
	}
	return nil
}

// InsertCampus 新增校区
func (cmr *CampusManagerRepository) InsertCampus(campus model.Campus) error {
	if err := cmr.Conn(); err != nil {
		return err
	}

	sqlStr := "insert into " + cmr.table + "(campus_name) values(?)"

	stmt, err := cmr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(campus.CampusName)
	return err
}

// DeleteCampus 删除校区
func (cmr *CampusManagerRepository) DeleteCampus(id int) error {
	if err := cmr.Conn(); err != nil {
		return err
	}

	sqlStr := "delete from " + cmr.table + " where id = ?"
	stmt, err := cmr.mysqlConn.Prepare(sqlStr)
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

// UpdateCampus 修改校区信息
func (cmr *CampusManagerRepository) UpdateCampus(campus model.Campus) error {
	if err := cmr.Conn(); err != nil {
		return err
	}

	sqlStr := "update " + cmr.table + " set campus_name = ? where id = ?"

	stmt, err := cmr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(campus.CampusName, campus.ID)
	if err != nil {
		return err
	}

	num, _ := result.RowsAffected()
	if num == 0 {
		return errors.New("no exists")
	}
	return nil
}

// SelectCount 获取有多少个校区
func (cmr *CampusManagerRepository) SelectCount() (count int, err error) {
	if err = cmr.Conn(); err != nil {
		return
	}

	sqlStr := "select count(*) from " + cmr.table
	err = cmr.mysqlConn.QueryRow(sqlStr).Scan(&count)
	return
}

// SelectCampusByPage 翻页查询校区
func (cmr *CampusManagerRepository) SelectCampusByPage(page, onePageCount int) (campus []model.Campus, err error) {
	if err = cmr.Conn(); err != nil {
		return
	}

	page -= 1
	startIndex := strconv.Itoa(page * onePageCount)
	sqlStr := "select id, campus_name from " + cmr.table + " order by id limit " + startIndex + ", " + strconv.Itoa(onePageCount)
	err = cmr.mysqlConn.Select(&campus, sqlStr)
	return
}

// SelectAllCampus 查询所有校区信息
func (cmr *CampusManagerRepository) SelectAllCampus() (campus []model.Campus, err error) {
	if err = cmr.Conn(); err != nil {
		return
	}

	sqlStr := "select id, campus_name from " + cmr.table

	err = cmr.mysqlConn.Select(&campus, sqlStr)
	return
}

func (cmr *CampusManagerRepository) SelectCampusByID(id int) (campus model.Campus, err error) {
	if err = cmr.Conn(); err != nil {
		return
	}

	sql := "select campus_name from " + cmr.table + " where id=? limit 1"
	err = cmr.mysqlConn.Get(&campus, sql, id)
	return
}

// IsExists 查询校区是否存在
func (cmr *CampusManagerRepository) IsCampusExists(id int) (bool, error) {
	if err := cmr.Conn(); err != nil {
		return false, err
	}

	sqlStr := "select 1 from " + cmr.table + " where id = ?"

	num := 0
	cmr.mysqlConn.QueryRow(sqlStr, id).Scan(&num)

	if num == 1 {
		return true, nil
	}
	return false, nil
}
