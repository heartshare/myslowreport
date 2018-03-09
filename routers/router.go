package routers

import (
	"github.com/ximply/myslowreport/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
