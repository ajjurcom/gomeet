package repository

import (
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"errors"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type IUserRepository interface {
	Conn() error
	Add(user *model.User) error
	Delete(id int) error
	Put(user *model.User) error
	UpdateAttr(id int, attrName, attrVal string) error
	SelectPasswordByAttr(attrName, attrVal string) (string, error)
	IsExistsByAttr(attrName, attrVal string) (bool, error)
	SelectStateByAttr(attrName, attrVal string) (string, error)
	SelectUsersByPage(page, onePageCount int, state string) ([]model.User, error)
	SelectUserByID(id int) (model.User, error)
	SelectCountByState(state string) (int, error)
	SelectIDNameByAtr(attrName, attrVal string) (int, string, error)
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

	sqlStr := "insert into " + umr.table + "(sno, phone, password, username, email) values (?, ?, ?, ?, ?)"
	stmt, err := umr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Sno, user.Phone, user.Password, user.Username, user.Email)
	return err
}

func (umr *UserManagerRepository) Delete(id int) error {
	if err := umr.Conn(); err != nil {
		return err
	}

	sqlStr := "delete from " + umr.table + " where id = ?"

	stmt, err := umr.mysqlConn.Prepare(sqlStr)
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
		return errors.New("the user no exists")
	}
	return nil
}

func (umr *UserManagerRepository) Put(user *model.User) error {
	if err := umr.Conn(); err != nil {
		return err
	}

	sqlStr := "update " + umr.table + " set sno=?, phone=?, username=?, email=? where id=?"
	stmt, err := umr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Sno, user.Phone, user.Username, user.Email, user.ID)
	return err
}

func (umr *UserManagerRepository) UpdateAttr(id int, attrName, attrVal string) error {
	if err := umr.Conn(); err != nil {
		return err
	}

	sqlStr := "update " + umr.table + " set " + attrName + "=? where id=?"
	stmt, err := umr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(attrVal, id)
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

func (umr *UserManagerRepository) SelectIDNameByAtr(attrName, attrVal string) (id int, username string, err error) {
	if err = umr.Conn(); err != nil {
		return
	}

	sqlStr := "select id, username from " + umr.table + " where " + attrName + " = ?"

	err = umr.mysqlConn.QueryRow(sqlStr, attrVal).Scan(&id, &username)
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

func (umr *UserManagerRepository) SelectUsersByPage(page, onePageCount int, state string) (userList []model.User, err error) {
	if err = umr.Conn(); err != nil {
		return
	}

	startIndex := strconv.Itoa(page * onePageCount)
	sqlStr := "select id, sno, state, ban, username from " + umr.table + " where state = ? limit " + startIndex + ", " + strconv.Itoa(onePageCount)
	err = umr.mysqlConn.Select(&userList, sqlStr, state)
	return
}

func (umr *UserManagerRepository) SelectUserByID(id int) (user model.User, err error) {
	if err = umr.Conn(); err != nil {
		return
	}
	sqlStr := "select id, sno, phone, state, ban, username, email from " + umr.table + " where id = ?"
	err = umr.mysqlConn.Get(&user, sqlStr, id)
	return
}

func (umr *UserManagerRepository) SelectCountByState(state string) (count int, err error) {
	if err = umr.Conn(); err != nil {
		return
	}

	sqlStr := "select count(*) from " + umr.table + " where state = ?"
	err = umr.mysqlConn.QueryRow(sqlStr, state).Scan(&count)
	return
}