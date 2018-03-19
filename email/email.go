package email


import (
	"strings"
	"github.com/go-gomail/gomail"
	"strconv"
	"fmt"
	"github.com/ximply/myslowreport/utils"
)

func sendEmailWithAdditionImpl(
	mailUserName string,
	mailPassWord string,
	mailHost string,
	mailPort string,
	mailFrom string,
	mailFromAlias string,
	mailTo string,
	mailCc string,
	mailSubject string,
	reportFile string,
	body string,
	mailType string) error {

	m := gomail.NewMessage()
	if len(mailCc) > 1 {
		m.SetHeaders(map[string][]string{
			"From":    {m.FormatAddress(mailFrom, mailFromAlias)},
			"To":      strings.Split(mailTo, ";"),
			"Cc":      strings.Split(mailCc, ";"),
			"Subject": {fmt.Sprintf(mailSubject)},
		})
	} else {
		m.SetHeaders(map[string][]string{
			"From":    {m.FormatAddress(mailFrom, mailFromAlias)},
			"To":      strings.Split(mailTo, ";"),
			"Subject": {fmt.Sprintf(mailSubject)},
		})
	}

	m.SetBody(fmt.Sprintf("text/%s;", mailType), body)
	fileList := strings.Split(reportFile, ";")
	for _, f := range fileList {
		if utils.FileExists(f) {
			m.Attach(f)
		}
	}
	port, _ := strconv.Atoi(mailPort)
	d := gomail.NewDialer(mailHost, port, mailUserName, mailPassWord)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmailWithAddition(
	mailUserName string,
	mailPassWord string,
	mailHost string,
	mailPort string,
	mailFrom string,
	mailFromAlias string,
	mailTo string,
	mailCc string,
	mailSubject string,
	reportFile string,
	body string,
	mailType string) (int, error) {

	if len(mailPassWord) < 1 ||
		len(mailHost) < 1 ||
	 	len(mailPort) < 1 ||
	 	len(mailFrom) < 1 ||
	 	len(mailTo) < 1 {
		return -1, nil
	}

	return 0, sendEmailWithAdditionImpl(
		mailUserName,
		mailPassWord,
		mailHost,
		mailPort,
		mailFrom,
		mailFromAlias,
		mailTo,
		mailCc,
		mailSubject,
		reportFile,
		body,
		mailType)
}
