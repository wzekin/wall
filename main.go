package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	_ "wall/routers"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "data.db")

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.SetStaticPath("/fonts/roboto", "static")
	beego.AddFuncMap("add", add)
}

func main() {
	beego.Run()
}

func add(a int) int {
	return a + 1
}
