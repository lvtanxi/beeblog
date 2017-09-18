package controllers

import (
	"github.com/astaxie/beego"
	"beeblog/models"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	c.Data["IsCategory"] = true
	op :=c.Input().Get("op")
	beego.Info(op)
	switch op {
	case "add":
		name:=c.Input().Get("name")
		if len(name)==0{
			break
		}
		err :=models.AddCategory(name)
		if err!=nil {
			beego.Error(err)
		}
		c.Redirect("/category",301)
		return
	case "del":
		id:=c.Input().Get("id")
		if len(id)==0{
			break
		}
		err :=models.DeleteCategoryById(id)
		if err!=nil {
			beego.Error(err)
		}
		c.Redirect("/category",301)
		return
	default:
	}

	var err error
	c.Data["Categories"] ,err =models.FindAllCategory()
	if err!=nil {
		beego.Error(err)
	}
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.TplName = "category.html"
}
