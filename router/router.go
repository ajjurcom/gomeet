package router

import (
	"com/mittacy/gomeet/controller"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由器，路由规则
func InitRouter() *gin.Engine {
	r := gin.New()
	
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CorsMiddleware())
	
	gin.SetMode("debug")	// debug / release / test

	const relativePath = "/api/v1"

	// 初始化控制器
	userController := controller.NewUserController()
	campusController := controller.NewCampusController()
	buildingController := controller.NewBuildingController()
	meetingController := controller.NewMeetingController()

	api := r.Group(relativePath)
	{
		api.POST("/user", userController.Post)
		api.POST("/session", userController.Login)
		/*
		 * 校区所需API
		 */
		api.POST("/campus", campusController.PostCampus)
		api.DELETE("/campus/:id", campusController.DeleteCampus)
		api.PUT("/campus", campusController.PutCampus)
		api.GET("/campus/:onePageCount/:page", campusController.GetCampusByPage)
		/*
		 * 建筑所需API
		 */
		api.POST("/building", buildingController.PostBuilding)
		api.DELETE("/building/:id", buildingController.DeleteBuilding)
		api.PUT("/building", buildingController.PutBuilding)
		api.GET("/campus", campusController.GetAllCampus)
		api.GET("/buildings/:onePageCount/:page", buildingController.GetBuildingsByPage)
		api.GET("/building/:id", buildingController.GetBuildingByID)
		/*
		 * 会议室所需API
		 */
		api.GET("/campus_layer/:campus_id", buildingController.GetBuildingLayer)
		api.GET("/campus_buildings/:campus_id", buildingController.GetAllBuildingsByCampus)
		api.POST("/meeting", meetingController.PostMeeting)
		api.DELETE("/meeting/:id", meetingController.DeleteMeeting)
		api.PUT("/meeting", meetingController.PutMeeting)
		api.GET("/meeting/:id", meetingController.GetMeetingByID)
		api.GET("/meetings/:onePageCount/:page", meetingController.GetMeetingsByPage)
		api.GET("/meeting_options", meetingController.GetMeetingOptions)
	}

	apiUser := r.Group(relativePath)
	apiUser.Use(VerifyPower("user"))
	{

	}

	apiAdmin := r.Group(relativePath)
	apiAdmin.Use(VerifyPower("admin"))
	{

	}

	apiRoot := r.Group(relativePath)
	apiRoot.Use(VerifyPower("root"))
	{
	}

	return r
}