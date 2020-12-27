package controller

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/e"
	"com/mittacy/gomeet/logger"
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"com/mittacy/gomeet/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ICampusController interface {
	PostCampus(c *gin.Context)
	DeleteCampus(c *gin.Context)
	PutCampus(c *gin.Context)
	GetAllCampus(c *gin.Context)
	GetCampusByPage(c *gin.Context)
}

func NewCampusController() ICampusController {
	repo := repository.NewCampusRepository("campus")
	campusService := service.NewCampusService(repo)
	return &CampusController{campusService}
}

type CampusController struct {
	CampusService service.ICampusService
}

// PostCampus 新增校区
func (cc *CampusController) PostCampus(c *gin.Context) {
	// 1. 解析请求数据
	campus := model.Campus{}
	if err := c.ShouldBindJSON(&campus); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 访问数据库
	if err := cc.CampusService.AddCampus(campus); err != nil {
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// DeleteCampus 删除某个校区
func (cc *CampusController) DeleteCampus(c *gin.Context) {
	// 1. 解析请求数据
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "the id must be a number")
		return
	}
	// 2. 访问数据库
	if err := cc.CampusService.DeleteCampus(id); err != nil {
		if err.Error() == "no exists" {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "the campus no exists")
		} else {
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// PutCampus 更新校区信息
func (cc *CampusController) PutCampus(c *gin.Context) {
	// 1. 解析请求数据
	campus := model.Campus{}
	if err := c.ShouldBindJSON(&campus); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 访问数据库
	if err := cc.CampusService.UpdateCampus(campus); err != nil {
		if err.Error() == "no exists" {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "the campus no exists")
		} else {
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// GetAllCampus 查询所有校区信息
func (cc *CampusController) GetAllCampus(c *gin.Context) {
	// 1. 操作数据库
	result := map[string]interface{} {
		"campusList": []model.Campus{},
	}
	campus, err := cc.CampusService.GetAllCampus()
	if err != nil {
		logger.Record("获取数据库出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	// 2. 返回结果
	result["campusList"] = campus
	common.ResolveResult(c, true, e.SUCCESS, result)
}

// GetCampusByPage 分页获取校区
func (cc *CampusController) GetCampusByPage(c *gin.Context) {
	result := map[string]interface{}{
		"campusList": []model.Campus{},
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
	// 2. 数据操作
	// 获取校区数量
	count, err := cc.CampusService.GetCampusCount()
	if err != nil {
		logger.Record("获取数据库出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	result["count"] = count
	// 获取校区
	campusList, err := cc.CampusService.GetCampusByPage(page, onePageCount)
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	result["campusList"] = campusList
	common.ResolveResult(c, true, e.SUCCESS, result)
}
