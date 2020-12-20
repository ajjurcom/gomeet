package router

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/config"
	"com/mittacy/gomeet/e"
	"com/mittacy/gomeet/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CorsMiddleware 跨域控制
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		tokenName := config.Cfg.Section("jwt").Key("tokenName").String()
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "POST, DELETE, PUT, GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, " + tokenName)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

// VerifyPower 验证权限
func VerifyPower(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取token
		var (
			code int = e.SUCCESS
			//token string = c.Request.Header.Get(config.Cfg.Section("jwt").Key("tokenName").String())
			session *model.Session
		)
		token, err := c.Cookie(config.Cfg.Section("jwt").Key("tokenName").String())
		// 2. 检验有效性
		if err != nil {
			code = e.NOT_POWER
		} else {
			var err error
			session, err = common.ParseToken(token)
			if err != nil || session == nil {
				code = e.NOT_POWER
			} else {
				switch role {
				case "admin":
					if !session.IsAdmin {
						code = e.NOT_POWER
					}
				case "root":
					if !session.IsRoot {
						code = e.NOT_POWER
					}
				}
			}
		}

		// 3. 回复结果
		if code != e.SUCCESS {
			common.ResolveResult(c, false, code, nil, "跳转登录页面")
			c.Abort()
			return
		}
		c.Next()
	}
}
