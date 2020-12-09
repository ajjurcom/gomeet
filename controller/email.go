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
		l := strings.Split(pve.appointment.Day, "/")
		time := l[2] + "-" + l[0] + "-" + l[1] + " " + pve.appointment.StartTime + "-" + pve.appointment.EndTime
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

	l := strings.Split(nme.appointment.Day, "/")
	time := l[2] + "-" + l[0] + "-" + l[1] + " " + nme.appointment.StartTime + "-" + nme.appointment.EndTime

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
