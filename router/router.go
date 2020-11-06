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

	api := r.Group(relativePath)
	{
		api.POST("/user", userController.Post)
		api.POST("/session", userController.Login)
		api.GET("/campus", campusController.GetAllCampus)
		api.GET("/campus/:onePageCount/:page", campusController.GetCampusByPage)
		api.GET("/buildings/:onePageCount/:page", buildingController.GetBuildingByCampusAndPage)
		api.GET("/building/:id", buildingController.GetBuildingByID)

		// apiAdmin
		api.POST("/campus", campusController.PostCampus)
		api.DELETE("/campus/:id", campusController.DeleteCampus)
		api.PUT("/campus", campusController.PutCampus)

		api.POST("/building", buildingController.PostBuilding)
		api.DELETE("/building/:id", buildingController.DeleteBuilding)
		api.PUT("/building", buildingController.PutBuilding)
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