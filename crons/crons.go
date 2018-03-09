package crons

import (
	_ "github.com/astaxie/beego"
	"github.com/robfig/cron"
	_ "github.com/ximply/myslowreport/models"
	//"time"
)

func syncMysqlSlowlog() {

}

func sendMysqlSlowlogReport() {

}

func Init() {
	c := cron.New()
	//c.AddFunc(models.SyncMysqlSlowlogSpec(), syncMysqlSlowlog)
	//c.AddFunc(models.SendMysqlSlowlogReportSpec(), sendMysqlSlowlogReport)
	c.Start()
}

