package main

import (
	_ "ions_zhiliao/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "ions_zhiliao/models/auth"
	_ "ions_zhiliao/models/my_center"
	_ "ions_zhiliao/models/caiwu"
	_ "ions_zhiliao/models/news"
	"ions_zhiliao/utils"
	"github.com/astaxie/beego/logs"
	"fmt"
)

func init()  {
	username := beego.AppConfig.String("username")
	pwd := beego.AppConfig.String("pwd")
	host := beego.AppConfig.String("host")
	port := beego.AppConfig.String("port")
	db := beego.AppConfig.String("db")

	// username:pwd@tcp(ip:port)/db?charset=utf8&loc=Local
	dataSource := username + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + db + "?charset=utf8&loc=Local"
	orm.RegisterDriver("mysql",orm.DRMySQL)
	orm.RegisterDataBase("default","mysql",dataSource)

	ret := fmt.Sprintf("host:%s|port:%s|db:%s",host,port,db)
	logs.Info(ret)
}


func main() {
	// 数据库命令行迁移
	orm.RunCommand()
	// 开启session
	beego.BConfig.WebConfig.Session.SessionOn = true

	// 未登录请求拦截
	beego.InsertFilter("/main/*",beego.BeforeRouter,utils.LoginFilter)

	// 日志
	logs.SetLogger(logs.AdapterMultiFile,`{"filename":"logs/ions.log","separate":["error","info"]}`)
	logs.SetLogFuncCallDepth(3)
	beego.SetStaticPath("/upload","upload")
	// 打印sql
	//orm.Debug = true
	beego.Run()
}

