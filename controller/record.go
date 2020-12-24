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
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	HotMeeting = "热门会议室Top10"
	ColdMeeting = "冷门会议室Top10"
	ZeroMeeting = "无预约会议室"
)

var options = []string{HotMeeting, ColdMeeting, ZeroMeeting}

type IRecordController interface {
	TransferExpireAppointment()
	StatisticsAppointment(c *gin.Context)
	StatisticsOptions(c *gin.Context)
}

func NewRecordController() IRecordController {
	appointmentRepo := repository.NewAppointmentRepository("appointment", "user")
	appointmentService := service.NewAppointmentService(appointmentRepo)

	meetingRepo := repository.NewMeetingRepository("meeting")
	meetingService := service.NewMeetingService(meetingRepo)

	recordRepo := repository.NewRecordRepository("appointment", "user", "record")
	recordService := service.NewRecordService(recordRepo)

	return &RecordController{
		AppointmentService: appointmentService,
		MeetingService: meetingService,
		RecordService:      recordService,
	}
}

type RecordController struct {
	AppointmentService service.IAppointmentService
	MeetingService service.IMeetingService
	RecordService      service.IRecordService
}

// 清除过期会议，转移到 日志表
func (rc *RecordController) TransferExpireAppointment() {
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
	appointments, err := rc.AppointmentService.GetExpireAppointment(day, endTime)
	if err != nil {
		logger.Record("!!! 获取过期会议错误, 请尽快处理 !!!")
		return
	}
	// 3. 遍历每个会议，将会议逐个转移到日志表
	for _, v := range appointments {
		members, _, err := rc.AppointmentService.GetAllMembersAndCreatorIDByID(v.ID)
		if err != nil {
			logger.Record("!!! 获取会议所有成员错误, 请尽快处理 !!!")
			return
		}
		if err = rc.RecordService.TransferAppointment(v, members); err != nil {
			logger.Record("!!! 转移会议错误, 请尽快处理 !!!")
			return
		}
	}
	fmt.Println(time.Now(), "会议转移成功!")
}

func (rc *RecordController) StatisticsOptions(c *gin.Context) {
	result := map[string][]string{
		"options": options,
	}
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (rc *RecordController) StatisticsAppointment(c *gin.Context) {
	result := map[string]interface{}{}
	isUpdate := c.DefaultQuery("updateDate", "false")

	now := time.Now()
	endDay := now.Format("20060102")
	startDay := now.AddDate(0, -1, 0).Format("20060102")
	startDay = c.DefaultQuery("startDay", startDay)
	endDay = c.DefaultQuery("endDay", endDay)

	// 1. 根据日期获取会议信息
	appointments, err := rc.RecordService.GetAppointmentByDay(startDay, endDay)
	if err != nil && err != sql.ErrNoRows {
		logger.Record("根据日期获取会议失败", err)
		common.ResolveResult(c, false, e.BACK_ERROR, nil)
		return
	}
	statisticsAppointment := map[string]int{
		model.AppointmentVerify: 0,
		model.AppointmentRefuse: 0,
		"pass":                  0,
	}
	for _, v := range appointments {
		statisticsAppointment[v.State] += 1
	}

	var count int
	var items []model.Item
	if isUpdate == "true" {
		count = len(appointments)
		ratio1 := float32(statisticsAppointment[model.AppointmentVerify]) / float32(count) * 100
		ratio2 := float32(statisticsAppointment[model.AppointmentRefuse]) / float32(count) * 100
		ratio3 := 100 - ratio1 - ratio2
		if count == 0 {
			ratio1 = 0
			ratio2 = 0
			ratio3 = 0
		}
		item1 := model.Item{
			Title: "未审核数",
			Num: statisticsAppointment[model.AppointmentVerify],
			Ratio: ratio1,
		}
		item2 := model.Item{
			Title: "审核被拒数",
			Num: statisticsAppointment[model.AppointmentRefuse],
			Ratio: ratio2,
		}
		tmp := count - statisticsAppointment[model.AppointmentVerify] - statisticsAppointment[model.AppointmentRefuse]
		item3 := model.Item{
			Title: "审核通过率",
			Num: tmp,
			Ratio: ratio3,
		}
		items = []model.Item{item1, item2, item3}
		result["count"] = count
		result["items"] = items
	}

	if isUpdate == "true" && count == 0 {		// 没有会议，统计其他数据没意义
		result["statisticsList"] = []string{}
		common.ResolveResult(c, true, e.SUCCESS, result)
		return
	}

	statisticsType := c.DefaultQuery("statisticsType", options[0])
	if statisticsType == HotMeeting || statisticsType == ColdMeeting {
		isHot := false
		if statisticsType == HotMeeting {
			isHot = true
		}
		meetings, err := rc.getStatisticsMeeting(appointments, isHot)
		if err != nil {
			common.ResolveResult(c, false, e.INVALID_PARAMS, result)
			return
		}
		result["statisticsList"] = meetings
	} else if statisticsType == ZeroMeeting {
		meetings, err := rc.getZeroMeetings(appointments)
		if err != nil {
			common.ResolveResult(c, false, e.INVALID_PARAMS, result)
			return
		}
		result["statisticsList"] = meetings
	} else {
		common.ResolveResult(c, false, e.INVALID_PARAMS, result, "统计参数错误")
		return
	}
	common.ResolveResult(c, true, e.SUCCESS, result)
}

func (rc *RecordController) getStatisticsMeeting(appointments []model.Appointment,
	isHot bool) (meetings []model.Meeting, err error) {
	reserveCountMap := map[int]int{}
	for _, v := range appointments {
		reserveCountMap[v.MeetingID] += 1
	}

	// 对map进行排序，按从大到小输出到orderAppointment
	type kv struct {
		Key   int
		Value int
	}
	var orderAppointment []kv
	for k, v := range reserveCountMap {
		orderAppointment = append(orderAppointment, kv{k, v})
	}
	if isHot {
		sort.Slice(orderAppointment, func(i, j int) bool {
			return orderAppointment[i].Value > orderAppointment[j].Value
		})
	} else {
		sort.Slice(orderAppointment, func(i, j int) bool {
			return orderAppointment[i].Value < orderAppointment[j].Value
		})
	}

	var ids string
	for i, v := range orderAppointment {
		if i == 10 {
			break
		}
		ids += strconv.Itoa(v.Key) + ","
	}
	ids = strings.Trim(ids, ",")
	meeting, err := rc.MeetingService.GetMeetingsByID(ids)
	if err != nil {
		return meeting, err
	}

	for i, v := range orderAppointment {
		if i == 10 {
			break
		}
		meeting[i].ReverseCount = v.Value
	}
	return meeting, nil
}

func (rc *RecordController) getZeroMeetings(appointments []model.Appointment) (meetings []model.Meeting, err error) {
	existsMeetings := map[int]struct{}{}
	for _, v := range appointments {
		existsMeetings[v.MeetingID] = struct{}{}
	}

	ids, err := rc.MeetingService.GetAllMeetingsID()
	resultID := make([]int, 0)
	for _, v := range ids {
		if _, ok := existsMeetings[v]; !ok {
			resultID = append(resultID, v)
		}
	}

	str := ""
	for _, v := range resultID {
		str += strconv.Itoa(v) + ","
	}
	str = strings.Trim(str, ",")
	return rc.MeetingService.GetMeetingsByID(str)
}
