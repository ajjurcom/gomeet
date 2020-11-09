package repository

import (
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	Conn() error
	Add(user *model.User) error
	SelectPasswordByAttr(attrName, attrVal string) (string, error)
	IsExistsByAttr(attrName, attrVal string) (bool, error)
	SelectStateByAttr(attrName, attrVal string) (string, error)
	SelectIsAdminByAttr(attrName, attrVal string) (bool, error)
}

func NewUserRepository(table string) IUserRepository {
	return &UserManagerRepository{table, database.MysqlDB}
}

type UserManagerRepository struct {
	table string
	mysqlConn *sqlx.DB
}

func (umr *UserManagerRepository) Conn() error {
	if umr.mysqlConn == nil {
		if err := database.InitMysql(); err != nil {
			return err
		}
		umr.mysqlConn = database.MysqlDB
	}
	if umr.table == "" {
		umr.table = "user"
	}
	return nil
}

func (umr *UserManagerRepository) Add(user *model.User) error {
	if err := umr.Conn(); err != nil {
		return err
	}

	sqlStr := "insert into " + umr.table + "(sno, phone, password, campus_id, username, email) values (?, ?, ?, ?, ?, ?)"
	stmt, err := umr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Sno, user.Phone, user.Password, user.CampusID, user.Username, user.Email)
	return err
}

func (umr *UserManagerRepository) SelectPasswordByAttr(attrName, attrVal string) (password string, err error) {
	if err = umr.Conn(); err != nil {
		return
	}

	sqlStr := "select password from " + umr.table + " where " + attrName + " = ?"

	err = umr.mysqlConn.QueryRow(sqlStr, attrVal).Scan(&password)
	return
}

func (umr *UserManagerRepository) IsExistsByAttr(attrName string, attrVal string) (bool, error) {
	if err := umr.Conn(); err != nil {
		return false, err
	}

	sqlStr := "select 1 from " + umr.table + " where " + attrName + " = ? limit 1"

	num := 0
	umr.mysqlConn.QueryRow(sqlStr, attrVal).Scan(&num)
	if num == 1 {
		return true, nil
	}
	return false, nil
}

func (umr *UserManagerRepository) SelectStateByAttr(attrName, attrVal string) (state string, err error) {
	if err = umr.Conn(); err != nil {
		return
	}

	sqlStr := "select state from " + umr.table + " where " + attrName + " = ? limit 1"

	err = umr.mysqlConn.QueryRow(sqlStr, attrVal).Scan(&state)
	return
}

func (umr *UserManagerRepository) SelectIsAdminByAttr(attrName, attrVal string) (bool, error) {
	if err := umr.Conn(); err != nil {
		return false, err
	}

	sqlStr := "select is_admin from " + umr.table + " where " + attrName + " = ? limit 1"

	var isAdmin int
	err := umr.mysqlConn.QueryRow(sqlStr, attrVal).Scan(&isAdmin)
	if isAdmin == 0 {
		return false, err
	}
	return true, err
}
