package controllers

import (
	"github.com/astaxie/beego"
	"net/url"
	"os"
	"io"
)

type AttachController struct {
	beego.Controller
}

func (a *AttachController) Get() {
	filePath, err := url.QueryUnescape(a.Ctx.Request.RequestURI[1:])
	if err !=nil{
		a.Ctx.WriteString(err.Error())
		return
	}
	f,err :=os.Open(filePath)
	if err !=nil{
		a.Ctx.WriteString(err.Error())
		return
	}
	defer f.Close()
	_,err = io.Copy(a.Ctx.ResponseWriter,f)
	if err !=nil{
		a.Ctx.WriteString(err.Error())
		return
	}
}
