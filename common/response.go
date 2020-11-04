package common

import (
	"com/mittacy/gomeet/e"
	"github.com/gin-gonic/gin"
)

// ResolveResult 封装返回信息
//
// 参数 (c *gin.Context, state bool, code int, data interface{}, msg ...string)
//
// 1. state 是否成功 true/false
//
// 2. code	返回状态码
//
// 3. data 返回数据
//
// 4. msg 提示信息，缺省该参数则使用默认提示信息
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
