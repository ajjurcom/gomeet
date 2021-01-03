package common

import (
	"errors"
	"gopkg.in/gomail.v2"
	"strconv"
)

const (
	PassUserVerify = "用户审核通过"
	FailedUserVerify = "用户审核不通过"
	PassAdminVerify = "管理员审核通过"
	FailedAdminVerify = "管理员审核不通过"
	PassAppointmentVerify = "会议审核通过"
	FailedAppointmentVerify = "会议审核不通过"
	PassNotifyMembers = "会议参会通知"
	FailedNotifyMembers = "会议退订通知"
	EmailVerifyCode = "注册验证码(邮箱)"
)

func SendEmail(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": "81679140@qq.com",
		"pass": "rssiujminkdmbhhc",
		"host": "smtp.qq.com",
		"port": "465",
	}

	port, err := strconv.Atoi(mailConn["port"])
	if err != nil {
		return errors.New("端口号必须为数字字符串")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "GoMEET系统"))
	m.SetHeader("To", mailTo...)		// 收件人
	m.SetHeader("Subject", subject) 	// 邮件主题
	m.SetBody("text/html", body)	// 邮件正文
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	return d.DialAndSend(m)
}
