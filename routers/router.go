package routers

import (
	"github.com/astaxie/beego"
	"wall/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/checked", &controllers.CheckController{}, "get:Checked")
	beego.Router("/not_checked", &controllers.CheckController{}, "get:NotChecked")
	beego.Router("/pass", &controllers.CheckController{}, "post:Pass")
	beego.Router("/login", &controllers.LoginController{})
}
