package login

import (
	"github.com/astaxie/beego"
	"ions_zhiliao/utils"
	"ions_zhiliao/models/auth"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"fmt"
)

type LoginController struct {
	beego.Controller
}


func (l *LoginController) Get()  {
	id,base64,err := utils.GetCaptcha()
	if err != nil{
		ret := fmt.Sprintf("登录的get请求，获取验证码错误，错误信息:%v",err)
		logs.Error(ret)
		return
	}

	l.Data["captcha"] = utils.Captcha{Id:id,BS64:base64}
	l.TplName = "login/login.html"

}

func (l *LoginController) Post()  {

	username := l.GetString("username")
	password := l.GetString("password")
	captcha := l.GetString("captcha")
	captcha_id := l.GetString("captcha_id")


	md5_pwd := utils.GetMd5Str(password)

	userinfo := auth.User{}
	o := orm.NewOrm()
	is_exist := o.QueryTable("sys_user").Filter("user_name",username).Filter("password",md5_pwd).Exist()

	o.QueryTable("sys_user").Filter("user_name",username).Filter("password",md5_pwd).One(&userinfo)

	// 验证码校验,需要验证码id和验证码的答案
	is_ok := utils.VerityCaptcha(captcha_id,captcha)
	ret_map := map[string]interface{}{}
	if !is_exist{
		ret := fmt.Sprintf("登录的post请求，用户名密码错误，登录信息：username:%s;pwd:%s",username,md5_pwd)
		logs.Info(ret)
		ret_map["code"] = 10001
		ret_map["msg"] = "用户名或密码错误"
		l.Data["json"] = ret_map
	}else if !is_ok {
		ret := fmt.Sprintf("登录的post请求，验证码错误，验证码信息:%t",is_ok)
		logs.Info(ret)
		ret_map["code"] = 10001
		ret_map["msg"] = "验证码错误"
		l.Data["json"] = ret_map
	}else if userinfo.IsActive == 0 {
		ret := fmt.Sprintf("登录的post请求，该用户已停用，用户名：%s，状态：停用",username)
		logs.Info(ret)
		ret_map["code"] = 10001
		ret_map["msg"] = "该用户已停用，请联系管理员"
		l.Data["json"] = ret_map
	}else {
		ret := fmt.Sprintf("登录的post请求，登录成功，登录信息：username:%s;pwd:%s",username,md5_pwd)
		logs.Info(ret)
		l.SetSession("id",userinfo.Id)
		ret_map["code"] = 200
		ret_map["msg"] = "登录成功"
		l.Data["json"] = ret_map
	}

	l.ServeJSON()


}

func (l *LoginController) ChangeCaptcha()  {

	message := map[string]interface{}{}
	id,base64,err := utils.GetCaptcha()

	if err != nil{  // 说明有错误
		ret := fmt.Sprintf("生成验证码失败,验证码信息：id:%s,错误信息:%v",id,err)
		logs.Error(ret)
		message["msg"] = "生成验证码失败"
		message["Code"] = 404
		l.Data["json"] = message
	}else {
		l.Data["json"] = utils.Captcha{Id:id,BS64:base64,Code:200}
	}
	l.ServeJSON()

}

func (l *LoginController) LogOut()  {

	l.DelSession("id")
	l.Redirect(beego.URLFor("LoginController.Get"),302)
	
}



