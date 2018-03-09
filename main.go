package main

import (
	_ "github.com/ximply/myslowreport/routers"
	"github.com/ximply/myslowreport/models"
	"github.com/ximply/myslowreport/crons"
	"github.com/astaxie/beego"
	"github.com/ximply/myslowreport/rsync"
)

func init() {
	models.Init()
}

func main() {
	crons.Init()
	rsync.SyncMysqlSlowlogFile("192.168.10.121", "myslow", "localhost-slow.log", "/home/www/log/")
	beego.Run()
}
