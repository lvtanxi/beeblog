package controllers

import (
	"github.com/astaxie/beego"
	"beeblog/models"
)

type MainController struct {
	beego.Controller
}

type u struct {
	Name string
	Age  int
	Sex  int
}

func (c *MainController) Get() {
	c.Data["IsHome"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	var err error
	c.Data["Topics"],err=models.FindAllTopic(c.Input().Get("label"),c.Input().Get("cate"),true)
	if err != nil {
		beego.Error(err)
	}
	categories,err:=models.FindAllCategory()
	if err != nil {
		beego.Error(err)
	}
	c.Data["Categories"]=categories
	c.TplName = "home.html"
}
