package routers

import (
	"bee-go-vue/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.UserController{}, "Get:ShowIndex")
	beego.Router("/users/:keyword/:page", &controllers.UserController{}, "Get:GetUsers")
	beego.Router("/users/:page", &controllers.UserController{}, "Get:GetAllUsers")
	beego.Router("/user", &controllers.UserController{}, "Post:PostUser")
	beego.Router("/user/:uid/:org/edit", &controllers.UserController{}, "Put:PutUser")
	beego.Router("/user/:uid/:org", &controllers.UserController{}, "Delete:DeleteUser")
	beego.Router("/user/:uid/:org/password", &controllers.UserController{}, "Put:ChangeUserPassword")
}
