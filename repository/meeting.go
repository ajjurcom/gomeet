package repository

import (
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"errors"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type IMeetingRepository interface {
	Conn() error
	InsertMeeting(meeting model.Meeting) error
	DeleteMeeting(id int) error
	UpdateMeeting(meeting model.Meeting) error
	SelectMeetingByID(id int) (model.Meeting, error)
	SelectMeetingsByBuilding(buildingID int, pageAndOnePageCount ...int) ([]model.Meeting, error)
	SearchMeetingsByKeyword(page, onePageCount int, keyword string) ([]model.Meeting, error)
	SelectAllMeetingsByParams(buildingID int, layer int, meetingType []string, scales []string) ([]model.Meeting, error)
	SelectMeetingCount(attrName, attrVal string, isEqual bool) (int, error)
	SelectAllMeetingTypes() []string
	SelectAllScaleTypes() []string
	SelectMeetingByInfo(meetingsID, campusID, meetingType, meetingScale string) (model.Meeting, error)
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
	sqlStr := "select id, meeting_name, layer, meeting_type, scale, room_number  from " + mr.table + " where building_id = ? order by id desc"
	if len(pageAndOnePageCount) >= 2 {
		page := pageAndOnePageCount[0] - 1
		onePageCount := pageAndOnePageCount[1]
		startIndex := strconv.Itoa(page * onePageCount)
		sqlStr += " limit " + startIndex + ", " + strconv.Itoa(onePageCount)
	}
	err = mr.mysqlConn.Select(&meetings, sqlStr, buildingID)
	return
}

func (mr *MeetingRepository) SearchMeetingsByKeyword(page, onePageCount int, keyword string) (meetings []model.Meeting, err error) {
	if err = mr.Conn(); err != nil {
		return
	}

	page -= 1
	startIndex := strconv.Itoa(page * onePageCount)
	sqlStr := "select id, meeting_name, layer, meeting_type, scale, room_number from " + mr.table + " where meeting_name like '%" + keyword + "%' limit " + startIndex + ", " + strconv.Itoa(onePageCount)
	err = mr.mysqlConn.Select(&meetings, sqlStr)
	return
}

func (mr *MeetingRepository) SelectMeetingCount(attrName, attrVal string, isEqual bool) (count int, err error) {
	if err = mr.Conn(); err != nil {
		return
	}

	sql := "select count(*) from " + mr.table + " where " + attrName
	if isEqual {
		sql += " = ?"
		err = mr.mysqlConn.QueryRow(sql, attrVal).Scan(&count)
		return
	}
	sql += " like ('%" + attrVal + "%')"
	err = mr.mysqlConn.QueryRow(sql).Scan(&count)
	return

}

func (mr *MeetingRepository) SelectAllMeetingTypes() []string {
	return model.GetMeetingTypeList()
}

func (mr *MeetingRepository) SelectAllScaleTypes() []string {
	return model.GetScaleTypeList()
}

func (mr *MeetingRepository) SelectAllMeetingsByParams(buildingID int, layer int, meetingType []string, scales []string) (meetings []model.Meeting, err error) {
	if err = mr.Conn(); err != nil {
		return
	}

	/* 是否有页和页码
	 * 1. 有 -> 分页查询
	 * 2. 没有 -> 查询全部
	 */
	//sqlStr := "select id, meeting_name, layer, meeting_type, scale, room_number  from " + mr.table + " where building_id = ?"
	sqlStr := "select id, meeting_name, layer, meeting_type, scale, room_number  from " + mr.table + " where building_id = " + strconv.Itoa(buildingID)
	if layer > 0 {
		sqlStr += " and layer = " + strconv.Itoa(layer)
	}
	if len(meetingType) > 0 {
		str := "'" + meetingType[0] + "'"
		for i := 1; i < len(meetingType); i++ {
			str += ", '" + meetingType[i] + "'"
		}
		sqlStr += " and meeting_type in (" + str + ")"
	}
	if len(scales) > 0 {
		str := "'" + scales[0] + "'"
		for i := 1; i < len(scales); i++ {
			str += ", '" + scales[i] + "'"
		}
		sqlStr += " and scale in (" + str + ")"
	}
	//err = mr.mysqlConn.Select(&meetings, sqlStr, buildingID)
	err = mr.mysqlConn.Select(&meetings, sqlStr)
	return
}

func (mr *MeetingRepository) SelectMeetingByInfo(meetingsID, campusID, meetingType, meetingScale string) (meeting model.Meeting, err error) {
	if err = mr.Conn(); err != nil {
		return
	}

	sql := "select meeting.id from building, meeting where "
	if meetingsID != "" {
		sql += "meeting.id not in(" + meetingsID +") and "
	}
	sql += "building.campus_id=? and meeting.building_id=building.id and meeting.meeting_type=? and meeting.scale=? limit 1"
	err = mr.mysqlConn.Get(&meeting, sql, campusID, meetingType, meetingScale)
	return
}

