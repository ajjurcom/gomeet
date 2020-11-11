package controller

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/config"
	"com/mittacy/gomeet/e"
	"com/mittacy/gomeet/logger"
	"database/sql"

	//"com/mittacy/gomeet/logger"
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"com/mittacy/gomeet/service"
	//"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type IUserController interface {
	Post(c *gin.Context)
	Login(c *gin.Context)
}

func NewUserController() IUserController {
	repo := repository.NewUserRepository("user")
	userService := service.NewUserService(repo)
	repoCampus := repository.NewCampusRepository("campus")
	campusService := service.NewCampusService(repoCampus)
	return &UserController{userService, campusService}
}

type UserController struct {
	UserService service.IUserService
	CampusService service.ICampusService
}

// 注册用户
func (uc *UserController) Post(c *gin.Context) {
	// 1. 解析请求
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 加密密码, 插入数据库
	user.Password = common.Encryption(user.Password)
	if err := uc.UserService.CreateUser(&user); err != nil {
		// 处理错误：用户是否存在、学号是否存在
		fmt.Println("创建用户错误: ", err)
		if strings.Contains(err.Error(), "Duplicate") {
			msg := "学号已存在"
			if strings.Contains(err.Error(), "uidx_phone") {
				msg = "手机号已存在"
			}
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, msg)
			return
		}
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// 用户、管理员登录
func (uc *UserController) Login(c *gin.Context) {
	name := config.Cfg.Section("jwt").Key("tokenName").String()
	result := map[string]string{name: ""}
	// 1. 解析请求
	session := model.Session{}
	if err := c.ShouldBindJSON(&session); err != nil {
		fmt.Println("登录用户错误: ", err)
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	fmt.Println(session)

	var (
		loginWay string	// 登录方式：手机/学号
		wayVal string
		state string	// 账号状态
	)
	// 2. 登录方式
	if len(session.Sno) > 0 {
		loginWay = "sno"
		wayVal = session.Sno
	} else if len(session.Phone) > 0 {
		loginWay = "phone"
		wayVal = session.Phone
	} else {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "请检查输入值")
		return
	}
	// 3. 管理员登录还是普通用户登录 -> 验证权限
	var err error
	state, err = uc.UserService.GetStateByAttr(loginWay, wayVal)
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, result, "账号不存在")
			return
		}
		logger.Record("查询用户状态错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}

	// 4. 根据不同身份有不同的验证
	msg := ""
	isOk := true
	if !session.IsAdmin {	// 普通用户登录
		fmt.Println("test...")
		switch state {
		case model.VerifyUser:
			msg = "账号还未通过审核"
			isOk = false
		case model.RefuseUser:
			msg = "账号审核不通过"
			isOk = false
		}
	} else {	// 管理员
		if state != model.NormalAdmin && state != model.Root {
			if state == model.VerifyAdmin {
				msg = "账号升级管理员还未通过审核"
			} else {
				msg = "您还不是管理员，可通过申请升级为管理员"
			}
			isOk = false
		}
	}
	if !isOk {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, msg)
		return
	}
	// 5. 检查密码
	password, err := uc.UserService.GetPasswordByAttr(loginWay, wayVal)
	if err != nil {
		logger.Record("获取用户密码错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	// 5. 比对结果
	if common.Encryption(session.Password) != password {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "密码错误")
		return
	}
	// 6. 生成token
	if state == model.Root {
		session.IsRoot = true
	}
	token, err := common.GenerateToken(&session)
	if err != nil {
		common.ResolveResult(c, false, e.SUCCESS, password, "登录成功，但生成token失败")
		return
	}
	//7. 返回数据
	result[name] = token
	common.ResolveResult(c, true, e.SUCCESS, result)
}
