package controller

import (
	"com/mittacy/gomeet/common"
	"com/mittacy/gomeet/config"
	"com/mittacy/gomeet/logger"
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
	"com/mittacy/gomeet/service"
	"strings"
)

type IEmail interface {
	SendEmail(bool) error
}

func NewEmail(emailName string, obj ...interface{}) IEmail {
	switch emailName {
	case "userVerify":
		return &UserVerifyEmail{obj[0].(model.User)}
	case "adminVerify":
		return &AdminVerifyEmail{obj[0].(model.User)}
	case "appointmentVerify":
		return &AppointmentVerifyEmail{obj[0].(model.User), obj[1].(model.Appointment)}
	case "notifyMembers":
		return &NotifyMembersEmail{obj[0].([]model.User), obj[1].(model.Appointment)}
	case "emailVerifyCode":
		return &EmailVerifyCodeEmail{obj[0].(model.User)}
	default:
		return nil
	}
}

type UserVerifyEmail struct {
	user model.User
}

type AdminVerifyEmail struct {
	user model.User
}

type AppointmentVerifyEmail struct {
	user model.User
	appointment model.Appointment
}

type NotifyMembersEmail struct {
	user []model.User
	appointment model.Appointment
}

type EmailVerifyCodeEmail struct {
	user model.User
}

func (pve *UserVerifyEmail) SendEmail(success bool) (err error) {
	/*
	 * 1. 获取邮件模板
	 * 2. 替换邮件模板变量
	 * 3. 发送邮件
	 */
	repo := repository.NewEmailRepository("email")
	emailService := service.NewEmailService(repo)
	var content, subject string
	if success {
		subject = common.PassUserVerify
	} else {
		subject = common.FailedUserVerify
	}
	if content, err = emailService.GetContent(subject); err != nil {
		return err
	}

	content = strings.Replace(content, "${{username}}", pve.user.Username, 1)
	content = strings.Replace(content, "${{sno}}", pve.user.Sno, 1)
	content = strings.Replace(content, "${{phone}}", pve.user.Phone, 1)
	content = strings.Replace(content, "${{gomeet_url}}", config.Cfg.Section("server").Key("gomeet_url").String(), 1)
	return common.SendEmail([]string{pve.user.Email}, subject, content)
}

func (ave *AdminVerifyEmail) SendEmail(success bool) (err error) {
	/*
	 * 1. 获取邮件模板
	 * 2. 替换邮件模板变量
	 * 3. 发送邮件
	 */
	repo := repository.NewEmailRepository("email")
	emailService := service.NewEmailService(repo)
	var content, subject string
	if success {
		subject = common.PassAdminVerify
	} else {
		subject = common.FailedAdminVerify
	}
	if content, err = emailService.GetContent(subject); err != nil {
		return err
	}

	content = strings.Replace(content, "${{username}}", ave.user.Username, 1)
	if success {
		content = strings.Replace(content, "${{sno}}", ave.user.Sno, 1)
		content = strings.Replace(content, "${{phone}}", ave.user.Phone, 1)
		content = strings.Replace(content, "${{gomeet_url}}", config.Cfg.Section("server").Key("gomeet_url").String(), 1)
	}
	return common.SendEmail([]string{ave.user.Email}, subject, content)
}

func (pve *AppointmentVerifyEmail) SendEmail(success bool) (err error) {
	/*
	 * 1. 获取邮件模板
	 * 2. 替换邮件模板变量
	 * 3. 发送邮件
	 */
	defer func() {
		if err := recover(); err != nil {
			logger.Record("发送邮件奔溃")
		}
	}()
	repo := repository.NewEmailRepository("email")
	emailService := service.NewEmailService(repo)
	var content, subject string
	if success {
		subject = common.PassAppointmentVerify
	} else {
		subject = common.FailedAppointmentVerify
	}
	if content, err = emailService.GetContent(subject); err != nil {
		return err
	}

	content = strings.Replace(content, "${{username}}", pve.user.Username, 1)
	content = strings.Replace(content, "${{theme}}", pve.appointment.Theme, 1)
	if success {
		day := pve.appointment.Day[:4] + "-" + pve.appointment.Day[4:6] + "-" + pve.appointment.Day[6:]
		time := day + " " + pve.appointment.StartTime + "-" + pve.appointment.EndTime
		content = strings.Replace(content, "${{time}}", time, 1)
		content = strings.Replace(content, "${{locate}}", pve.appointment.Locate, 1)
		content = strings.Replace(content, "${{content}}", pve.appointment.Content, 1)
		content = strings.Replace(content, "${{gomeet_url}}", config.Cfg.Section("server").Key("gomeet_url").String(), 1)
	}
	return common.SendEmail([]string{pve.user.Email}, subject, content)
}

func (nme *NotifyMembersEmail) SendEmail(success bool) (err error) {
	/*
	 * 1. 获取邮件模板
	 * 2. 替换邮件模板变量
	 * 3. 发送邮件
	 */
	repo := repository.NewEmailRepository("email")
	emailService := service.NewEmailService(repo)
	var content, subject string
	if success {
		subject = common.PassNotifyMembers
	} else {
		subject = common.FailedNotifyMembers
	}
	if content, err = emailService.GetContent(subject); err != nil {
		return err
	}

	day := nme.appointment.Day[:4] + "-" + nme.appointment.Day[4:6] + "-" + nme.appointment.Day[6:]
	time := day + " " + nme.appointment.StartTime + "-" + nme.appointment.EndTime

	content = strings.Replace(content, "${{time}}", time, 1)
	content = strings.Replace(content, "${{theme}}", nme.appointment.Theme, 1)
	content = strings.Replace(content, "${{creator}}", nme.appointment.CreatorName, 1)
	content = strings.Replace(content, "${{locate}}", nme.appointment.Locate, 1)
	content = strings.Replace(content, "${{content}}", nme.appointment.Content, 1)
	content = strings.Replace(content, "${{gomeet_url}}", config.Cfg.Section("server").Key("gomeet_url").String(), 1)
	for _, v := range nme.user {
		tmp := strings.Replace(content, "${{username}}", v.Username, 1)
		if err = common.SendEmail([]string{v.Email}, subject, tmp); err != nil {
			logger.Record("发送会议参会通知邮件错误", err)
		}
	}
	return nil
}

func (vce *EmailVerifyCodeEmail) SendEmail(bool) (err error) {
	/*
	 * 1. 获取邮件模板
	 * 2. 替换邮件模板变量
	 * 3. 发送邮件
	 */
	// 1. 获取邮件模板
	repo := repository.NewEmailRepository("email")
	emailService := service.NewEmailService(repo)
	var content string
	subject := common.EmailVerifyCode
	if content, err = emailService.GetContent(subject); err != nil {
		return err
	}

	// 2. 替换邮件模板变量
	// 生成随机6位数字，存储到redis
	content = strings.Replace(content, "${{code}}", vce.user.Code, 1)
	content = strings.Replace(content, "${{gomeet_url}}", config.Cfg.Section("server").Key("gomeet_url").String(), 1)
	// 3. 发送邮件
	return common.SendEmail([]string{vce.user.Email}, subject, content)
}
