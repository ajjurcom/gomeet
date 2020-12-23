package controller

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/e"
	"com/mittacy/gomeet/logger"
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"com/mittacy/gomeet/service"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type IAppointmentController interface {
	Post(c *gin.Context)
	FastPost(c *gin.Context)
	Delete(c *gin.Context)
	Put(c *gin.Context)
	PutState(c *gin.Context)
	GetAllReserve(c *gin.Context)
	GetMyAppointments(c *gin.Context)
	GetAppointment(c *gin.Context)
	GetAppointmentByPage(c *gin.Context)
	GetAppointmentStates(c *gin.Context)
	TransferExpireAppointment()
}

func NewAppointmentController() IAppointmentController {
	groupRepo := repository.NewGroupRepository("user_group", "user")
	groupService := service.NewGroupService(groupRepo)

	appointmentRepo := repository.NewAppointmentRepository("appointment", "user")
	appointmentService := service.NewAppointmentService(appointmentRepo)

	userRepo := repository.NewUserRepository("user")
	userService := service.NewUserService(userRepo)

	meetingRepo := repository.NewMeetingRepository("meeting")
	meetingService := service.NewMeetingService(meetingRepo)

	buildingRepo := repository.NewBuildingRepository("building")
	buildingService := service.NewBuildingService(buildingRepo)

	campusRepo := repository.NewCampusRepository("campus")
	campusService := service.NewCampusService(campusRepo)

	return &AppointmentController{
		AppointmentService: appointmentService,
		GroupService: groupService,
		UserService: userService,
		MeetingService: meetingService,
		BuildingService: buildingService,
		CampusService: campusService,
	}
}

type AppointmentController struct {
	AppointmentService service.IAppointmentService
	GroupService service.IGroupService
	UserService service.IUserService
	MeetingService service.IMeetingService
	BuildingService service.IBuildingService
	CampusService service.ICampusService
}

func (ac *AppointmentController) Post(c *gin.Context) {
	/*
	 * 1. 解析请求
	 * 2. 查看会议时间是否冲突
	 * 3. 对members中的成员去重
	 * 4. 添加会议
	 */
	// 1. 解析请求
	appointment := model.Appointment{}
	if err := c.ShouldBindJSON(&appointment); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	if appointment.StartTime >= appointment.EndTime {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "开始时间不能等于或晚于结束时间")
		return
	}
	// 2. 查看会议是否冲突, isConflict 会议是否冲突
	isConflict, err := ac.AppointmentService.IsAppointmentConflict(appointment, "post")
	if err != nil {
		logger.Record("检查会议是否冲突错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	if isConflict {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议时间冲突")
		return
	}
	// 3. 对members中的成员去重
	appointment.Members = common.MemberListToStr(common.RemoveDuplicateEle(common.MemberStrToList(appointment.Members)))
	// 4. 添加会议室
	if err := ac.AppointmentService.CreateAppointment(appointment); err != nil {
		logger.Record("添加会议错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

func (ac *AppointmentController) FastPost(c *gin.Context) {
	/*
	 * 1. 解析请求
	 * 2. 查找一个随机的空余会议室ID
	 * 3. 添加会议
	 * 4. 返回会议室详细信息
	 */
	// 1. 解析请求
	campusID := c.Query("campus_id")
	meetingScale := c.Query("meeting_scale")
	meetingType := c.Query("meeting_type")
	appointment := model.Appointment{}
	if err := c.ShouldBindJSON(&appointment); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 查找一个空闲的会议室
	appointments, err := ac.AppointmentService.GetAppointmentsIDByTime(appointment)
	if err != nil {
		fmt.Println(err)
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	meetingID := make([]string, len(appointments))
	for i, v := range appointments {
		meetingID[i] = strconv.Itoa(v.MeetingID)
	}
	meetingIDList := common.MemberListToStr(meetingID)
	meeting, err := ac.MeetingService.GetMeetingByInfo(meetingIDList, campusID, meetingType, meetingScale)
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "该时间段会议室已满, 您可以更换时间或会议室类型")
		} else {
			logger.Record("获取空闲会议室失败", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	// 3. 会议室详细信息
	result := map[string]interface{}{
		"appointment": model.Appointment{},
	}
	appointment.MeetingID = meeting.ID
	if appointment.Locate, err = ac.getLocate(appointment.MeetingID); err != nil {
		logger.Record("获取会议室的位置错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	// 4. 添加会议
	if err := ac.AppointmentService.CreateAppointment(appointment); err != nil {
		logger.Record("添加会议错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	result["appointment"] = appointment
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (ac *AppointmentController) Delete(c *gin.Context) {
	/*
	 * 1. 解析请求, 获得删除会议的ID
	 * 2. 检查删除的会议是否为该用户创建
	 * 3. 将会议中的所有成员查询出来
	 * 4. 将该会议标记为退订: cancel / adopt-cancel
	 */
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "id无效")
		return
	}
	queryCreatorID, err := strconv.Atoi(c.Query("creator_id"))
	if err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "creator_id无效")
		return
	}

	// 2. 检查删除的会议是否为该用户创建
	// 3. 将会议中的所有成员查询出来
	members, creatorID, err := ac.AppointmentService.GetAllMembersAndCreatorIDByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议不存在")
		} else {
			logger.Record("获取全部成员出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	if queryCreatorID != creatorID {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "不能删除其他用户创建的会议")
		return
	}
	// 4. 删除会议前先保存会议
	appointment, err := ac.AppointmentService.GetAppointmentById(id)
	if err != nil {
		logger.Record("会议通过审核，获取会议详情失败", err)
		return
	}
	state := model.AppointmentCancel
	if appointment.State == model.AppointmentAdopt {
		state = model.AppointmentAdoptCancel
	}
	if err := ac.AppointmentService.PutState(id, state); err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议不存在")
		} else {
			logger.Record("删除会议，修改会议状态错误", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	common.ResolveResult(c, true, e.SUCCESS, nil)
	// 只有会议通过审核才发送退订通知
	if appointment.State == model.AppointmentAdopt {
		go func() {
			if appointment.Locate, err = ac.getLocate(appointment.MeetingID); err != nil {
				logger.Record("获取会议室的位置错误", err)
				return
			}
			users, err := ac.UserService.GetUsersEmailByID(members)
			if err != nil {
				logger.Record("通知会议参会成员时查询参会人员出错")
				return
			}
			email := NewEmail("notifyMembers", users, appointment)
			if email == nil {
				logger.Record("获取邮件接口错误")
				return
			}
			if err = email.SendEmail(false); err != nil {
				logger.Record("通知会议参会成员失败", err)
				return
			}
		}()
	}
}

func (ac *AppointmentController) Put(c *gin.Context) {
	/*
	 * 1. 解析请求
	 * 2. 查看新会议时间是否冲突
	 * 3. 查询旧成员
	 * 4. 比较得出计算新增加成员和删除成员
	 * 5. 操作数据库添加会议室
	 */
	// 1. 解析请求
	appointment := model.Appointment{}
	if err := c.ShouldBindJSON(&appointment); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}
	//if appointment.StartTime >= appointment.EndTime {
	//	common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "开始时间不能等于或晚于结束时间")
	//	return
	//}
	// 2. 查看会议是否冲突, isConflict 会议是否冲突
	//isConflict, err := ac.AppointmentService.IsAppointmentConflict(appointment, "put")
	//if err != nil {
	//	logger.Record("检查会议是否冲突错误", err)
	//	common.ResolveResult(c, false, e.BACK_ERROR, nil)
	//	return
	//}
	//if isConflict {
	//	common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议时间冲突")
	//	return
	//}
	// 3. 查询旧成员
	members, _, err := ac.AppointmentService.GetAllMembersAndCreatorIDByID(appointment.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议不存在")
		} else {
			logger.Record("获取全部成员出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	oldMembers := common.MemberStrToList(members)
	newMembers := common.RemoveDuplicateEle(common.MemberStrToList(appointment.Members))
	appointment.Members = common.MemberListToStr(newMembers)

	// 5. 计算新成员和删除成员, 所有成员
	deleteMembers, addMembers := common.DiffMember(oldMembers, newMembers)

	// 6. 修改会议信息
	if err := ac.AppointmentService.PutAppointment(appointment,
		common.MemberListToStr(addMembers),
		common.MemberListToStr(deleteMembers)); err != nil {
		logger.Record("更新会议错误", err)
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "请检查输入值")
		return
	}
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

func (ac *AppointmentController) PutState(c *gin.Context) {
	// 1. 解析请求
	appointment := model.Appointment{}
	if err := c.ShouldBindJSON(&appointment); err != nil {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}

	if !model.IsAppointmentState(appointment.State) {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "状态参数错误")
		return
	}

	// 2. 操作数据库
	err := ac.AppointmentService.PutState(appointment.ID, appointment.State)
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议不存在")
		} else {
			logger.Record("修改会议状态错误", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}

	common.ResolveResult(c, true, e.SUCCESS, nil)
	// 发送邮件通知创建者会议审核结果
	go func() {
		// 查询必要信息
		if appointment, err = ac.AppointmentService.GetAppointmentById(appointment.ID); err != nil {
			logger.Record("会议通过审核，获取会议详情失败", err)
			return
		}
		var user model.User
		if user, err = ac.UserService.GetUserByID(appointment.CreatorID); err != nil {
			logger.Record("获取会议创建者失败", err)
			return
		}
		isPass := true
		if appointment.State == model.AppointmentAdopt {
			if appointment.Locate, err = ac.getLocate(appointment.MeetingID); err != nil {
				logger.Record("获取会议室的位置错误", err)
				return
			}
		} else if appointment.State == model.AppointmentRefuse {
			isPass = false
		}
		email := NewEmail("appointmentVerify", user, appointment)
		if email == nil {
			logger.Record("获取邮件接口错误")
			return
		}
		if err = email.SendEmail(isPass); err != nil {
			logger.Record("发送会议审核邮件失败", err)
			return
		}
		// 如果是通过会议，需要通知所有参会成员
		if appointment.State == model.AppointmentAdopt {
			users, err := ac.UserService.GetUsersEmailByID(appointment.Members)
			if err != nil {
				logger.Record("通知会议参会成员时查询参会人员出错")
				return
			}
			if email = NewEmail("notifyMembers", users, appointment); email == nil {
				logger.Record("获取邮件接口错误")
				return
			}
			if err = email.SendEmail(true); err != nil {
				logger.Record("通知会议参会成员失败", err)
				return
			}
		}
	}()
}

// api/v1/all_creator?day=11/30/2020&meeting_id[]=...
func (ac *AppointmentController) GetAllReserve(c *gin.Context) {
	result := map[string]interface{}{
		"appointments": []model.Appointment{},
	}

	day := c.Query("day")
	startTime := c.Query("start_time")
	meetingID := c.QueryArray("meeting_id[]")
	if day == "" || len(meetingID) == 0 {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}
	appointments, err := ac.AppointmentService.GetAllReserve(day, startTime, common.MemberListToStr(meetingID))
	if err != nil {
		logger.Record("查询预定情况错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}

	result["appointments"] = appointments
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (ac *AppointmentController) GetMyAppointments(c *gin.Context) {
	result := map[string]interface{}{
		"myReserve": []model.Appointment{},
		"otherReserve": []model.Appointment{},
	}

	id, err := strconv.Atoi(c.Query("creator_id"))
	if err != nil || id == 0 {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}

	// 获取我预定的会议
	myReserve, err := ac.AppointmentService.GetMyAllReserve(id)
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		} else {
			logger.Record("获取我的所有会议错误", err)
			common.ResolveResult(c, false, e.BACK_ERROR, result)
		}
		return
	}

	// 获取受邀的会议
	s, err := ac.UserService.GetMyAppointmentsID(id)
	if err != nil {
		logger.Record("获取user的appointments错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}
	otherAppointments := common.RemoveDuplicateEle(common.MemberStrToList(s))
	var myAppointments []string
	for _, item := range myReserve {
		myAppointments = append(myAppointments, strconv.Itoa(item.ID))
	}
	otherAppointments = common.RemoveSameEle(otherAppointments, myAppointments)
	otherReserve, err := ac.AppointmentService.GetAppointmentsByID(common.MemberListToStr(otherAppointments))
	if err != nil {
		logger.Record("根据多个id字符串获取appointments错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}


	result["myReserve"] = myReserve
	result["otherReserve"] = otherReserve
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (ac *AppointmentController) GetAppointment(c *gin.Context) {
	result := map[string]interface{}{
		"appointment": model.Appointment{},
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}

	appointment, err := ac.AppointmentService.GetAppointmentById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, result, "该会议不存在")
		} else {
			logger.Record("获取会议详细信息错误", err)
			common.ResolveResult(c, false, e.BACK_ERROR, result)
		}
		return
	}

	if appointment.Locate, err = ac.getLocate(appointment.MeetingID); err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, result, "该会议不存在")
		} else {
			logger.Record("获取会议室的位置错误", err)
			common.ResolveResult(c, false, e.BACK_ERROR, result)
		}
		return
	}

	result["appointment"] = appointment
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (ac *AppointmentController) GetAppointmentByPage(c *gin.Context) {
	result := map[string]interface{}{
		"appointments": []model.Appointment{},
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
	if !model.IsAppointmentState(state) {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result)
		return
	}

	/*
	 * 1. 获取count
	 * 2. 获取分页数据
	 */
	count, err := ac.AppointmentService.GetCountByState(state)
	if err != nil {
		logger.Record("获取会议室数量错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}

	appointments, err := ac.AppointmentService.GetAppointmentsByPage(page, onePageCount, state)
	if err != nil && err != sql.ErrNoRows {
		logger.Record("分页获取会议错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, result)
		return
	}

	result["count"] = count
	result["appointments"] = appointments
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (ac *AppointmentController) GetAppointmentStates(c *gin.Context) {
	result := map[string][]string {
		"states": model.AppointmentStates(),
	}
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (ac *AppointmentController) getLocate(meetingID int) (string, error) {
	/*
	 * 1. 根据 会议室id 获取 会议室楼层、名字、建筑id、name
	 * 2. 根据 建筑id 获取 校区name
	 * 3. 组成 locate
	 */
	locate := ""
	meeting, err := ac.MeetingService.GetMeetingByID(meetingID)
	if err != nil {
		return locate, err
	}

	building, err := ac.BuildingService.GetBuildingByID(meeting.BuildingID)
	if err != nil {
		return locate, err
	}

	campus, err := ac.CampusService.GetCampusByID(building.CampusID)
	if err != nil {
		return locate, err
	}

	locate = campus.CampusName + " - " + building.BuildingName + " - F" + strconv.Itoa(meeting.Layer) + "-" + meeting.RoomNumber + "（" + meeting.MeetingName + "）"
	return locate, nil
}

// 清除过期会议，转移到 日志表
func (ac *AppointmentController) TransferExpireAppointment() {
	/*
	 * 1. 获取当前时间
	 * 2. 会议室查询过期会议
	 * 3. 遍历每个会议，将会议逐个转移到日志表
	 */
	fmt.Println(time.Now(), "开始转移过期会议")
	// 1. 获取当前时间
	now := time.Now()
	day := now.Format("20060102")
	endTime := now.Format("15") + ":00"
	// 2. 会议室查询过期会议
	appointments, err := ac.AppointmentService.GetExpireAppointment(day, endTime)
	if err != nil {
		logger.Record("!!! 获取过期会议错误, 请尽快处理 !!!")
		return
	}
	// 3. 遍历每个会议，将会议逐个转移到日志表
	for _, v := range appointments {
		fmt.Println("转移会议: ", v)
		members, _, err := ac.AppointmentService.GetAllMembersAndCreatorIDByID(v.ID)
		if err != nil {
			fmt.Println("会议转移失败!")
			logger.Record("!!! 获取会议所有成员错误, 请尽快处理 !!!")
			return
		}
		if err = ac.AppointmentService.TransferAppointment(v, members); err != nil {
			fmt.Println("会议转移失败!")
			logger.Record("!!! 转移会议错误, 请尽快处理 !!!")
			return
		}
	}
	fmt.Println(time.Now(), "会议转移成功!")
}
