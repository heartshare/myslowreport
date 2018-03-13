package email


import (
	"net/smtp"
	"strings"
	"github.com/go-gomail/gomail"
	"strconv"
	"fmt"
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

func sendEmailWithAdditionImpl(mailUserAlias string,
	mailUserName string,
	mailPassWord string,
	mailHost string,
	mailPort string,
	mailTo string,
	mailSubject string,
	reportFile string,
	body string,
	mailType string) error {

	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress("kmmp@ktvme.com", mailUserAlias)},
		"To":      strings.Split(mailTo, ";"),
		"Subject": {fmt.Sprintf(mailSubject)},
	})
	m.SetBody(fmt.Sprintf("text/%s;", mailType), body)
	m.Attach(reportFile)
	port, _ := strconv.Atoi(mailPort)
	d := gomail.NewDialer(mailHost, port, mailUserName, mailPassWord)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmailWithAddition(mailUserAlias string,
	mailUserName string,
	mailPassWord string,
	mailHost string,
	mailPort string,
	mailTo string,
	mailSubject string,
	reportFile string,
	body string,
	mailType string) (int, error) {

	if len(mailUserAlias) < 1 || len(mailPassWord) < 1 || len(mailHost) < 1 || len(mailPort) < 1 || len(mailTo) < 1 {
		return -1, nil
	}

	return 0, sendEmailWithAdditionImpl(mailUserAlias, mailUserName,
		mailPassWord, mailHost, mailPort, mailTo, mailSubject, reportFile, body, mailType)
}
