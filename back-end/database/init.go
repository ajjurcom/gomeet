package database

import (
	"com/mittacy/gomeet/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MysqlDB mysql数据库操作句柄
var MysqlDB *sqlx.DB

// InitMysql 连接mysql
//
// 返回值 error，如果返回错误需要 panic
func InitMysql() error {
	mysqlCfg, err := config.Cfg.GetSection("mysql")
	if err != nil {
		return err
	}
	var (
		user = mysqlCfg.Key("user").String()
		port = mysqlCfg.Key("port").String()
		password = mysqlCfg.Key("password").String()
		host = mysqlCfg.Key("host").String() + ":" + port
		database = mysqlCfg.Key("database").String()
	)
	
	par := user + ":" + password + "@tcp(" + host + ")/" + database
	MysqlDB, err = sqlx.Open("mysql", par)
	if err != nil {
		return err
	}
	return MysqlDB.Ping()
}

// CloseMysql 关闭mysql连接
func CloseMysql() {
	MysqlDB.Close()
}
