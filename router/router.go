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
	scheduleController := controller.NewScheduleController()
	groupController := controller.NewGroupController()

	api := r.Group(relativePath)
	{
		/*
		 * 用户
		 */
		api.GET("/members/:id", userController.GetAllUserByIDs)
		api.GET("/users", userController.SearchUsers)
		api.POST("/user", userController.Post)
		api.POST("/session", userController.Login)
		api.GET("/user/:id", userController.GetUserByID)
		api.GET("/users/:onePageCount/:page", userController.GetUsersByPage)
		api.GET("/user_options", userController.GetUserStateOptions)
		/*
		 * 用户组
		 */
		api.GET("/user_group/:onePageCount/:page", groupController.GetMeetingsByPage)
		/*
		 * 校区所需API
		 */
		api.GET("/campus", campusController.GetAllCampus)
		api.GET("/campus/:onePageCount/:page", campusController.GetCampusByPage)
		/*
		 * 建筑所需API
		 */
		api.GET("/buildings/:onePageCount/:page", buildingController.GetBuildingsByPage)
		api.GET("/building/:id", buildingController.GetBuildingByID)
		/*
		 * 会议室所需API
		 */
		api.GET("/campus_layer/:campus_id", buildingController.GetBuildingLayer)
		api.GET("/campus_buildings/:campus_id", buildingController.GetAllBuildingsByCampus)
		api.GET("/meeting/:id", meetingController.GetMeetingByID)
		api.GET("/meetings/:onePageCount/:page", meetingController.GetMeetingsByPage)
		api.GET("/meeting_options", meetingController.GetMeetingOptions)
		/*
		 * 预定页面所需API
		 */
		api.GET("/schedule_options", scheduleController.GetOptions)
		api.GET("/options", scheduleController.UpdateOptions)
	}

	apiUser := r.Group(relativePath)
	apiUser.Use(VerifyPower("user"))
	{
		/*
		 * 用户
		 */
		apiUser.PUT("/user", userController.Put)
		apiUser.PUT("/user_password", userController.PutPassword)

		/*
		 * 用户组
		 */
		apiUser.POST("/user_group", groupController.Post)
		apiUser.DELETE("/user_group/:id", groupController.Delete)
		apiUser.PUT("/user_name", groupController.PutName)
		apiUser.PUT("/user_member", groupController.PutMember)
	}

	apiAdmin := r.Group(relativePath)
	apiAdmin.Use(VerifyPower("admin"))
	{
		/*
		 * 管理用户API
		 */
		apiAdmin.DELETE("/user/:id", userController.Delete)
		apiAdmin.PUT("/user_state/:id", userController.PutState)
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