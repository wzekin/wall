package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"wall/models"
	_ "wall/routers"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "data.db")

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.SetStaticPath("/fonts/roboto", "static")
	beego.AddFuncMap("add", add)
	beego.AddFuncMap("format", models.FormatTime)
}

func main() {
	d := models.Data{Content: "1231", Date: time.Now().AddDate(0, 0, -1)}
	d.Insert()
	beego.Run()
}

func add(a int) int {
	return a + 1
}
