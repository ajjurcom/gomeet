package controller

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/e"
	"com/mittacy/gomeet/logger"
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"com/mittacy/gomeet/service"
	"database/sql"
	"github.com/gin-gonic/gin"
	"strings"
)

type IAppointmentController interface {
	Post(c *gin.Context)
	Delete(c *gin.Context)
}

func NewAppointmentController() IAppointmentController {
	groupRepo := repository.NewGroupRepository("user_group", "user")
	groupService := service.NewGroupService(groupRepo)

	appointmentRepo := repository.NewAppointmentRepository("appointment", "user")
	appointmentService := service.NewAppointmentService(appointmentRepo)

	return &AppointmentController{appointmentService,groupService}
}

type AppointmentController struct {
	AppointmentService service.IAppointmentService
	GroupService service.IGroupService
}

func (ac *AppointmentController) Post(c *gin.Context) {
	/*
	 * 1. 解析请求
	 * 2. 查看会议时间是否冲突
	 * 3. 将groups中的组成员查询出来加入到members中
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
	// 2. 查看会议是否已占用, isConflict 会议是否冲突
	isConflict, err := ac.AppointmentService.IsAppointmentConflict(appointment)
	if err != nil {
		logger.Record("检查会议是否冲突错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	if isConflict {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议时间冲突")
		return
	}
	// 3. 将groups中的组成员查询出来加入到members中, 并去除相同用户ID
	members := appointment.Members
	if strings.Trim(appointment.Groups, ",") != "" {
		groups, err := ac.GroupService.GetMembersByGroups(appointment.Groups)
		if err != nil {
			logger.Record("查询组的所有用户ID出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
			return
		}
		for i := 0; i < len(groups); i++ {
			members += "," + groups[i].Members
		}
	}
	tempMembers := common.RemoveDuplicateEle(common.MemberStrToList(strings.Trim(members, ",")))
	members = common.MemberListToStr(tempMembers)
	// 4. 添加会议室
	if err := ac.AppointmentService.CreateAppointment(appointment, members); err != nil {
		logger.Record("添加会议错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	common.ResolveResult(c, true, e.SUCCESS, nil)
}

func (ac *AppointmentController) Delete(c *gin.Context) {
	/*
	 * 1. 解析请求, 获得删除会议的ID
	 * 2. 检查删除的会议是否为该用户创建
	 * 3. 将会议中的所有成员查询出来
	 * 4. 删除会议
	 */
	appointment := model.Appointment{}
	if err := c.ShouldBindJSON(&appointment); err != nil || appointment.ID == 0 || appointment.CreatorID == 0 {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil)
		return
	}

	// 2. 检查删除的会议是否为该用户创建
	// 3. 将会议中的所有成员查询出来
	members, creatorID, err := ac.AppointmentService.GetAllMembersAndCreatorIDByID(appointment.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "会议不存在")
		} else {
			logger.Record("获取全部成员出错", err)
			common.ResolveResult(c, false, e.BACK_ERROR, nil)
		}
		return
	}
	if appointment.CreatorID != creatorID {
		common.ResolveResult(c, false, e.INVALID_PARAMS, nil, "不能删除其他用户创建的会议")
		return
	}
	// 4. 删除会议
	if err := ac.AppointmentService.DeleteAppointment(appointment.ID, members); err != nil {
		logger.Record("删除会议错误", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}

	common.ResolveResult(c, true, e.SUCCESS, nil)
}
