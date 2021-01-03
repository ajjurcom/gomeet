package config

import "gopkg.in/ini.v1"

// Cfg 是配置文件操作的句柄
var Cfg *ini.File

// InitConfig 初始化配置文件操作句柄，返回类型 error，如果返回错误，必须 panic
func InitConfig() (err error) {
	if Cfg, err = ini.Load("./config/my.ini"); err != nil {
		return err
	}
	return nil
}
