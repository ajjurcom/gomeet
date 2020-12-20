package controller

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/config"
	"com/mittacy/gomeet/e"
	"com/mittacy/gomeet/logger"
	"database/sql"
	"strconv"

	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"com/mittacy/gomeet/service"
	"github.com/gin-gonic/gin"
	"strings"
)

type IUserController interface {
	Post(c *gin.Context)
	Delete(c *gin.Context)
	Put(c *gin.Context)
	PutPassword(c *gin.Context)
	PutState(c *gin.Context)
	Login(c *gin.Context)
	GetUsersByPage(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetUserStateOptions(c *gin.Context)
	SearchUsers(c *gin.Context)
	GetAllUserByIDs(c *gin.Context)
	ApplyAdmin(c *gin.Context)
}

func NewUserController() IUserController {
	repo := repository.NewUserRepository("user")
	userService := service.NewUserService(repo)

	repoCampus := repository.NewCampusRepository("campus")
	campusService := service.NewCampusService(repoCampus)

	repoGroup := repository.NewGroupRepository("user_group", "user")
	groupService := service.NewGroupService(repoGroup)

	appointmentRepo := repository.NewAppointmentRepository("appointment", "user")
	appointmentService := service.NewAppointmentService(appointmentRepo)

	return &UserController{userService, campusService, groupService, appointmentService}
}

type UserController struct {
	UserService service.IUserService
	CampusService service.ICampusService
	GroupService service.IGroupService
	AppointmentService service.IAppointmentService
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
		if strings.Contains(err.Error(), "Duplicate") {
			msg := "学号已存在"
			if strings.Contains(err.Error(), "uidx_phone") {
				msg = "手机号已存在"
			}
			logger.Record("新建用户错误", err)
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
	result := map[string]interface{}{name: "", "id": 0, "username": ""}
	// 1. 解析请求
	session := model.Session{}
	if err := c.ShouldBindJSON(&session); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}

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
	// 7. 查询用户ID和名字
	id, username, err := uc.UserService.GetIDNameByAttr(loginWay, wayVal)
	if err != nil {
		common.ResolveResult(c, false, e.SUCCESS, password, "登录成功，但获取个人ID和名字失败")
		return
	}
	// 8. 返回数据
	result[name] = token
	result["expire"], _ = config.Cfg.Section("jwt").Key("expire").Int()
	result["id"] = id
	result["username"] = username
	result["isRoot"] = session.IsRoot
	result["state"] = state
	common.ResolveResult(c, true, e.SUCCESS, result)
}

// Delete 删除用户
func (uc *UserController) Delete(c *gin.Context) {
	// 1. 解析请求数据
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "用户ID必须为数字")
		return
	}
	// 2. 访问数据库
	if err := uc.UserService.DeleteUser(id); err != nil {
		if strings.Contains(err.Error(), "no exists") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "用户不存在")
		} else {
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// Put 更新用户信息
func (uc *UserController) Put(c *gin.Context) {
	// 1. 解析请求
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}

	// 2. 验证用户信息，避免利用自己的token修改他人的信息
	token := c.Request.Header.Get(config.Cfg.Section("jwt").Key("tokenName").String())
	if token == "" {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "你不可修改他人信息")
		return
	}
	if session, _ := common.ParseToken(token); session.Sno != user.Sno {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "你不可修改他人信息")
		return
	}

	// 3. 操作数据库
	if err := uc.UserService.PutUser(&user); err != nil {
		// 处理错误：用户是否存在、学号是否存在
		if strings.Contains(err.Error(), "Duplicate") {
			msg := "学号已存在"
			if strings.Contains(err.Error(), "uidx_phone") {
				msg = "手机号已存在"
			}
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, msg)
			return
		}
		logger.Record("更新用户错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 4. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// Put 更新用户密码
func (uc *UserController) PutPassword(c *gin.Context) {
	// 1. 解析请求
	user := model.Session{}
	if err := c.ShouldBindJSON(&user); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	/*
	1. 验证旧密码
	2. 修改新密码
	 */
	oldPwd, err := uc.UserService.GetPasswordByAttr("id", strconv.Itoa(user.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "用户不存在")
			return
		}
		logger.Record("通过ID查询用户密码", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	if common.Encryption(user.OldPassword) != oldPwd {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "旧密码不正确")
		return
	}

	if err := uc.UserService.PutPassword(user.ID, common.Encryption(user.Password)); err != nil {
		logger.Record("更新用户密码错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// PutState 修改状态
func (uc *UserController) PutState(c *gin.Context) {
	// 1. 解析请求
	newState := c.Query("state")
	if newState == "" {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "缺少新状态参数")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "用户ID必须为数字")
		return
	}
	// 2. 验证权限
	token := c.Request.Header.Get(config.Cfg.Section("jwt").Key("tokenName").String())
	session, _ := common.ParseToken(token)
	if !session.IsRoot && (newState == "verify_admin" || newState == "normal_admin" || newState == "root") {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "该操作只能由Root发起")
		return
	}
	// 3. 修改到数据库
	oldState, err := uc.UserService.GetStateByAttr("id", strconv.Itoa(id))
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "no exists")
		} else {
			logger.Record("更新用户状态出错", err)
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		}
		return
	}
	if err := uc.UserService.PutUserState(id, newState); err != nil {
		if strings.Contains(err.Error(), "Data truncated for column 'state'") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "状态参数错误")
			return
		}
		logger.Record("更新用户状态出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil, "更新用户状态失败")
		return
	}
	// 4. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
	// 5. 如果是审核通过用户，发送邮件通知
	go func() {
		if oldState == model.VerifyUser {
			// 查询user信息
			user, err := uc.UserService.GetUserByID(id)
			if err != nil {
				logger.Record("用户通过审核，获取用户信息失败", err)
				return
			}
			email := NewEmail("userVerify", user)
			if email == nil {
				logger.Record("获取邮件接口错误")
				return
			}

			isPass := true
			if newState == model.RefuseUser {
				isPass = false
				if err = uc.UserService.DeleteUser(id); err != nil {
					logger.Record("拒绝用户申请后，删除用户失败", err)
				}
			}
			if err = email.SendEmail(isPass); err != nil {
				logger.Record("发送用户审核邮件失败", err)
				return
			}
		} else if oldState == model.VerifyAdmin {
			// 查询user信息
			user, err := uc.UserService.GetUserByID(id)
			if err != nil {
				logger.Record("管理员审核中，获取用户信息失败", err)
				return
			}
			email := NewEmail("adminVerify", user)
			if email == nil {
				logger.Record("获取邮件接口错误")
				return
			}

			isPass := true
			if newState == model.NormalUser {
				isPass = false
			}
			if err = email.SendEmail(isPass); err != nil {
				logger.Record("发送管理员审核邮件失败", err)
				return
			}
		}
	}()
}

// ApplyAdmin 申请升级为管理员
func (uc *UserController) ApplyAdmin(c *gin.Context) {
	// 1. 解析请求
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "用户ID必须为数字")
		return
	}
	// 2. 修改到数据库
	if err = uc.UserService.PutUserState(id, model.VerifyAdmin); err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "no exists")
		} else {
			logger.Record("申请升级管理员出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 返回结果
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

// GetUsersByPage 获取用户
func (uc *UserController) GetUsersByPage(c *gin.Context) {
	result := map[string]interface{}{
		"userList": []model.User{},
		"count": 0,
	}
	// 1. 解析请求
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	onePageCount, err := strconv.Atoi(c.Param("onePageCount"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	state := c.Query("state")
	/* 数据操作
	 * 1. 获取state状态的用户总数量
	 * 2. 获取state用户列表
	 */
	count, err := uc.UserService.GetCountByState(state)
	if err != nil {
		logger.Record("获取用户数量出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}

	userList, err := uc.UserService.GetUserByPage(page, onePageCount, state)
	if err != nil && err != sql.ErrNoRows {
		logger.Record("获取用户错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	result["count"] = count
	result["userList"] = userList
	common.ResolveResult(c, true, e.SUCCESS, result)
}

// GetUserByID 获取用户详细信息
func (uc *UserController) GetUserByID(c *gin.Context) {
	result := map[string]interface{}{
		"user": model.User{},
	}
	// 1. 解析请求
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "用户id必须为数字")
		return
	}
	// 2. 数据库查询
	user, err := uc.UserService.GetUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, result, "用户不存在")
		} else {
			logger.Record("获取用户详细信息出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, result)
		}
		return
	}
	// 3. 返回结果
	result["user"] = user
	common.ResolveResult(c, true, e.SUCCESS, result)
}

/* GetUserStateOptionsByAdmin 用户管理初始化数据
 * 1. 参数 ?role=root  root
 * 2. 参数 ?role=admin 管理者
 */
func (uc *UserController) GetUserStateOptions(c *gin.Context) {
	role := c.Query("role")
	if role == "" || (role != "admin" && role != "root") {
		msg := "role值错误"
		if role == "" {
			msg = "缺少role参数"
		}
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, msg)
		return
	}
	result := map[string]interface{}{
		"stateList": model.StateOptions(role),
	}
	common.ResolveResult(c, true, e.SUCCESS, result)
}

/* 通过关键字搜索用户, 返回 [{"id": 1, "sno/phone", "username": "陈铭涛"}, ……]
 * 1. 学号
 * 2. 手机号
 * 3. 用户名
 */
func (uc *UserController) SearchUsers(c *gin.Context) {
	result := map[string]interface{}{
		"userList": []model.User{},
	}
	// 1. 解析参数
	searchWay := c.Query("searchWay")
	if searchWay != "phone" && searchWay != "sno" && searchWay != "username" {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "查询方法只能为学号、手机号、用户名")
		return
	}
	// 2. 请求数据库
	keyword := c.Query("keyword")
	userList, err := uc.UserService.SearchUsersByAttr(searchWay, keyword)
	if err != nil {
		if strings.Contains(err.Error(), "查询方法只能为") {
			common.ResolveResult(c, false, e.INVALID_PARAMS, result, err.Error())
		} else {
			logger.Record("根据关键字搜索用户错误", err)
			common.ResolveResult(c, false, e.BACK_ERROR, result)
		}
		return
	}

	result["userList"] = userList
	common.ResolveResult(c, true, e.SUCCESS, result)
}

/* 获取多个成员
 * way: 1. user_group
 * 		2. appointment
 */
func (uc *UserController) GetAllUserByIDs(c *gin.Context) {
	result := map[string]interface{}{
		"idList": []int{},
		"userList": []model.User{},
	}
	// 1. 解析请求
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "用户id必须为数字")
		return
	}
	way := c.Query("way")
	// 2. 数据库查询
	// 2.1 获取用户组id列表
	var group model.Group
	var appointment model.Appointment
	var ids string
	if way == "user_group" {
		group, err = uc.GroupService.GetGroupByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				common.ResolveResult(c, false, e.INVALID_PARAMS, result, "用户组不存在")
			} else {
				logger.Record("获取用户组成员信息出错", err)
				common.ResolveResult(c, false, e.BACK_ERROR, result)
			}
			return
		}
		ids = strings.Trim(group.Members, ",")
	} else if way == "appointment" {
		appointment, err = uc.AppointmentService.GetAppointmentById(id)
		if err != nil {
			if err == sql.ErrNoRows {
				common.ResolveResult(c, false, e.INVALID_PARAMS, result, "用户组不存在")
			} else {
				logger.Record("获取用户组成员信息出错", err)
				common.ResolveResult(c, false, e.BACK_ERROR, result)
			}
			return
		}
		ids = strings.Trim(appointment.Members, ",")
	} else {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "way参数错误")
		return
	}
	// 2.2 获取用户组用户信息
	userList, err := uc.UserService.GetAllUsersByIDs(ids)
	if err != nil {
		logger.Record("获取用户组成员信息出错", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}

	// 转化 ID字符串 -> 整数类型列表
	idStrList := common.MemberStrToList(ids)
	idList := make([]int, len(idStrList))
	for i := 0; i < len(idStrList); i++ {
		idList[i], _ = strconv.Atoi(idStrList[i])
	}

	// 3. 返回结果
	result["idList"] = idList
	result["userList"] = userList
	common.ResolveResult(c, true, e.SUCCESS, result)
}
