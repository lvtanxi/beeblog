package routers

import (
	"beeblog/controllers"
	"github.com/astaxie/beego"
	"os"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/topic", &controllers.TopicController{})
	//自动路由,必须以Controller结尾
	beego.AutoRouter(&controllers.TopicController{})

	beego.Router("/reply", &controllers.ReplyController{})
	beego.AutoRouter(&controllers.ReplyController{})

	//创建附件文件目录
	os.Mkdir("/static/upload/",os.ModePerm)
	//作为单独一个控制器来处理
	beego.Router("/attachment/:all",&controllers.AttachController{})

}
