package common

import (
	"com/mittacy/gomeet/e"
	"github.com/gin-gonic/gin"
)

// ResolveResult 封装返回信息
func ResolveResult(c *gin.Context, success bool, code int, data interface{}, msg ...string) {
	returnMsg := ""
	if len(msg) == 0 {
		returnMsg = e.GetMsg(code)
	} else {
		returnMsg = msg[0]
	}
	c.JSON(code, gin.H{
		"success": success,
		"msg": returnMsg,
		"data": data,
	})
}
