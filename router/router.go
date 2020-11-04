package router

import (
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由器，路由规则
func InitRouter() *gin.Engine {
	r := gin.New()
	
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	gin.SetMode("debug")	// debug / release / test
	

	return r
}