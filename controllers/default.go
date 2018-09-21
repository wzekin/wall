package controllers

import (
	"github.com/astaxie/beego"
	"time"
	"wall/models"
)

type MainController struct {
	beego.Controller
}

//主页面
func (c *MainController) Get() {
	c.Data["msgs"], _ = models.GetData(true, time.Now().AddDate(0, 0, -1))
	c.TplName = "wall.html"
}

//提交表白信息
func (c *MainController) Post() {
	d := models.Data{Content: c.GetString("content")}
	err := d.Insert()
	if err != nil {
		beego.Error(err)
	}
	c.Ctx.WriteString("<script>alert('提交成功');window.location.href = '/';</script>")
	return
}

//后台
type UnderController struct {
	beego.Controller
}

//查询表白审核
//所有表白审核 所有通过表白 所有今天表白审核 所有今天通过表白
func (c *UnderController) Post() {
	checked, err := c.GetBool("checked")
	if err != nil {
		beego.Error(err)
	}
	t, _ := time.Parse("2006-01-02", c.GetString("time"))
	c.Data["json"], err = models.GetData(checked, t)
	if err != nil {
		return
	}
	c.ServeJSON()
}

//登陆
type LoginController struct {
	beego.Controller
}

//登陆页面
func (c *LoginController) Get() {
	c.TplName = "login.html"
}

//登陆提交
func (c *LoginController) Post() {
	if c.GetString("auth") == "goland" && c.GetString("password") == "beego" {
		c.SetSession("login", true)
		c.Ctx.WriteString("<script>alert('登陆成功');window.location.href = '/checked?isToday=true';</script>")
	}
	c.Ctx.WriteString("<script>alert('登陆失败');window.location.href = '/login';</script>")
}

//后端渲染
type CheckController struct {
	beego.Controller
}

func (c *CheckController) Prepare() {
	if c.GetSession("login") == nil {
		c.Abort("403")
	}
}

func (c *CheckController) Checked() {
	b, err := c.GetBool("isToday")
	if err != nil {
		beego.Error(err)
		return
	}
	var datas []models.Data
	if b {
		datas, err = models.GetData(true, time.Now())
		c.Data["title"] = "今日表白通过列表"
		if err != nil {
			beego.Error(err)
			return
		}
	} else {
		datas, err = models.GetData(true, time.Time{})
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["title"] = "所有表白通过列表"
	}
	c.Data["msgs"] = datas
	c.TplName = "checked.html"
}

func (c *CheckController) NotChecked() {
	b, err := c.GetBool("isToday")
	if err != nil {
		beego.Error(err)
		return
	}
	var datas []models.Data
	if b {
		datas, err = models.GetData(false, time.Now())
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["title"] = "今日表白审核"
	} else {
		datas, err = models.GetData(false, time.Time{})
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["title"] = "所有表白审核"
	}
	c.Data["msgs"] = datas
	c.TplName = "not_checked.html"
}

func (c *CheckController) Pass() {
	id, err := c.GetInt("id")
	if err != nil {
		beego.Error(err)
		return
	}
	d := models.Data{Id: id}
	err = d.Update()
	if err != nil {
		beego.Error()
		return
	}
	c.Ctx.WriteString("<script>alert('提交成功')</script>")
}
