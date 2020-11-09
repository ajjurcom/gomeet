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
	GetBuildingsByPage(c *gin.Context)
	GetBuildingByID(c *gin.Context)
	GetAllBuildingsByCampus(c *gin.Context)
	GetBuildingLayer(c *gin.Context)
}

func NewBuildingController() IBuildingController {
	buildingRepo := repository.NewBuildingRepository("building")
	campusRepo := repository.NewCampusRepository("campus")
	buildingService := service.NewBuildingService(buildingRepo)
	campusService := service.NewCampusService(campusRepo)
	return &BuildingController{buildingService, campusService}
}

type BuildingController struct {
	BuildingService service.IBuildingService
	CampusService service.ICampusService
}

// PostBuilding 添加建筑
func (bc *BuildingController) PostBuilding(c *gin.Context) {
	// 解析请求数据
	building := model.Building{}
	if err := c.ShouldBindJSON(&building); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	/* 访问数据库
	 * 1. 保证建筑所属校区存在
	 * 2. 添加建筑
	 */
	isExists, err := bc.CampusService.IsCampusExists(building.CampusID)
	if err != nil {
		logger.Record("添加建筑时，查询校区是否存在错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	if !isExists {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "所选校区不存在")
		return
	}

	if err := bc.BuildingService.AddBuilding(building); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "该建筑已存在")
			return
		}
		logger.Record("添加建筑错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// DeleteBuilding 删除建筑
func (bc *BuildingController) DeleteBuilding(c *gin.Context) {
	// 1. 解析请求数据
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "建筑ID必须为数字")
		return
	}
	// 2. 访问数据库
	if err := bc.BuildingService.DeleteBuilding(id); err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "建筑不存在")
		} else {
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// PutBuilding 修改建筑信息
func (bc *BuildingController) PutBuilding(c *gin.Context) {
	// 1. 解析请求数据
	building := model.Building{}
	if err := c.ShouldBindJSON(&building); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	/* 访问数据库
	 * 1. 保证建筑所属校区存在
	 * 2. 添加建筑
	 */
	isExists, err := bc.CampusService.IsCampusExists(building.CampusID)
	if err != nil {
		logger.Record("添加建筑时，查询校区是否存在错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	if !isExists {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "所选校区不存在")
		return
	}

	if err := bc.BuildingService.UpdateBuilding(building); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "新建筑名已存在")
			return
		}
		logger.Record("更新建筑错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// GetBuildingsByPage 翻页获取校区的建筑
func (bc *BuildingController) GetBuildingsByPage(c *gin.Context) {
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
	buildingList, err := bc.BuildingService.GetBuildingsByPage(page, onePageCount, campusID)
	if err != nil {
		fmt.Println(err)
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	result["count"] = count
	result["buildingList"] = buildingList
	common.ResolveResult(c, true, e.SUCCESS, result)
}

// GetBuildingByID 根据ID获取建筑详细信息
func (bc *BuildingController) GetBuildingByID(c *gin.Context) {
	result := map[string]interface{}{
		"building": model.Building{},
	}

	// 1. 解析请求数据
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "建筑ID必须为数字")
		return
	}
	// 2. 访问数据库
	building, err := bc.BuildingService.GetBuildingByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "建筑不存在")
			return
		}
		logger.Record("根据ID获取建筑详细信息出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 3. 返回结果
	result["building"] = building
	common.ResolveResult(c, true, e.SUCCESS, result)
}

// GetAllBuildingsByCampus 获取校区的全部建筑
func (bc *BuildingController) GetAllBuildingsByCampus(c *gin.Context) {
	result := map[string]interface{}{
		"buildings": []model.Building{},
	}
	// 1. 解析请求数据
	campusID, err := strconv.Atoi(c.Param("campus_id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "校区ID必须为数字")
		return
	}
	// 2. 访问数据库
	buildings, err := bc.BuildingService.GetAllBuildingsByCampus(campusID)
	if err != nil {
		logger.Record("获取校区的全部建筑出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}

	result["buildings"] = buildings
	common.ResolveResult(c, true, e.SUCCESS, result)
}

// GetBuildingLayer 获取建筑楼高
func (bc *BuildingController) GetBuildingLayer(c *gin.Context) {
	result := map[string]interface{}{
		"building_layer": model.Building{},
	}

	// 1. 解析请求数据
	id, err := strconv.Atoi(c.Param("campus_id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "建筑ID必须为数字")
		return
	}
	// 2. 访问数据库
	layer, err := bc.BuildingService.GetBuildingLayer(id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "所选建筑不存在")
			return
		}
		logger.Record("根据ID获取建筑详细信息出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 3. 返回结果
	result["building_layer"] = layer
	common.ResolveResult(c, true, e.SUCCESS, result)
}
