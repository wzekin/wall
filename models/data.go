package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Data struct {
	Id      int       `orm:"auto"`
	Content string    `orm:"type(text)"`
	Date    time.Time `orm:"auto_now_add;type(date)"`
	Check   bool
}

func init() {
	orm.RegisterModel(new(Data))
	orm.Debug = true
}

func GetData(checked bool, t time.Time) (datas []Data, err error) {
	o := orm.NewOrm()
	hasTime := FormatTime(t) != FormatTime(time.Time{})
	if !hasTime {
		_, err = o.QueryTable("data").Filter("check", checked).All(&datas)
	} else {
		_, err = o.QueryTable("data").Filter("check", checked).Filter("date", t).All(&datas)
	}
	return
}

func (d *Data) Insert() (err error) {
	o := orm.NewOrm()
	_, err = o.Insert(d)
	return
}

func (d *Data) Update() (err error) {
	o := orm.NewOrm()
	if d.Id <= 0 {
		return fmt.Errorf("非法Id：%d", d.Id)
	}
	err = o.Read(d)
	if err != nil {
		return err
	}
	d.Check = !d.Check
	beego.Informational(d.Id)
	beego.Informational(d.Check)
	_, err = o.Update(d)
	return
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02")
}
