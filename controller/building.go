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
	"strings"
)

type IBuildingController interface {
	PostBuilding(c *gin.Context)
	DeleteBuilding(c *gin.Context)
	PutBuilding(c *gin.Context)
	GetBuildingByCampusAndPage(c *gin.Context)
	GetBuildingCountByCampus(c *gin.Context)
	GetBuildingByID(c *gin.Context)
}

func NewBuildingController() IBuildingController {
	buildingRepo := repository.NewBuildingRepository("building")
	campusRepo := repository.NewCampusRepository("campus")
	buildingService := service.NewBuildingService(buildingRepo, campusRepo)
	return &BuildingController{buildingService}
}

type BuildingController struct {
	BuildingService service.IBuildingService
}

func (bc *BuildingController) PostBuilding(c *gin.Context) {
	// 1. 解析请求数据
	building := model.Building{}
	if err := c.ShouldBindJSON(&building); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 访问数据库
	if err := bc.BuildingService.AddBuilding(building); err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, err.Error())
			return
		}
		if strings.Contains(err.Error(), "Duplicate") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "该建筑已存在")
			return
		}
		logger.Record("修改建筑错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

func (bc *BuildingController) DeleteBuilding(c *gin.Context) {
	// 1. 解析请求数据
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "the id must be a number")
		return
	}
	// 2. 访问数据库
	if err := bc.BuildingService.DeleteBuilding(id); err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "the campus no exists")
		} else {
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

func (bc *BuildingController) PutBuilding(c *gin.Context) {
	// 1. 解析请求数据
	building := model.Building{}
	if err := c.ShouldBindJSON(&building); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 访问数据库
	if err := bc.BuildingService.UpdateBuilding(building); err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, err.Error())
		} else {
			logger.Record("修改建筑错误", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

func (bc *BuildingController) GetBuildingCountByCampus(c *gin.Context) {
	result := map[string]interface{}{"count": 0}
	// 1. 解析获得campus_id值
	campusID, err := strconv.Atoi(c.Param("campus_id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "the campus_id must be a number")
		return
	}
	// 2. 数据库操作
	count, err := bc.BuildingService.GetBuildingCountByCampus(campusID)
	if err != nil {
		logger.Record("获取数据库出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	// 3. 返回结果
	result["count"] = count
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (bc *BuildingController) GetBuildingByCampusAndPage(c *gin.Context) {
	result := map[string]interface{}{
		"buildingList": []model.Campus{},
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
	campusID, err := strconv.Atoi(c.Query("campus_id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	// 2. 数据操作
	// 获取校区数量
	count, err := bc.BuildingService.GetBuildingCountByCampus(campusID)
	if err != nil {
		logger.Record("获取数据库出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	// 获取校区
	buildingList, err := bc.BuildingService.GetBuildingByCampusAndPage(page, onePageCount, campusID)
	if err != nil {
		fmt.Println(err)
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	result["count"] = count
	result["buildingList"] = buildingList
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (bc *BuildingController) GetBuildingByID(c *gin.Context) {
	result := map[string]interface{}{
		"building": model.Building{},
	}

	// 1. 解析请求数据
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "the id must be a number")
		return
	}
	// 2. 访问数据库
	building, err := bc.BuildingService.GetBuildingByID(id)
	if err != nil {
		fmt.Println("找不到id err -> ", err)
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "the campus no exists")
		} else {
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	result["building"] = building
	common.ResolveResult(c, true, e.SUCCESS, result)
}
