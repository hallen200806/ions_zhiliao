package user

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ions_zhiliao/models/auth"
	"ions_zhiliao/utils"
	"strings"
)

type MyCenterController struct {
	beego.Controller
}

func (m *MyCenterController) Get()  {
	id := m.GetSession("id")

	o := orm.NewOrm()
	qs := o.QueryTable("sys_user")
	user := auth.User{}
	qs.Filter("id",id).One(&user)
	m.Data["user"] = user
	m.TplName = "user/my_center_edit.html"

}

func (m *MyCenterController) Post()  {
	uid,_ := m.GetInt("uid")

	username := m.GetString("username")
	old_pwd := m.GetString("old_pwd")
	new_pwd := m.GetString("new_pwd")
	age,_ := m.GetInt("age")
	gender,_ := m.GetInt("gender")
	phone,_ := m.GetInt64("phone")
	addr := m.GetString("addr")
	is_active,_ := m.GetInt("is_active")

	addr_new := strings.Trim(addr," ")

	o := orm.NewOrm()
	qs := o.QueryTable("sys_user")
	user := auth.User{}
	qs.Filter("id",uid).One(&user)

	old_pwd_md5 := utils.GetMd5Str(old_pwd)

	message_map := map[string]interface{}{}
	if old_pwd_md5 != user.Password {
		message_map["code"] = 10001
		message_map["msg"] = "原密码错误"
	}else {

		qs.Filter("id",uid).Update(orm.Params{
			"username":username,
			"password":utils.GetMd5Str(new_pwd),
			"age":age,
			"gender":gender,
			"phone":phone,
			"addr":addr_new,
			"is_active":is_active,
		})
		message_map["code"] = 200
		message_map["msg"] = "修改成功"

	}


	m.Data["json"] = message_map
	m.ServeJSON()



}
