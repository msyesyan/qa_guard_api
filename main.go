package main

import (
	_ "github.com/lib/pq"
	"github.com/astaxie/beego/orm"
	_ "qa_guard_api/routers"

	"github.com/astaxie/beego"
)

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	err := orm.RegisterDataBase("default", "postgres", beego.AppConfig.String("dbconn"))

	if err != nil {
		panic(err)
	}

	orm.RunCommand()
}

func main() {
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

	o := orm.NewOrm()
	o.Using("default")

	beego.Run()
}
