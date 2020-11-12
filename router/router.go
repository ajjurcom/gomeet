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
	}

	apiUser := r.Group(relativePath)
	apiUser.Use(VerifyPower("user"))
	{
		/*
		 * 用户
		 */
		apiUser.PUT("/user", userController.Put)
		apiUser.PUT("/user_password", userController.PutPassword)
		apiUser.GET("/user/:id", userController.GetUserByID)
		apiUser.GET("/users/:onePageCount/:page", userController.GetUsersByPage)
		/*
		 * 校区所需API
		 */
		apiUser.GET("/campus", campusController.GetAllCampus)
		apiUser.GET("/campus/:onePageCount/:page", campusController.GetCampusByPage)
		/*
		 * 建筑所需API
		 */

		apiUser.GET("/buildings/:onePageCount/:page", buildingController.GetBuildingsByPage)
		apiUser.GET("/building/:id", buildingController.GetBuildingByID)
		/*
		 * 会议室所需API
		 */
		apiUser.GET("/campus_layer/:campus_id", buildingController.GetBuildingLayer)
		apiUser.GET("/campus_buildings/:campus_id", buildingController.GetAllBuildingsByCampus)
		apiUser.GET("/meeting/:id", meetingController.GetMeetingByID)
		apiUser.GET("/meetings/:onePageCount/:page", meetingController.GetMeetingsByPage)
		apiUser.GET("/meeting_options", meetingController.GetMeetingOptions)
	}

	apiAdmin := r.Group(relativePath)
	apiAdmin.Use(VerifyPower("admin"))
	{
		/*
		 * 管理用户API
		 */
		apiAdmin.DELETE("/user", userController.Delete)
		/*
		 * 校区所需API
		 */
		apiAdmin.POST("/campus", campusController.PostCampus)
		apiAdmin.DELETE("/campus/:id", campusController.DeleteCampus)
		apiAdmin.PUT("/campus", campusController.PutCampus)
		/*
		 * 建筑所需API
		 */
		apiAdmin.POST("/building", buildingController.PostBuilding)
		apiAdmin.DELETE("/building/:id", buildingController.DeleteBuilding)
		apiAdmin.PUT("/building", buildingController.PutBuilding)
		/*
		 * 会议室所需API
		 */
		apiAdmin.POST("/meeting", meetingController.PostMeeting)
		apiAdmin.DELETE("/meeting/:id", meetingController.DeleteMeeting)
		apiAdmin.PUT("/meeting", meetingController.PutMeeting)
	}

	apiRoot := r.Group(relativePath)
	apiRoot.Use(VerifyPower("root"))
	{
	}

	return r
}