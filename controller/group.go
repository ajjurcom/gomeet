package controller

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/e"
	"com/mittacy/gomeet/logger"
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"com/mittacy/gomeet/service"
	"database/sql"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type IGroupController interface {
	Post(c *gin.Context)
	Delete(c *gin.Context)
	PutName(c *gin.Context)
	PutMember(c *gin.Context)
	GetGroupsByPage(c *gin.Context)
	GetAllGroups(c *gin.Context)
}

func NewGroupController() IGroupController {
	repo := repository.NewGroupRepository("user_group", "user")
	groupService := service.NewGroupService(repo)
	return &GroupController{groupService}
}

type GroupController struct {
	GroupService service.IGroupService
}

// 新增用户组
func (gc *GroupController) Post(c *gin.Context) {
	// 1. 解析请求
	group := model.Group{}
	if err := c.ShouldBindJSON(&group); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 操作数据库
	if err := gc.GroupService.CreateGroup(group); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "用户组已存在")
		} else if strings.Contains(err.Error(), "don't exist") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "用户组成员部分不存在")
		} else {
			logger.Record("新建成员组出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// 删除用户组
func (gc *GroupController) Delete(c *gin.Context) {
	// 1. 解析请求数据
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "用户ID必须为数字")
		return
	}
	// 2. 操作数据库
	if err := gc.GroupService.DeleteGroup(id); err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "用户组不存在")
		} else {
			logger.Record("删除成员组出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// 修改用户组名字
func (gc *GroupController) PutName(c *gin.Context) {
	// 1. 解析请求
	group := model.Group{}
	if err := c.ShouldBindJSON(&group); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 操作数据库
	if err := gc.GroupService.PutName(group); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "用户组已存在")
		} else {
			logger.Record("修改成员组名字出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// 修改用户组成员
func (gc *GroupController) PutMember(c *gin.Context) {
	// 1. 解析请求
	group := model.Group{}
	if err := c.ShouldBindJSON(&group); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 操作数据库
	if err := gc.GroupService.PutMember(group); err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "部分用户不存在")
		} else {
			logger.Record("修改成员组成员出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// 分页获取创建者的分组
func (gc *GroupController) GetGroupsByPage(c *gin.Context) {
	result := map[string]interface{}{
		"groupList": []model.Group{},
		"count": 0,
	}
	// 1. 解析请求
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "page参数必须为数字")
		return
	}
	onePageCount, err := strconv.Atoi(c.Param("onePageCount"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "onePageCount参数必须为数字")
		return
	}
	creator, err := strconv.Atoi(c.Query("creator"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "creator参数必须为数字")
		return
	}
	// 2. 数据操作
	// 2.1 获取创建者的所有组数量
	count, err := gc.GroupService.GetGroupCountByCreator(creator)
	if err != nil {
		logger.Record("获取用户组数量出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	// 2.2 分页获取组
	groupList, err := gc.GroupService.GetGroupsByCreatorAndPage(creator, page, onePageCount)
	if err != nil && err != sql.ErrNoRows {
		logger.Record("获取用户组错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}

	result["count"] = count
	result["groupList"] = groupList
	common.ResolveResult(c, true, e.SUCCESS, result)
}

// 获取creator的所有分组
func (gc *GroupController) GetAllGroups(c *gin.Context) {
	result := map[string]interface{}{
	}
	creator, err := strconv.Atoi(c.Param("creator"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "creator参数必须为数字")
		return
	}
	// 操作数据库
	groupList, err := gc.GroupService.GetAllGroupsByCreator(creator)
	if err != nil && err != sql.ErrNoRows {
		logger.Record("获取用户组错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	result["groupList"] = groupList
	common.ResolveResult(c, true, e.SUCCESS, result)
}
