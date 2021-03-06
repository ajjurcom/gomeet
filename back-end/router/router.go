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
	appointmentController := controller.NewAppointmentController()
	recordController := controller.NewRecordController()

	api := r.Group(relativePath)
	{
		/*
		 * 统计数据
		 */
		api.GET("/statistics", recordController.StatisticsAppointment)
		api.GET("/statistics_options", recordController.StatisticsOptions)
		/*
		 * 会议
		 */
		api.GET("/reserve", appointmentController.GetAllReserve)
		api.GET("/my_reserve", appointmentController.GetMyAppointments)
		api.GET("/appointment/:id", appointmentController.GetAppointment)
		api.GET("/appointments/:onePageCount/:page", appointmentController.GetAppointmentByPage)
		api.GET("/appointment_states", appointmentController.GetAppointmentStates)
		/*
		 * 用户
		 */
		api.GET("/verify_code", userController.VerifyCode)
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
		api.GET("/user_group/:onePageCount/:page", groupController.GetGroupsByPage)
		api.GET("/user_groups/:creator", groupController.GetAllGroups)
		/*
		 * 校区所需API
		 */
		api.GET("/campus", campusController.GetAllCampus)
		api.GET("/campus/:onePageCount/:page", campusController.GetCampusByPage)
		/*
		 * 建筑所需API
		 */
		api.GET("/buildings/:onePageCount/:page", buildingController.GetBuildingsByPage)
		api.GET("/search_buildings/:onePageCount/:page", buildingController.SearchBuildings)
		api.GET("/building/:id", buildingController.GetBuildingByID)
		/*
		 * 会议室所需API
		 */
		api.GET("/campus_layer/:campus_id", buildingController.GetBuildingLayer)
		api.GET("/campus_buildings/:campus_id", buildingController.GetAllBuildingsByCampus)
		api.GET("/meeting/:id", meetingController.GetMeetingByID)
		api.GET("/meetings/:onePageCount/:page", meetingController.GetMeetingsByPage)
		api.GET("/search_meetings/:onePageCount/:page", meetingController.SearchMeetings)
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
		 * 会议
		 */
		apiUser.POST("/appointment", appointmentController.Post)	// 创建会议
		apiUser.POST("/appointment_fast", appointmentController.FastPost) 	// 快速创建会议，智能推荐
		apiUser.DELETE("/appointment/:id", appointmentController.Delete)	// 删除会议
		apiUser.PUT("/appointment", appointmentController.Put)	// 修改会议信息
		/*
		 * 用户
		 */
		apiUser.PUT("/user", userController.Put)	// 修改用户信息
		apiUser.PUT("/user_password", userController.PutPassword)	// 修改用户密码
		apiUser.PUT("/apply_admin/:id", userController.ApplyAdmin)	// 用户申请升级为管理员
		/*
		 * 用户组
		 */
		apiUser.POST("/user_group", groupController.Post)		// 创建用户组
		apiUser.DELETE("/user_group/:id", groupController.Delete)	// 删除用户组
		apiUser.PUT("/user_name", groupController.PutName)	// 修改用户组名字
		apiUser.PUT("/user_member", groupController.PutMember)	// 修改用户组成员
	}

	apiAdmin := r.Group(relativePath)
	apiAdmin.Use(VerifyPower("admin"))
	{
		/*
		 * 会议API
		 */
		apiAdmin.PUT("/appointment_state", appointmentController.PutState)	// 修改会议状态
		/*
		 * 管理用户API
		 */
		apiAdmin.DELETE("/user/:id", userController.Delete)	// 删除用户
		apiAdmin.PUT("/user_state/:id", userController.PutState)	// 修改用户状态
		/*
		 * 校区所需API
		 */
		apiAdmin.POST("/campus", campusController.PostCampus)	// 创建校区
		apiAdmin.DELETE("/campus/:id", campusController.DeleteCampus)	// 删除校区
		apiAdmin.PUT("/campus", campusController.PutCampus)	// 修改校区信息
		/*
		 * 建筑所需API
		 */
		apiAdmin.POST("/building", buildingController.PostBuilding)		// 增加建筑
		apiAdmin.DELETE("/building/:id", buildingController.DeleteBuilding)	// 删除建筑
		apiAdmin.PUT("/building", buildingController.PutBuilding)	// 修改建筑信息
		/*
		 * 会议室所需API
		 */
		apiAdmin.POST("/meeting", meetingController.PostMeeting)	// 创建会议室
		apiAdmin.DELETE("/meeting/:id", meetingController.DeleteMeeting)	// 删除会议室
		apiAdmin.PUT("/meeting", meetingController.PutMeeting)	// 修改会议室信息
	}

	apiRoot := r.Group(relativePath)
	apiRoot.Use(VerifyPower("root"))
	{
	}

	return r
}