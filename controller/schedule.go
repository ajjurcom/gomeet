package controller

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/e"
	"com/mittacy/gomeet/logger"
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"com/mittacy/gomeet/service"
	"github.com/gin-gonic/gin"
)

type IScheduleController interface {
	GetOptions(c *gin.Context)
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