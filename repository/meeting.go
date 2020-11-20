package repository

import (
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"github.com/jmoiron/sqlx"
	"errors"
	"strconv"
)

type IMeetingRepository interface {
	Conn() error
	InsertMeeting(meeting model.Meeting) error
	DeleteMeeting(id int) error
	UpdateMeeting(meeting model.Meeting) error
	SelectMeetingByID(id int) (model.Meeting, error)
	SelectMeetingsByBuilding(buildingID int, pageAndOnePageCount ...int) ([]model.Meeting, error)
	SelectMeetingCountCountByBuilding(buildingID int) (int, error)
	SelectAllMeetingTypes() []string
	SelectAllScaleTypes() []string
}

type MeetingRepository struct {
	table     string
	mysqlConn *sqlx.DB
}

func NewMeetingRepository(table string) IMeetingRepository {
	return &MeetingRepository{table, database.MysqlDB}
}

func (mr *MeetingRepository) Conn() error {
	if mr.mysqlConn == nil {
		if err := database.InitMysql(); err != nil {
			return err
		}
		mr.mysqlConn = database.MysqlDB
	}
	if mr.table == "" {
		mr.table = "meeting"
	}
	return nil
}

func (mr *MeetingRepository) InsertMeeting(meeting model.Meeting) error {
	if err := mr.Conn(); err != nil {
		return err
	}

	sqlStr := "insert into " + mr.table + "(meeting_name, building_id, layer, meeting_type, scale, room_number) values(?, ?, ?, ?, ?, ?)"

	stmt, err := mr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(meeting.MeetingName, meeting.BuildingID, meeting.Layer, meeting.MeetingType, meeting.Scale, meeting.RoomNumber)
	return err
}

func (mr *MeetingRepository) DeleteMeeting(id int) error {
	if err := mr.Conn(); err != nil {
		return err
	}

	sqlStr := "delete from " + mr.table + " where id = ?"

	stmt, err := mr.mysqlConn.Prepare(sqlStr)
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
		return errors.New("the meeting is no exists")
	}
	return nil
}

func (mr *MeetingRepository) UpdateMeeting(meeting model.Meeting) error {
	if err := mr.Conn(); err != nil {
		return err
	}

	sqlStr := "update " + mr.table + " set meeting_name=?, building_id=?, layer=?, meeting_type=?, scale=?, room_number=? where id = ?"

	stmt, err := mr.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(meeting.MeetingName, meeting.BuildingID, meeting.Layer, meeting.MeetingType, meeting.Scale, meeting.RoomNumber, meeting.ID)
	if err != nil {
		return err
	}

	num, _ := result.RowsAffected()
	if num == 0 {
		return errors.New("the meeting no exists")
	}
	return nil
}

func (mr *MeetingRepository) SelectMeetingByID(id int) (meeting model.Meeting, err error) {
	if err = mr.Conn(); err != nil {
		return
	}

	sqlStr := "select * from " + mr.table + " where id = ? limit 1"

	err = mr.mysqlConn.Get(&meeting, sqlStr, id)
	return
}

/* SelectMeetingsByBuilding 查询建筑中 全部/分页 会议室
 * 1. buildingID
 * 2. pageAndOnePageCount[0]: page 第几页, 从0开始
 * 3. pageAndOnePageCount[1]: onePageCount 一页多少个
 */
func (mr *MeetingRepository) SelectMeetingsByBuilding(buildingID int, pageAndOnePageCount ...int) (meetings []model.Meeting, err error) {
	if err = mr.Conn(); err != nil {
		return
	}

	/* 是否有页和页码
	 * 1. 有 -> 分页查询
	 * 2. 没有 -> 查询全部
	 */
	sqlStr := "select id, meeting_name, layer, meeting_type, scale, room_number  from " + mr.table + " where building_id = ?"
	if len(pageAndOnePageCount) >= 2 {
		page := pageAndOnePageCount[0]
		onePageCount := pageAndOnePageCount[1]
		startIndex := strconv.Itoa(page * onePageCount)
		sqlStr += " limit " + startIndex + ", " + strconv.Itoa(onePageCount)
	}
	err = mr.mysqlConn.Select(&meetings, sqlStr, buildingID)
	return
}

func (mr *MeetingRepository) SelectMeetingCountCountByBuilding(buildingID int) (count int, err error) {
	if err = mr.Conn(); err != nil {
		return
	}

	sqlStr := "select count(*) from " + mr.table + " where building_id = ?"
	err = mr.mysqlConn.QueryRow(sqlStr, buildingID).Scan(&count)
	return
}

func (mr *MeetingRepository) SelectAllMeetingTypes() []string {
	return model.GetMeetingTypeList()
}

func (mr *MeetingRepository) SelectAllScaleTypes() []string {
	return model.GetScaleTypeList()
}

