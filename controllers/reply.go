package controllers

import (
	"github.com/astaxie/beego"
	"beeblog/models"
)

type ReplyController struct {
	beego.Controller
}

func (r *ReplyController) Add() {
	tid := r.Input().Get("tid")
	err := models.AddReply(tid, r.Input().Get("nickname"), r.Input().Get("content"))
	if err != nil {
		beego.Error(err)
	}
	r.Redirect("/topic/view/"+tid, 302)
}

func (r *ReplyController) Delete() {
	if !checkAccount(r.Ctx) {
		r.Redirect("/login", 302)
		return
	}
	err := models.DeleteReply(r.Ctx.Input.Params()["0"])
	if err != nil {
		beego.Error(err)
	}
	r.Redirect("/topic/view/"+r.Ctx.Input.Params()["1"], 302)
}
