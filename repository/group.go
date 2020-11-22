package repository

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/model"
	"errors"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type IGroupRepository interface {
	Conn() error
	Add(group model.Group) error
	Delete(id int) error
	PutName(group model.Group) error
	PutMember(newGroup model.Group) error
	SelectGroupByID(id int) (group model.Group, err error)
	SelectGroupsByCreator(creator int, pageAndOnePageCount ...int) ([]model.Group, error)
	SelectGroupCountByCreator(creator int) (int, error)
}

type GroupRepository struct {
	groupTable string
	userTable string
	mysqlConn *sqlx.DB
}

func NewGroupRepository(groupTable, userTable string) IGroupRepository {
	return &GroupRepository{groupTable, userTable, database.MysqlDB}
}

func (gr *GroupRepository) Conn() error {
	if gr.mysqlConn == nil {
		if err := database.InitMysql(); err != nil {
			return err
		}
		gr.mysqlConn = database.MysqlDB
	}
	if gr.groupTable == "" {
		gr.groupTable = "user_group"
	}
	if gr.userTable == "" {
		gr.userTable = "user"
	}
	return nil
}

func (gr *GroupRepository) Add(group model.Group) error {
	/*
	 * 1. 将memberList转为字符串保存到member_list
	 * 2. 获取新组的ID，添加到member_list每个user的group_list中
	 */
	if err := gr.Conn(); err != nil {
		return err
	}

	tx, err := gr.mysqlConn.Begin()
	if err != nil {
		return err
	}

	var id int64 = -1
	// 1. 将memberList转为字符串保存到member_list
	sqlStr := "insert into " + gr.groupTable + "(creator, group_name, member_list) values (?, ?, ?)"
	result, err := tx.Exec(sqlStr, group.Creator, group.GroupName, group.MemberList)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err = result.LastInsertId()
	if err != nil || id == -1 {
		tx.Rollback()
		return err
	}
	// 2. 将新组的ID添加到member_list每个user的group_list中
	sqlStr = "update " + gr.userTable + " set group_list=concat(group_list, '," + strconv.Itoa(int(id)) + "') where id in (" + group.MemberList + ")"
	result, err = tx.Exec(sqlStr)
	if err != nil {
		tx.Rollback()
		return err
	}
	aff, _ := result.RowsAffected()
	if int(aff) < len(common.MemberStrToList(group.MemberList)) {
		tx.Rollback()
		return errors.New("some ID don't exist")
	}
	tx.Commit()
	return nil
}

func (gr *GroupRepository) Delete(id int) error {
	/*
	 * 1. 查询用户组中的所有成员
	 * 2. 查询所有成员的group_list
	 * 3. 删除该用户组
	 * 4. 从所有成员中删除该用户组
	 */
	if err := gr.Conn(); err != nil {
		return err
	}
	// 1. 查询用户组中的所有成员
	memberList := ""
	sqlStr := "select member_list from user_group where id = ?"
	err := gr.mysqlConn.QueryRow(sqlStr, id).Scan(&memberList)
	if err != nil {
		return err
	}
	if len(common.MemberStrToList(memberList)) == 0 {
		return nil
	}
	// 2. 查询所有成员的group_list
	var userList []model.User
	sqlStr = "select id, group_list from user where id in (" + memberList + ")"
	if err = gr.mysqlConn.Select(&userList, sqlStr); err != nil {
		return err
	}
	// 3. 删除用户组
	tx, err := gr.mysqlConn.Begin()
	if err != nil {
		return err
	}
	sqlStr = "delete from " + gr.groupTable + " where id=?"
	if _, err = tx.Exec(sqlStr, id); err != nil {
		tx.Rollback()
		return err
	}
	// 4. 从所有成员中删除该组id
	sqlStr = "update " + gr.userTable + " set group_list=? where id=?"
	idStr := strconv.Itoa(id)
	for _, user := range userList {
		tmpList := common.MemberStrToList(user.GroupList)
		index := common.StrIndexOf(tmpList, idStr)
		if index != -1 {
			tmpList = append(tmpList[:index], tmpList[index+1:]...)
		}
		_, err := tx.Exec(sqlStr, common.MemberListToStr(tmpList), user.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (gr *GroupRepository) PutName(group model.Group) error {
	if err := gr.Conn(); err != nil {
		return err
	}

	sqlStr := "update " + gr.groupTable + " set group_name=? where id=?"
	_, err := gr.mysqlConn.Exec(sqlStr, group.GroupName, group.ID)
	return err
}

func (gr *GroupRepository) PutMember(newGroup model.Group) error {
	/*
	 * 1. 获取旧的member_list
	 * 2. 计算新旧list的差异: deleteList, newList
	 * 3. 查询deleteList所有成员信息
	 * 5. 更新member_list到group
	 * 6. 将deleteUserList中的每个成员删除该ID
	 * 7. 向addUserList中的每个成员添加该ID
	 */
	if err := gr.Conn(); err != nil {
		return err
	}

	// 1. 获取旧的member_list
	oldGroup, err := gr.SelectGroupByID(newGroup.ID)
	if err != nil {
		return err
	}

	// 2. 计算新旧list的差异: deleteList, addList
	deleteList, addList := common.DiffMember(common.MemberStrToList(oldGroup.MemberList), common.MemberStrToList(newGroup.MemberList))

	// 3. 查询deleteList所有成员信息
	var deleteUserList []model.User
	if len(deleteList) > 0 {
		sqlStr := "select id, group_list from user where id in (" + common.MemberListToStr(deleteList) + ")"
		if err = gr.mysqlConn.Select(&deleteUserList, sqlStr); err != nil {
			return err
		}
	}

	// 5. 更新member_list到group
	tx, err := gr.mysqlConn.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	sqlStr := "update " + gr.groupTable + " set member_list=? where id=?"
	if _, err := tx.Exec(sqlStr, newGroup.MemberList, oldGroup.ID); err != nil {
		tx.Rollback()
		return err
	}

	// 6. 将deleteUserList中的每个成员删除该ID
	idStr := strconv.Itoa(newGroup.ID)
	if len(deleteUserList) > 0 {
		sqlStr = "update " + gr.userTable + " set group_list=? where id=?"
		for _, user := range deleteUserList {
			tmpList := common.MemberStrToList(user.GroupList)
			index := common.StrIndexOf(tmpList, idStr)
			if index != -1 {
				tmpList = append(tmpList[:index], tmpList[index+1:]...)
			}
			_, err := tx.Exec(sqlStr, common.MemberListToStr(tmpList), user.ID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	// 7. 向addUserList中的每个成员添加该ID
	if len(addList) > 0 {
		sqlStr = "update " + gr.userTable + " set group_list=concat(group_list, '," + idStr + "') where id in (" + common.MemberListToStr(addList) + ")"
		result, err := tx.Exec(sqlStr)
		if err != nil {
			tx.Rollback()
			return err
		}
		aff, _ := result.RowsAffected()
		if int(aff) < len(addList) {
			tx.Rollback()
			return errors.New("some ID don't exist")
		}
	}
	tx.Commit()
	return nil
}

func (gr *GroupRepository) SelectGroupByID(id int) (group model.Group, err error) {
	if err = gr.Conn(); err != nil {
		return
	}
	sqlStr := "select * from user_group where id = ?"
	err = gr.mysqlConn.Get(&group, sqlStr, id)
	return
}

/* SelectMeetingsByBuilding 查询建筑中 全部/分页 会议室
 * 1. creator 创建者
 * 2. pageAndOnePageCount[0]: page 第几页, 从0开始
 * 3. pageAndOnePageCount[1]: onePageCount 一页多少个
 */
func (gr *GroupRepository) SelectGroupsByCreator(creator int, pageAndOnePageCount ...int) (groups []model.Group, err error) {
	if err = gr.Conn(); err != nil {
		return
	}

	/* 是否有页和页码
	 * 1. 有 -> 分页查询
	 * 2. 没有 -> 查询全部
	 */
	sqlStr := "select id, group_name  from " + gr.groupTable + " where creator = ?"
	if len(pageAndOnePageCount) >= 2 {
		page := pageAndOnePageCount[0] - 1
		onePageCount := pageAndOnePageCount[1]
		startIndex := strconv.Itoa(page * onePageCount)
		sqlStr += " limit " + startIndex + ", " + strconv.Itoa(onePageCount)
	}
	err = gr.mysqlConn.Select(&groups, sqlStr, creator)
	return
}

func (gr *GroupRepository) SelectGroupCountByCreator(creator int) (count int, err error) {
	if err = gr.Conn(); err != nil {
		return
	}

	sqlStr := "select count(*) from " + gr.groupTable + " where creator = ?"
	err = gr.mysqlConn.QueryRow(sqlStr, creator).Scan(&count)
	return
}
