package email


import (
	"net/smtp"
	"strings"
)

func sendMailImpl(user, password, host, port, from, to, subject, body, mailType string) error{
	auth := smtp.PlainAuth("", user, password, host)
	var contentType string
	if mailType == "html" {
		contentType = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + from + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)

	err := smtp.SendMail(host + ":" + port, auth, user, strings.Split(to, ";"), msg)
	return err
}

func SendEmail(mailUserAlias string,
		mailUserName string,
		mailPassWord string,
		mailHost string,
		mailPort string,
		mailTo string,
		mailSubject string,
		report string,
		mailType string) (int, error) {

	if len(mailUserAlias) < 1 || len(mailPassWord) < 1 || len(mailHost) < 1 || len(mailPort) < 1 || len(mailTo) < 1 {
		return -1, nil
	}

	return 0, sendMailImpl(mailUserName,
		mailPassWord, mailHost, mailPort, mailUserAlias, mailTo, mailSubject, report, mailType)
}

