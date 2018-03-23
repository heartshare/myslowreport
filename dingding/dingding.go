package dingding

import (
	"strings"
	"fmt"
	"time"
	"net/http"
	"github.com/ximply/myslowreport/models"
)

func SendTextToDingding(content string, at []string, atAll bool, webHook string) error {
	s := `{
    	"msgtype": "text",
    	"text": {
        	"content": "{{.content}}"
    	},
    	"at": {
        	"atMobiles": [
            	{{.at}}
        	],
        	"isAtAll": {{.atAll}}
    	}
	}
	`
	s = strings.Replace(s, "{{.content}}", content, -1)
	if atAll {
		s = strings.Replace(s, "{{.atAll}}", "true", -1)
	} else {
		s = strings.Replace(s, "{{.atAll}}", "false", -1)
		atList := ""
		for _, a := range at {
			atList += fmt.Sprintf("\"%s\",", a)
		}
		atList = strings.TrimRight(atList, ",")
		s = strings.Replace(s, "{{.at}}", atList, -1)
	}

	retry := models.MyslowReportDdRetry()
	timeout :=models.MyslowReportDdTimeout()
	var ret error = nil
	for i := 0; i < retry; i++ {
		postReq, err := http.NewRequest("POST",
			models.MyslowReportDdWebhook(),
			strings.NewReader(s))
		ret = err
		if err != nil {
			continue
		}

		postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

		client := &http.Client{}
		client.Timeout = time.Second * time.Duration(timeout)
		resp, err := client.Do(postReq)
		defer resp.Body.Close()
		ret = err
		if err != nil {
			continue
		} else {
			break
		}
	}
	return ret
}
