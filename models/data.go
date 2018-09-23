package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Data struct {
	Id      int `orm:"auto"`
	WxId    string
	Date    string `orm:"size(10)"`
	Content string `orm:"type(text)"`
	Check   bool
}

func init() {
	orm.RegisterModel(new(Data))
	orm.Debug = true
}

func GetData(checked bool, time string) (datas []Data, err error) {
	o := orm.NewOrm()
	hasTime := time != ""
	if !hasTime {
		_, err = o.QueryTable("data").Filter("check", checked).All(&datas)
	} else {
		_, err = o.QueryTable("data").Filter("check", checked).Filter("date", time).All(&datas)
	}
	if err != nil {
		if err.Error() == "<QuerySeter> no row found" {
			return nil, nil
		}
	}
	beego.Informational(datas)
	return
}

func (d *Data) Insert() (err error) {
	o := orm.NewOrm()
	if d.Date == "" {
		d.Date = FormatTime(time.Now())
	}
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
	beego.Informational(d)
	_, err = o.Update(d)
	return
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02")
}

func (d Data) CheckExist() bool {
	if d.Date == "" {
		d.Date = FormatTime(time.Now())
	}
	o := orm.NewOrm()
	err := o.Read(&d, "date", "wx_id")
	if err != nil {
		return true
	}
	if d.Id != 0 {
		return true
	} else {
		return false
	}
}
