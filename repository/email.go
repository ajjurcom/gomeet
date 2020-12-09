package repository

import (
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"errors"
	"github.com/jmoiron/sqlx"
	"strings"
)

type IEmailRepository interface {
	Conn() error
	GetContentByName(name string) (string, error)
}

func NewEmailRepository(table string) IEmailRepository {
	return &EmailRepository{table, database.MysqlDB}
}

type EmailRepository struct {
	table string
	mysqlConn *sqlx.DB
}

func (er *EmailRepository) Conn() error {
	if er.mysqlConn == nil {
		if err := database.InitMysql(); err != nil {
			return err
		}
		er.mysqlConn = database.MysqlDB
	}
	if er.table == "" {
		er.table = "email"
	}
	return nil
}

func (er *EmailRepository) GetContentByName(name string) (content string, err error) {
	if strings.Trim(name, " ") == "" {
		err = errors.New("no exists")
		return
	}

	sql := "select content from " + er.table + " where name=?"
	var email model.Email
	err = er.mysqlConn.Get(&email, sql, name)
	return email.Content, err
}