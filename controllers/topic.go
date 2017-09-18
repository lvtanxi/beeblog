package controllers

import (
	"github.com/astaxie/beego"
	"beeblog/models"
	"strings"
	"path"
)

type TopicController struct {
	beego.Controller
}

func (t *TopicController) Get() {
	t.Data["IsTopic"] = true
	t.Data["IsLogin"] = checkAccount(t.Ctx)
	var err error
	t.Data["Topics"], err = models.FindAllTopic("", "", false)
	if err != nil {
		beego.Error(err)
	}
	t.TplName = "topic.html"
}

func (t *TopicController) Post() {
	if !checkAccount(t.Ctx) {
		t.Redirect("/login", 302)
		return
	}
	title := t.Input().Get("title")
	content := t.Input().Get("content")
	tid := t.Input().Get("tid")
	category := t.Input().Get("category")
	label := t.Input().Get("label")
	//获取附件
	f, fh, err := t.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var fileName string
	if fh != nil {
		fileName = fh.Filename
		beego.Info("this is upload",fileName)
		defer f.Close()
		err = t.SaveToFile("attachment", path.Join("static/upload/",fileName))
		if err != nil {
			beego.Error(err)
		}
	}
	if len(tid) == 0 {
		err = models.AddTopic(title, category, label, content,fileName)
	} else {
		err = models.ModifyTopic(tid, title, category, label, content,fileName)
	}

	if err != nil {
		beego.Error(err)
	}
	t.Redirect("/topic", 302)
}

func (t *TopicController) Add() {
	t.Data["IsTopic"] = true
	t.TplName = "topic_add.html"
}

func (t *TopicController) View() {
	tid := t.Ctx.Input.Params()["0"]
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		t.Redirect("/", 302)
		return
	}
	t.Data["Topic"] = topic
	t.Data["Labels"] = strings.Split(topic.Labels, " ")
	t.Data["IsTopic"] = true
	t.Data["Tid"] = tid

	replies, err := models.FindRepliesByTid(tid)
	if err != nil {
		beego.Error(err)
		return
	}
	t.Data["Replies"] = replies
	t.TplName = "topic_view.html"
}

func (t *TopicController) Modify() {
	if !checkAccount(t.Ctx) {
		t.Redirect("/login", 302)
		return
	}
	tid := t.Input().Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		t.Redirect("/", 302)
		return
	}
	t.Data["IsTopic"] = true
	t.Data["Topic"] = topic
	t.Data["Tid"] = tid
	t.TplName = "topic_add.html"
}

func (t *TopicController) Delete() {
	if !checkAccount(t.Ctx) {
		t.Redirect("/login", 302)
		return
	}
	err := models.DeleteTopicById(t.Ctx.Input.Params()["0"])
	if err != nil {
		beego.Error(err)
	}
	t.Redirect("/topic", 302)

}
