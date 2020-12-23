package repository

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"errors"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

type IAppointmentRepository interface {
	Conn() error
	Add(appointment model.Appointment) error
	Delete(id int, members string) error
	Put(appointment model.Appointment, addMembers, deleteMembers string) error
	PutState(id int, state string) error
	SelectConflictAppointments(appointment model.Appointment, limit int, attrNames ...string) ([]model.Appointment, error)
	SelectAppointmentsIDByTime(appointment model.Appointment) ([]model.Appointment, error)
	SelectOneByCondition(conditionName, conditionVal string, attrNames ...string) (model.Appointment, error)
	SelectCreator(day, startTime, meetingID string) ([]model.Appointment, error)
	SelectAppointmentsByCondition(conditionName, conditionVal string) ([]model.Appointment, error)
	SelectAppointmentsByID(ids string) ([]model.Appointment, error)
	SelectAppointmentsByPage(page, onePageCount int, state string) ([]model.Appointment, error)
	SelectCountByState(state string) (int, error)
	SelectExpireAppointment(day string, endTime string) ([]model.Appointment, error)
	TransferAppointment(appointment model.Appointment, members string) error
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

func (ar *AppointmentRepository) Add(appointment model.Appointment) error {
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
		"end_time, theme, content, members) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.Exec(sqlStr, appointment.CreatorID, appointment.CreatorName, appointment.MeetingID, appointment.Day,
		appointment.StartTime, appointment.EndTime, appointment.Theme, appointment.Content,
		appointment.Members)
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
	if appointment.Members != "" {
		sqlStr = "update " + ar.userTable + " " +
			"set appointments=concat(appointments, '," + strconv.Itoa(int(id)) + "') " +
			"where id in (" + appointment.Members + ")"
		_, err = tx.Exec(sqlStr)
		if err != nil {
			tx.Rollback()
			return err
		}
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
	var users []model.User
	sql := ""
	if members != "" {
		sql = "select id, appointments from " + ar.userTable + " where id in (" + members + ")"
		if err := ar.mysqlConn.Select(&users, sql); err != nil {
			return err
		}
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
	if members != "" {
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
	}
	tx.Commit()
	return nil
}

func (ar *AppointmentRepository) Put(appointment model.Appointment, addMembers, deleteMembers string) error {
	/*
	 * 1. 将deleteMembers中的每个成员删除该会议
	 * 2. 将addMembers中的每个成员添加该会议
	 * 3. 更新会议
	 */
	if err := ar.Conn(); err != nil {
		return err
	}

	idStr := strconv.Itoa(appointment.ID)
	// 1. 将deleteMembers中的每个成员删除该会议
	sql := ""
	tx, err := ar.mysqlConn.Begin()
	if err != nil {
		return err
	}
	if deleteMembers != "" {
		sql = "select id, appointments from " + ar.userTable + " where id in (" + deleteMembers + ")"
		var users []model.User
		if err := ar.mysqlConn.Select(&users, sql); err != nil {
			return err
		}
		sql = "update " + ar.userTable + " set appointments=? where id=?"
		for _, user := range users {
			tmp := common.MemberStrToList(user.Appointments)
			index := common.StrIndexOf(tmp, idStr)
			if index != -1 {
				tmp = append(tmp[:index], tmp[index+1:]...)
				if _, err := tx.Exec(sql, common.MemberListToStr(tmp), user.ID); err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	// 2. 将addMembers中的每个成员添加该会议
	if addMembers != "" {
		sql = "update " + ar.userTable + " set appointments=concat(appointments, '," + idStr + "') where id in (" + addMembers + ")"
		if _, err = tx.Exec(sql); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 3. 更新会议
	//sql = "update " + ar.appointmentTable + " set day=?, start_time=?, " +
	//	"end_time=?, theme=?, content=?, members=? where id =?"
	sql = "update " + ar.appointmentTable + " set theme=?, content=?, members=? where id =?"
	if _, err = tx.Exec(sql, appointment.Theme, appointment.Content, appointment.Members,
		appointment.ID); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (ar *AppointmentRepository) PutState(id int, state string) error {
	if err := ar.Conn(); err != nil {
		return err
	}

	sql := "update " + ar.appointmentTable + " set state=? where id=?"
	stmt, err := ar.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(state, id)
	if err != nil {
		return err
	}

	num, _ := result.RowsAffected()
	if num == 0 {
		return errors.New("no exists")
	}
	return nil
}

func (ar *AppointmentRepository) SelectConflictAppointments(appointment model.Appointment, limit int, attrNames ...string) (appointments []model.Appointment, err error) {
	if err = ar.Conn(); err != nil {
		return
	}

	attrs := ""
	if len(attrNames) == 0 {
		attrs = "id"
	} else {
		for _, v := range attrNames {
			attrs += v + ","
		}
		attrs = strings.Trim(attrs, ",")
	}

	sql := "select " + attrs + " from " + ar.appointmentTable + " where meeting_id=? and day=? and state != ? and state != ? " +
		"and (((start_time >= ? and start_time < ?) " +
		"or (end_time > ? and end_time <= ?)) " +
		"or (start_time <= ? and end_time >= ?))"

	if limit != 0 {
		sql += " limit ?"
		err = ar.mysqlConn.Select(&appointments, sql, appointment.MeetingID, model.AppointmentCancel,
			model.AppointmentAdoptCancel, appointment.Day, appointment.StartTime,
			appointment.EndTime, appointment.StartTime, appointment.EndTime,
			appointment.StartTime, appointment.EndTime, limit)
	} else {
		err = ar.mysqlConn.Select(&appointments, sql, appointment.MeetingID, appointment.Day, appointment.StartTime,
			appointment.EndTime, appointment.StartTime, appointment.EndTime, appointment.StartTime, appointment.EndTime)
	}
	return
}

func (ar *AppointmentRepository) SelectAppointmentsIDByTime(appointment model.Appointment) (appointments []model.Appointment, err error) {
	if err = ar.Conn(); err != nil {
		return
	}
	sql := "select meeting_id from " + ar.appointmentTable + " where day=? and state != ? and state != ? " +
		"and (((start_time >= ? and start_time < ?) " +
		"or (end_time > ? and end_time <= ?)) " +
		"or (start_time <= ? and end_time >= ?))"
	err = ar.mysqlConn.Select(&appointments, sql, appointment.Day, model.AppointmentCancel,
		model.AppointmentAdoptCancel, appointment.StartTime,
		appointment.EndTime, appointment.StartTime, appointment.EndTime, appointment.StartTime, appointment.EndTime)
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

func (ar *AppointmentRepository) SelectCreator(day, startTime, meetingID string) (appointments []model.Appointment, err error) {
	if err = ar.Conn(); err != nil {
		return
	}

	sql := "select meeting_id, creator_name, day, start_time, end_time from " + ar.appointmentTable + " where day=? and start_time>=? and state != ? and state != ?"
	if meetingID != "" {
		sql += " and meeting_id in (" + meetingID + ")"
	}

	err = ar.mysqlConn.Select(&appointments, sql, day, startTime, model.AppointmentCancel, model.AppointmentAdoptCancel)
	return
}

func (ar *AppointmentRepository) SelectAppointmentsByCondition(conditionName, conditionVal string) (appointments []model.Appointment, err error) {
	if err = ar.Conn(); err != nil {
		return
	}

	now := time.Now()
	day := now.Format("20060102")
	startTime := now.Format("15") + ":00"
	sql := "select id, meeting_id, day, start_time, end_time, state, theme, content " +
		"from " + ar.appointmentTable + " where " + conditionName + " = ? and state != ? and state != ? and ((day = ? and start_time >= ?) or (day > ?)) order by day, start_time"
	err = ar.mysqlConn.Select(&appointments, sql, conditionVal, model.AppointmentCancel, model.AppointmentAdoptCancel, day, startTime, day)
	return
}

func (ar *AppointmentRepository) SelectAppointmentsByID(ids string) (appointments []model.Appointment, err error) {
	if ids == "" {
		return
	}
	if err = ar.Conn(); err != nil {
		return
	}

	now := time.Now()
	day := now.Format("20060102")
	startTime := now.Format("15") + ":00"
	sql := "select id, meeting_id, day, start_time, end_time, state, theme, content " +
		"from " + ar.appointmentTable + " where id in (" + ids + ") and state != ? and state != ? and ((day = ? and start_time >= ?) or (day > ?)) order by day, start_time"
	err = ar.mysqlConn.Select(&appointments, sql, model.AppointmentCancel, model.AppointmentAdoptCancel, day, startTime, day)
	return
}

func (ar *AppointmentRepository) SelectAppointmentsByPage(page, onePageCount int, state string) (appointments []model.Appointment, err error) {
	if err = ar.Conn(); err != nil {
		return
	}

	now := time.Now()
	day := now.Format("20060102")
	startTime := now.Format("15") + ":00"
	page -= 1
	startIndex := strconv.Itoa(page * onePageCount)
	sql := "select id, day, start_time, end_time, state, theme from " + ar.appointmentTable + " where " +
		"state = ? and ((day = ? and start_time >= ?) or (day > ?)) order by id desc limit " + startIndex + "," + strconv.Itoa(onePageCount)
	err = ar.mysqlConn.Select(&appointments, sql, state, day, startTime, day)
	return
}

func (ar *AppointmentRepository) SelectCountByState(state string) (count int, err error) {
	if err = ar.Conn(); err != nil {
		return
	}

	now := time.Now()
	day := now.Format("20060102")
	startTime := now.Format("15") + ":00"
	sql := "select count(*) from " + ar.appointmentTable + " where state=? and ((day = ? and start_time >= ?) or (day > ?))"
	err = ar.mysqlConn.QueryRow(sql, state, day, startTime, day).Scan(&count)
	return
}

func (ar *AppointmentRepository) SelectExpireAppointment(day string, endTime string) (appointments []model.Appointment, err error) {
	if err = ar.Conn(); err != nil {
		return
	}

	sql := "select id, creator_id, creator_name, meeting_id, day, start_time, end_time, state" +
		" from " + ar.appointmentTable + " where day < ? or (day = ? and end_time <= ?)"
	err = ar.mysqlConn.Select(&appointments, sql, day, day, endTime)
	return
}

func (ar *AppointmentRepository) TransferAppointment(appointment model.Appointment, members string) error {
	/*
	 * 1. 插入 日志表
	 * 2. 从会议表删除
	 */
	if err := ar.Conn(); err != nil {
		return err
	}

	tx, err := ar.mysqlConn.Begin()
	if err != nil {
		return err
	}
	// 1. 插入 日志表
	sql := "insert into record(creator_id, creator_name, meeting_id, day, start_time, end_time, state) " +
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
		sql = "select id, appointments from " + ar.userTable + " where id in (" + members + ")"
		if err := ar.mysqlConn.Select(&users, sql); err != nil {
			return err
		}
	}
	// 删除该会议室
	sql = "delete from " + ar.appointmentTable + " where id = ?"
	if _, err = tx.Exec(sql, appointment.ID); err != nil {
		tx.Rollback()
		return err
	}
	// 从所有成员中删除该会议
	if members != "" {
		idStr := strconv.Itoa(appointment.ID)
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
	}
	tx.Commit()
	return nil
}
