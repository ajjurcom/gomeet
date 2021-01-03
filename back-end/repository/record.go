package repository

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type IRecordRepository interface {
	Conn() error
	TransferAppointment(appointment model.Appointment, members string) error
	SelectAppointmentByDay(startDay, endDay string) ([]model.Appointment, error)
}

type RecordRepository struct {
	appointmentTable string
	recordTable string
	userTable string
	mysqlConn *sqlx.DB
}

func NewRecordRepository(appointmentTable, userTable, recordTable string) IRecordRepository {
	return &RecordRepository{
		appointmentTable: appointmentTable,
		recordTable: recordTable,
		userTable: userTable,
		mysqlConn: database.MysqlDB,
	}
}

func (rr *RecordRepository) Conn() error {
	if rr.mysqlConn == nil {
		if err := database.InitMysql(); err != nil {
			return err
		}
		rr.mysqlConn = database.MysqlDB
	}
	if rr.appointmentTable == "" {
		rr.appointmentTable = "appointment"
	}
	if rr.recordTable == "" {
		rr.recordTable = "record"
	}
	return nil
}

func (rr *RecordRepository) TransferAppointment(appointment model.Appointment, members string) error {
	/*
	 * 1. 插入 日志表
	 * 2. 从会议表删除
	 */
	if err := rr.Conn(); err != nil {
		return err
	}

	tx, err := rr.mysqlConn.Begin()
	if err != nil {
		return err
	}
	// 1. 插入 日志表
	sql := "insert into " + rr.recordTable + "(creator_id, creator_name, meeting_id, day, start_time, end_time, state) " +
		"values(?,?,?,?,?,?,?)"
	if _, err = tx.Exec(sql, appointment.CreatorID, appointment.CreatorName, appointment.MeetingID, appointment.Day,
		appointment.StartTime, appointment.EndTime, appointment.State); err != nil {
		tx.Rollback()
		return err
	}
	// 2. 从会议表删除
	// 查询所有成员的group_list
	var users []model.User
	if members != "" {
		sql = "select id, appointments from " + rr.userTable + " where id in (" + members + ")"
		if err := rr.mysqlConn.Select(&users, sql); err != nil {
			return err
		}
	}
	// 删除该会议室
	sql = "delete from " + rr.appointmentTable + " where id = ?"
	if _, err = tx.Exec(sql, appointment.ID); err != nil {
		tx.Rollback()
		return err
	}
	// 从所有成员中删除该会议
	if members != "" {
		idStr := strconv.Itoa(appointment.ID)
		sql = "update " + rr.userTable + " set appointments=? where id=?"
		for _, user := range users {
			tmp := common.MemberStrToList(user.Appointments)
			index := common.StrIndexOf(tmp, idStr)
			if index != -1 {
				tmp = append(tmp[:index], tmp[index+1:]...)
				_, err := tx.Exec(sql, common.MemberListToStr(tmp), user.ID)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	tx.Commit()
	return nil
}

func (rr *RecordRepository) SelectAppointmentByDay(startDay, endDay string) (appointments []model.Appointment, err error) {
	if err = rr.Conn(); err != nil {
		return
	}

	sql := "select * from " + rr.recordTable + " where day >= ? and day <= ?"
	err = rr.mysqlConn.Select(&appointments, sql, startDay, endDay)
	return
}

