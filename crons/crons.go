package crons

import (
	_ "github.com/astaxie/beego"
	"github.com/robfig/cron"
	"github.com/ximply/myslowreport/models"
	"time"
)

func summaryReport() {
	models.Get(time.Now(), time.Now())
}

func Init() {
	c := cron.New()
	c.AddFunc(models.MyslowReportSummaryspec(), summaryReport)
	c.Start()
}

