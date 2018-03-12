package main

import (
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
