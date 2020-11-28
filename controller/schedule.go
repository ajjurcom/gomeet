package controller

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/e"
	"com/mittacy/gomeet/logger"
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"com/mittacy/gomeet/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type IScheduleController interface {
	GetOptions(c *gin.Context)
	UpdateOptions(c *gin.Context)
}

func NewScheduleController() IScheduleController {
	meetingRepo := repository.NewMeetingRepository("meeting")
	buildingRepo := repository.NewBuildingRepository("building")
	campusRepo := repository.NewCampusRepository("campus")

	meetingService := service.NewMeetingService(meetingRepo)
	buildingService := service.NewBuildingService(buildingRepo)
	campusService := service.NewCampusService(campusRepo)
	return &ScheduleController{meetingService, buildingService, campusService}
}

type ScheduleController struct {
	MeetingService service.IMeetingService
	BuildingService service.IBuildingService
	CampusService service.ICampusService
}

func (sc *ScheduleController) GetOptions(c *gin.Context) {
	result := map[string]interface{}{
		"campusList": []model.Campus{},
		"buildingList": []model.Building{},
		"meetingList": []model.Meeting{},
		"meetingTypes": []string{},
		"meetingScales": []string{},
	}
	// 1. 获取campusList
	campus, err := sc.CampusService.GetAllCampus()
	if err != nil {
		logger.Record("获取所有校区出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	result["campusList"] = campus
	// 2. 获取campusList[0]的buildingList
	if len(campus) == 0 {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "该校区中没有建筑")
		return
	}
	buildings, err := sc.BuildingService.GetAllBuildingsByCampus(campus[0].ID)
	if err != nil {
		logger.Record("获取校区的全部建筑出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	result["buildingList"] = buildings
	// 3. 获取buildingList[0]的全部会议室meetingList
	if len(buildings) == 0 {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "该建筑中没有会议室")
		return
	}
	meetings, err := sc.MeetingService.GetAllMeetingByBuilding(buildings[0].ID)
	if err != nil {
		logger.Record("获取校区的全部建筑出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}

	result["meetingList"] = meetings
	result["meetingTypes"] = sc.MeetingService.GetAllMeetingTypes()
	result["meetingScales"] = sc.MeetingService.GetAllScaleTypes()
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (sc *ScheduleController) UpdateOptions(c *gin.Context) {
	/* 1. 获取所有请求参数
	 * 2. 获取改变事件参数
	 * 3. 根据不同参数，获取不同内容
	 * 4. 操作数据库
	 * 5. 返回结果
	 */
	// 1. 变量声明
	var (
		campusID, buildingID, layer int
		meetingTypes, meetingScales []string
		way string
	)
	result := map[string]interface{}{
		"campusList": []model.Campus{},
		"buildingList": []model.Building{},
		"meetingList": []model.Meeting{},
	}
	// 2. 解析请求
	var err1, err2, err3 error
	campusID, err1 = strconv.Atoi(c.Query("campusID"))
	buildingID, err2 = strconv.Atoi(c.Query("buildingID"))
	layer, err3 = strconv.Atoi(c.Query("layer"))
	meetingTypes, _ = c.GetQueryArray("meetingTypes[]")
	meetingScales, _ = c.GetQueryArray("meetingScales[]")
	way = c.Query("way")
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("err: ", err1, err2, err3)
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "校区、建筑、楼层ID需为数字")
		return
	}
	fmt.Println("campusID: ", campusID)
	fmt.Println("buildingID: ", buildingID)
	fmt.Println("layer: ", layer)
	fmt.Println("meetingTypes: ", meetingTypes)
	fmt.Println("meetingScales: ", meetingScales)
	fmt.Println("way: ", way)
	// 3. 根据不同参数修改不同变量
	if way == "campus" {	// 更新建筑
		buildings, err := sc.BuildingService.GetAllBuildingsByCampus(campusID)
		if err != nil {
			logger.Record("获取校区的全部建筑出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
			return
		}
		result["buildingList"] = buildings
		if len(buildings) == 0 {
			common.ResolveResult(c, true, e.SUCCESS, result, "该校区中没有建筑")
			return
		}
		buildingID = buildings[0].ID
	}
	meetings, err := sc.MeetingService.GetAllMeetingsByParams(buildingID, layer, meetingTypes, meetingScales)
	if err != nil {
		logger.Record("获取校区的全部建筑出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	result["meetingList"] = meetings

	common.ResolveResult(c, true, e.SUCCESS, result)
}
