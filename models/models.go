package models

import (
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"os"
	"path"
)

const (
	_DB_NAME        = "beeblog"
	_SQLITLE_DRIVER = "mysql"
)

type Category struct {
	Id              int64
	Title           string
	Create          time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Labels          string
	Category        string
	Content         string    `orm:"size(500)"`
	Attachment      string
	Create          time.Time `orm:"index"`
	Update          time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}

type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(500)"`
	Create  time.Time `orm:"index"`
}

func RegisterDB() {
	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	orm.RegisterDriver(_SQLITLE_DRIVER, orm.DRMySQL)
	orm.RegisterDataBase("default", _SQLITLE_DRIVER, "root:123456@/"+_DB_NAME+"?charset=utf8")
}

func AddCategory(name string) error {
	o := orm.NewOrm()
	cate := &Category{Title: name, Create: time.Now(), TopicTime: time.Now()}
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}
	_, err = o.Insert(cate)
	return err
}

func FindAllCategory() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

func DeleteCategoryById(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	return err
}

func AddTopic(title, category, label, content, attachment string) error {
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"
	o := orm.NewOrm()
	topic := &Topic{
		Title:      title,
		Category:   category,
		Content:    content,
		Labels:     label,
		Attachment: attachment,
		Create:     time.Now(),
		Update:     time.Now(),
		ReplyTime:  time.Now()}
	_, err := o.Insert(topic)
	if err != nil {
		return err
	}
	cate := new(Category)
	qs := orm.NewOrm().QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return err
}

func GetTopic(id string) (*Topic, error) {
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tid).One(topic)
	if err != nil {
		return nil, err
	}
	topic.Views ++
	topic.Update = time.Now()
	_, err = o.Update(topic)
	topic.Labels = strings.Replace(strings.Replace(topic.Labels, "#", " ", -1), "$", "", -1)
	return topic, err
}

func ModifyTopic(tid, title, category, label, content, attachment string) error {
	topic, err := GetTopic(tid)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	//这里是判断是否存在记录
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"
	var oldCate, oldAttachment string
	if o.Read(topic) == nil {
		oldCate = topic.Category
		oldAttachment = topic.Attachment
		topic.Title = title
		topic.Labels = label
		topic.Attachment = attachment
		topic.Content = content
		topic.Category = category
		topic.Update = time.Now()
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	if attachment != oldAttachment && len(oldAttachment) > 0 {
		os.Remove(path.Join("static/upload/", oldAttachment))
	}
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", oldCate).One(cate)
	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return err
}

func DeleteTopicById(tid string) error {
	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	topic := &Topic{Id: id}
	var oldCate string
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount --
			_, err = o.Update(cate)
		}
	}
	return err
}

func FindAllTopic(label, cate string, isDesc bool) ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	var err error
	if isDesc {
		if len(cate) > 0 {
			qs = qs.Filter("category", cate)
		}
		if len(label) > 0 {
			qs = qs.Filter("labels__contains", "$"+label+"#")
		}
		qs = qs.OrderBy("-create")
	}
	_, err = qs.All(&topics)
	return topics, err
}

func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	reply := &Comment{Tid: tidNum, Name: nickname, Content: content, Create: time.Now()}
	_, err = orm.NewOrm().Insert(reply)
	return err
}

func FindRepliesByTid(tid string) ([]*Comment, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	replies := make([]*Comment, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).OrderBy("-create").All(&replies)
	return replies, err
}

func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}
	repl := &Comment{Id: ridNum}
	_, err = orm.NewOrm().Delete(repl)
	return err
}
