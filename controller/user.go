package controller

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/config"
	"com/mittacy/gomeet/e"
	"com/mittacy/gomeet/logger"
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"com/mittacy/gomeet/service"
	"database/sql"
	"github.com/gin-gonic/gin"
)

type IUserController interface {
	Post(c *gin.Context)
	Login(c *gin.Context)
}

func NewUserController() IUserController {
	repo := repository.NewUserRepository("user")
	userService := service.NewUserService(repo)
	return &UserController{userService}
}

type UserController struct {
	UserService service.IUserService
}

// 注册用户
func (uc *UserController) Post(c *gin.Context) {
	// 1. 解析请求
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 操作数据库
	// 判断学号是否存在
	exists, err := uc.UserService.IsExistsByAttr("sno", user.Sno)
	if err != nil {
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	if exists {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "学号已存在")
		return
	}
	// 判断手机号是否存在
	exists, err = uc.UserService.IsExistsByAttr("phone", user.Phone)
	if err != nil {
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	if exists {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "手机号已存在")
		return
	}
	// 加密密码
	user.Password = common.Encryption(user.Password)
	if err := uc.UserService.CreateUser(&user); err != nil {
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// 用户、管理员登录
func (uc *UserController) Login(c *gin.Context) {
	// 1. 解析请求
	user := model.Session{}
	if err := c.ShouldBindJSON(&user); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}

	var (
		loginWay = ""	// 登录方式：手机/学号
		wayVal = ""
		state = ""	// 账号状态
	)

	// 2. 用户登录方式
	if len(user.Phone) > 0 {
		loginWay = "phone"
		wayVal = user.Phone
	} else if len(user.Sno) > 0 {
		loginWay = "sno"
		wayVal = user.Sno
	} else {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "请检查输入值")
		return
	}
	// 3. 管理员登录还是普通用户登录 -> 验证权限
	var err error
	state, err = uc.UserService.GetStateByAttr(loginWay, wayVal)
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "账号不存在")
			return
		}
		logger.Record("数据库查询错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	if user.IsAdmin == 0 {	// 普通用户
		switch state {
		case model.StateVerifyUser:
			common.ResolveResult(c, false, e.NOT_POWER, nil, "账号还未通过审核")
			return
		case model.StateRefuse:
			common.ResolveResult(c, false, e.NOT_POWER, nil, "账号审核不通过")
			return
		}
	} else {	// 管理员
		isAdmin, err := uc.UserService.CheckAdminByAttr(loginWay, wayVal)
		if err != nil || !isAdmin {
			common.ResolveResult(c, false, e.NOT_POWER, nil, "您还不是管理员，可通过申请升级为管理员")
			return
		}
		if state == model.StateVerifyAdmin {
			common.ResolveResult(c, false, e.NOT_POWER, nil, "账号升级管理员还未通过审核")
			return
		}
	}
	// 4. 检查密码
	password, err := uc.UserService.GetPasswordByAttr(loginWay, wayVal)
	if err != nil {
		logger.Record("数据库操作错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 5. 比对结果
	if common.Encryption(user.Password) != password {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "密码错误")
		return
	}
	// 6. 生成token
	if user.Sno == "root" {
		user.IsRoot = 1
	}
	token, err := common.GenerateToken(&user)
	if err != nil {
		common.ResolveResult(c, false, e.SUCCESS, password, "登录成功，但生成token失败")
		return
	}
	// 7. 返回数据
	name := config.Cfg.Section("jwt").Key("tokenName").String()
	result := map[string]string{name: token}
	common.ResolveResult(c, true, e.SUCCESS, result)
}
