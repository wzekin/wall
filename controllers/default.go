package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"time"
	"wall/models"
)

const APPID string = "wx98695df4fa0decb7"
const APPSECRET string = "610cdfd7e2f61641466e9e442708447a"

type MainController struct {
	beego.Controller
}

//主页面
func (c *MainController) Get() {
	c.Data["json"], _ = models.GetData(true, models.FormatTime(time.Now().AddDate(0, 0, -1)))
	c.ServeJSON()
}

//提交表白信息
func (c *MainController) Post() {
	code := c.GetString("code")
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		APPID, APPSECRET, code)
	res, err := http.Get(url)
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	j := make(map[string]string)
	beego.Informational(string(b))
	err = json.Unmarshal(b, &j)
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	d := models.Data{Content: c.GetString("content"), WxId: j["openid"]}
	if d.CheckExist() {
		c.Ctx.WriteString("今日已提交")
	} else {
		err = d.Insert()
		if err != nil {
			beego.Error(err)
			c.Abort("500")
		}
		c.Ctx.WriteString("提交成功")
	}
}

//登陆
type LoginController struct {
	beego.Controller
}

//登陆页面
func (c *LoginController) Get() {
	c.TplName = "login.html"
}

func (c *LoginController) Prepare() {
	s := c.GetSession("fail")
	if s != nil {
		if s.(int) > 2 {
			c.Abort("403")
		}
	}
}

//登陆提交
func (c *LoginController) Post() {
	if c.GetString("auth") == "goland" && c.GetString("password") == "beego" {
		c.SetSession("login", true)
		c.Ctx.WriteString("<script>alert('登陆成功');window.location.href = '/checked?isToday=true';</script>")
	} else {
		n := c.GetSession("fail")
		var num int
		if n == nil {
			num = 1
		} else {
			num = n.(int) + 1
		}
		c.Ctx.WriteString(fmt.Sprintf("<script>alert('登陆失败,你还有%d次登陆机会');window.location.href = '/login';</script>", 3-num))
		c.SetSession("fail", num)
	}
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
		datas, err = models.GetData(true, models.FormatTime(time.Now()))
		c.Data["title"] = "今日表白通过列表"
		if err != nil {
			beego.Error(err)
			return
		}
	} else {
		datas, err = models.GetData(true, "")
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
		datas, err = models.GetData(false, models.FormatTime(time.Now()))
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["title"] = "今日表白审核"
	} else {
		datas, err = models.GetData(false, "")
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
	c.Ctx.WriteString("<script>alert('提交成功');window.location.href = document.referrer;</script>")
}
