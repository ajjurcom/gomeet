package controller

import "github.com/robfig/cron"

/*
 * 定时任务
 * 1. 定时将过期的会议转移到 日志表 保存
 */
func InitCron() {
	recordController := NewRecordController()

	c := cron.New()
	const transferInterval = "@hourly"	// 整点执行
	//const transferInterval = "*/60 * * * * ?"	// 调式整分执行
	c.AddFunc(transferInterval, recordController.TransferExpireAppointment)
	c.Start()
}
