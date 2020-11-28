package repository

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type IAppointmentRepository interface {
	Conn() error
	Add(appointment model.Appointment, members string) error
	Delete(id int, members string) error
	IsAppointmentExists(appointment model.Appointment) (bool, error)
	SelectOneByCondition(conditionName, conditionVal string, attrNames ...string) (model.Appointment, error)
}

func NewAppointmentRepository(appointmentTable, userTable string) IAppointmentRepository {
	return &AppointmentRepository{appointmentTable, userTable, database.MysqlDB}
}

type AppointmentRepository struct {
	appointmentTable string
	userTable string
	mysqlConn *sqlx.DB
}

func (ar *AppointmentRepository) Conn() error {
	if ar.mysqlConn == nil {
		if err := database.InitMysql(); err != nil {
			return err
		}
		ar.mysqlConn = database.MysqlDB
	}
	if ar.appointmentTable == "" {
		ar.appointmentTable = "appointment"
	}
	if ar.userTable == "" {
		ar.userTable = "user"
	}
	return nil
}

func (ar *AppointmentRepository) Add(appointment model.Appointment, members string) error {
	/*
	 * 1. 创建新会议，获得新会议id号
	 * 2. 将新会议id号添加到members每个user的appointments中
	 */
	if err := ar.Conn(); err != nil {
		return err
	}

	tx, err := ar.mysqlConn.Begin()
	if err != nil {
		return err
	}

	// 1. 创建新会议，获得新会议id号
	var id int64 = -1
	sqlStr := "insert into " + ar.appointmentTable + "(creator_id, creator_name, meeting_id, day, start_time, " +
		"end_time, theme, content, group_list, members, all_members) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	fmt.Println("sqlStr -> ", sqlStr)
	result, err := tx.Exec(sqlStr, appointment.CreatorID, appointment.CreateName, appointment.MeetingID, appointment.Day,
		appointment.StartTime, appointment.EndTime, appointment.Theme, appointment.Content,
		appointment.Groups, appointment.Members, members)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err = result.LastInsertId()
	if err != nil || id == -1 {
		tx.Rollback()
		return err
	}
	// 2. 将新会议id号添加到members每个user的appointments中
	sqlStr = "update " + ar.userTable + " set appointments=concat(appointments, '," + strconv.Itoa(int(id)) + "') where id in (" + members + ")"
	_, err = tx.Exec(sqlStr)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (ar *AppointmentRepository) Delete(id int, members string) error {
	/*
	 * 1. 查询所有成员的group_list
	 * 2. 删除该会议室
	 * 3. 从所有成员中删除该会议
	 */
	if err := ar.Conn(); err != nil {
		return err
	}

	// 1. 查询所有成员的group_list
	sql := "select id, appointments from " + ar.userTable + " where id in (" + members + ")"
	var users []model.User
	if err := ar.mysqlConn.Select(&users, sql); err != nil {
		return err
	}

	// 2. 删除该会议室
	tx, err := ar.mysqlConn.Begin()
	if err != nil {
		return err
	}
	sql = "delete from " + ar.appointmentTable + " where id = ?"
	if _, err = tx.Exec(sql, id); err != nil {
		tx.Rollback()
		return err
	}
	// 3. 从所有成员中删除该会议
	idStr := strconv.Itoa(id)
	sql = "update " + ar.userTable + " set appointments=? where id=?"
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
	tx.Commit()
	return nil
}

func (ar *AppointmentRepository) IsAppointmentExists(appointment model.Appointment) (isExists bool, err error) {
	if err = ar.Conn(); err != nil {
		return
	}

	sql := "select 1 from " + ar.appointmentTable + " where meeting_id=? and day=? " +
		"and ((start_time >= ? and start_time < ?) " +
		"or (end_time > ? and end_time <= ?)) " +
		"or (start_time <= ? and end_time >= ?)"
	num := 0
	ar.mysqlConn.QueryRow(sql, appointment.MeetingID, appointment.Day, appointment.StartTime,
		appointment.EndTime, appointment.StartTime, appointment.EndTime, appointment.StartTime, appointment.EndTime).Scan(&num)
	if num == 1 {
		isExists = true
	}
	return
}

func (ar *AppointmentRepository) SelectOneByCondition(conditionName, conditionVal string, attrNames ...string) (appointment model.Appointment, err error) {
	if err = ar.Conn(); err != nil {
		return
	}

	attrs := ""
	if len(attrNames) == 0 {
		attrs = "*"
	} else {
		for _, v := range attrNames {
			attrs += v + ","
		}
		attrs = strings.Trim(attrs, ",")
	}

	sql := "select " + attrs + " from " + ar.appointmentTable + " where " + conditionName + " = ? limit 1"
	err = ar.mysqlConn.Get(&appointment, sql, conditionVal)
	return
}
