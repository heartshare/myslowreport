package main

import (
	_ "github.com/ximply/myslowreport/routers"
	"github.com/ximply/myslowreport/models"
	"github.com/ximply/myslowreport/crons"
	"github.com/astaxie/beego"
)

func init() {
	models.Init()
}

func main() {
	crons.Init()
	beego.Run()
}
