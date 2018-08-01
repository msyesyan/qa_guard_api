package main

import (
	"github.com/astaxie/beego/plugins/cors"
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

	if beego.AppConfig.String("runmode") == "dev" {
		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
			AllowOrigins: []string{"http://localhost:*", "http://127.0.0.1:*"},
			AllowCredentials: true,
		}))

		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
}

func main() {
	o := orm.NewOrm()
	o.Using("default")

	beego.Run()
}
