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

type IMeetingController interface {
	PostMeeting(c *gin.Context)
	DeleteMeeting(c *gin.Context)
	PutMeeting(c *gin.Context)
	GetMeetingByID(c *gin.Context)
	GetMeetingsByPage(c *gin.Context)
	GetMeetingOptions(c *gin.Context)
	SearchMeetings(c *gin.Context)
}

func NewMeetingController() IMeetingController {
	meetingRepo := repository.NewMeetingRepository("meeting")
	buildingRepo := repository.NewBuildingRepository("building")
	campusRepo := repository.NewCampusRepository("campus")

	meetingService := service.NewMeetingService(meetingRepo)
	buildingService := service.NewBuildingService(buildingRepo)
	campusService := service.NewCampusService(campusRepo)
	return &MeetingController{meetingService, buildingService, campusService}
}

type MeetingController struct {
	MeetingService service.IMeetingService
	BuildingService service.IBuildingService
	CampusService service.ICampusService
}

// PostMeeting 添加会议室
func (mc *MeetingController) PostMeeting(c *gin.Context) {
	// 1. 解析请求数据
	meeting := model.Meeting{}
	if err := c.ShouldBindJSON(&meeting); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	/* 访问数据库
	 * 1. 保证会议室所选类型存在
	 * 2. 保证会议室大小类型存在
	 * 3. 保证会议室所属建筑存在
	 # 4. 保证会议室所在楼层<=建筑楼高
	 * 5. 添加会议室
	 */
	if !model.IsMeetingType(meeting.MeetingType) {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "所选会议室类型不存在")
		return
	}

	if !model.IsScaleType(meeting.Scale) {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "所选会议室容量大小不存在")
		return
	}

	buildingLayer, err := mc.BuildingService.GetBuildingLayer(meeting.BuildingID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "所选建筑不存在")
			return
		}
		logger.Record("添加会议室时，查询建筑是否存在错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	if meeting.Layer > buildingLayer {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "选择的会议室所在楼层高于建筑楼层")
		return
	}

	if err := mc.MeetingService.AddMeeting(meeting); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议室已存在")
			return
		}
		logger.Record("添加会议室错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// DeleteMeeting 删除会议室
func (mc *MeetingController) DeleteMeeting(c *gin.Context) {
	// 1. 解析请求数据
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议室ID必须为数字")
		return
	}
	// 2. 访问数据库
	if err := mc.MeetingService.DeleteMeeting(id); err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议室不存在")
		} else {
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// PutMeeting 修改会议室
func (mc *MeetingController) PutMeeting(c *gin.Context) {
	// 1. 解析请求数据
	meeting := model.Meeting{}
	if err := c.ShouldBindJSON(&meeting); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	/* 访问数据库
	* 0. 保证ID
	* 1. 保证会议室所选类型存在
	* 2. 保证会议室大小类型存在
	* 3. 保证会议室所属建筑存在
	# 4. 保证会议室所在楼层<=建筑楼高
	* 5. 添加会议室
	*/
	if meeting.ID == 0 {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议室ID不能缺省")
		return
	}
	if !model.IsMeetingType(meeting.MeetingType) {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "所选会议室类型不存在")
		return
	}

	if !model.IsScaleType(meeting.Scale) {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "所选会议室容量大小不存在")
		return
	}

	buildingLayer, err := mc.BuildingService.GetBuildingLayer(meeting.BuildingID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "所选建筑不存在")
			return
		}
		logger.Record("添加会议室时，查询建筑是否存在错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	if meeting.Layer > buildingLayer {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "选择的会议室所在楼层高于建筑楼层")
		return
	}

	if err := mc.MeetingService.UpdateMeeting(meeting); err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, err.Error())
		} else {
			logger.Record("修改会议室错误", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// GetMeetingByID 查询会议室详情
func (mc *MeetingController) GetMeetingByID(c *gin.Context) {
	result := map[string]interface{}{
		"meeting": model.Meeting{},
		"building": model.Building{},
	}

	// 1. 解析请求数据
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "the id must be a number")
		return
	}
	// 2. 访问数据库
	meeting, err := mc.MeetingService.GetMeetingByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "the meeting no exists")
		} else {
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	building, err := mc.BuildingService.GetBuildingByID(meeting.BuildingID)
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	// 3. 返回结果
	result["meeting"] = meeting
	result["building"] = building
	common.ResolveResult(c, true, e.SUCCESS, result)
}

// GetMeetingByPage 翻页查询建筑中的会议室
func (mc *MeetingController) GetMeetingsByPage(c *gin.Context) {
	result := map[string]interface{}{
		"meetingList": []model.Meeting{},
		"count": 0,
	}
	// 1. 解析请求
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	onePageCount, err := strconv.Atoi(c.Param("onePageCount"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	buildingID, err := strconv.Atoi(c.Query("building_id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	// 2. 数据操作
	// 获取建筑会议室总数量
	count, err := mc.MeetingService.GetMeetingCountByBuilding(buildingID)
	if err != nil {
		logger.Record("获取建筑会议室数量出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	// 获取会议室
	meetingList, err := mc.MeetingService.GetMeetingsByPage(page, onePageCount, buildingID)
	if err != nil && err != sql.ErrNoRows {
		logger.Record("获取建筑的会议室错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	result["count"] = count
	result["meetingList"] = meetingList
	common.ResolveResult(c, true, e.SUCCESS, result)
}

// GetOptions 获取会议室属性可选性
func (mc *MeetingController) GetMeetingOptions(c *gin.Context) {
	result := map[string]interface{}{
		"meetingTypes": []string{},
		"meetingScales": []string{},
		"campusList": []model.Campus{},
	}
	var err error
	result["campusList"], err = mc.CampusService.GetAllCampus()
	if err != nil {
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}

	result["meetingTypes"] = mc.MeetingService.GetAllMeetingTypes()
	result["meetingScales"] = mc.MeetingService.GetAllScaleTypes()
	common.ResolveResult(c, true, e.SUCCESS, result)
}

// SearchMeeting 查找会议室，分页返回
func (mc *MeetingController) SearchMeetings(c *gin.Context) {
	result := map[string]interface{}{
		"meetingList": []model.Meeting{},
		"count": 0,
	}
	// 1. 解析请求
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	onePageCount, err := strconv.Atoi(c.Param("onePageCount"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	keyword := c.Query("keyword")
	if strings.Trim(keyword, " ") == "" {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "keyword参数不能为空")
		return
	}
	// 2. 数据操作
	// 获取建筑会议室总数量
	count, err := mc.MeetingService.GetMeetingCountByKeyword(keyword)
	if err != nil {
		logger.Record("获取建筑会议室数量出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	// 获取会议室
	meetingList, err := mc.MeetingService.GetMeetingsByKeyword(page, onePageCount, keyword)
	if err != nil && err != sql.ErrNoRows {
		logger.Record("获取建筑的会议室错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	result["count"] = count
	result["meetingList"] = meetingList
	common.ResolveResult(c, true, e.SUCCESS, result)
}
