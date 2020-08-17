package user

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ions_zhiliao/models/auth"
	"math"
	"ions_zhiliao/utils"
	"strconv"
	"strings"
	"github.com/astaxie/beego/logs"
	"fmt"
)

type UserController struct {
	beego.Controller
}

func (u *UserController) List()  {

	//o := orm.NewOrm()
	//
	//qs := o.QueryTable("sys_user")
	//
	//users := []user.User{}
	//// 每页显示的条数
	//pagePerNum := 2
	//// 当前页
	//currentPage,err := u.GetInt("page")
	//if err != nil {   // 说明没有获取到当前页
	//	currentPage = 1
	//}
	//offsetNum := pagePerNum * (currentPage - 1)
	//
	//// 总数
	//count,_ := qs.Filter("is_delete",0).Count()
	//// 总页数
	//countPage := int(math.Ceil(float64(count) / float64(pagePerNum)))
	//
	//qs.Filter("is_delete",0).Limit(pagePerNum).Offset(offsetNum).All(&users)
	//
	///*
	//分页逻辑：
	//	当前第几页     offset     limit
	//	1              0          2           2 *  （1 - 1）
	//    2              2          2			  2 *  （2 -1）
	//	3              4		  2			  2 *  （3 -1）
	//
	//										 limitNum * (currentPage - 1)
	// */
	//prePage := 1
	//if currentPage == 1{
	//	 prePage = currentPage
	//}else if currentPage > 1{
	//	prePage = currentPage -1
	//}
	//
	//nextPage := 1
	//if currentPage < countPage {
	//	nextPage = currentPage + 1
	//}else if currentPage >= countPage {
	//	nextPage = currentPage
	//}


	o := orm.NewOrm()

	qs := o.QueryTable("sys_user")

	users := []auth.User{}
	// 每页显示的条数
	pagePerNum := 8
	// 当前页
	currentPage,err := u.GetInt("page")

	offsetNum := pagePerNum * (currentPage - 1)


	kw := u.GetString("kw")
	var count int64 = 0

	ret := fmt.Sprintf("当前页;%d,查询条件：%s",currentPage,kw)
	logs.Info(ret)
	if kw != ""{   // 有查询条件的
		// 总数
		count,_ = qs.Filter("is_delete",0).Filter("user_name__contains",kw).Count()
		qs.Filter("is_delete",0).Filter("user_name__contains",kw).Limit(pagePerNum).Offset(offsetNum).All(&users)
	}else {
		count,_ = qs.Filter("is_delete",0).Count()
		qs.Filter("is_delete",0).Limit(pagePerNum).Offset(offsetNum).All(&users)

	}
	if err != nil {   // 说明没有获取到当前页
		currentPage = 1
	}



	// 总页数
	countPage := int(math.Ceil(float64(count) / float64(pagePerNum)))



	prePage := 1
	if currentPage == 1{
		prePage = currentPage
	}else if currentPage > 1{
		prePage = currentPage -1
	}

	nextPage := 1
	if currentPage < countPage {
		nextPage = currentPage + 1
	}else if currentPage >= countPage {
		nextPage = currentPage
	}


	page_map := utils.Paginator(currentPage,pagePerNum,count)

	u.Data["users"] = users
	u.Data["prePage"] =prePage
	u.Data["nextPage"] = nextPage
	u.Data["currentPage"] = currentPage
	u.Data["countPage"] = countPage
	u.Data["count"] = count
	u.Data["page_map"] = page_map
	u.Data["kw"] = kw
	u.TplName = "user/user_list.html"

}

func (u *UserController) ToAdd()  {
	u.TplName = "user/salary_slip_import.html"

}

func (u *UserController) DoAdd()  {
	username := u.GetString("username")
	password := u.GetString("password")
	age,_ := u.GetInt("age")
	gender,_ := u.GetInt("gender")
	phone := u.GetString("phone")
	addr := u.GetString("addr")
	is_active,_ := u.GetInt("is_active")

	new_password := utils.GetMd5Str(password)
	phone_int64,_ := strconv.ParseInt(phone,10,64)
	o := orm.NewOrm()
	user_data := auth.User{UserName:username,Password:new_password,Age:age,Gender:gender,Phone:phone_int64,Addr:addr,IsActive:is_active}
	_,err := o.Insert(&user_data)


	message_map := map[string]interface{}{}
	if err != nil {  //说明插入数据有问题
		ret1 := fmt.Sprintf("插入数据信息：username:%s|md5_password:%s|age:%d|gender:%d|phone:%s|" +
			"addr:%s;is_active:%d",username,new_password,age,gender,phone,addr,is_active)
		ret := fmt.Sprintf("添加数据出错,错误信息:%v",err)
		logs.Error(ret1)
		logs.Error(ret)
		message_map["code"] = 10001
		message_map["msg"] = "添加数据出错，请重新添加"
		u.Data["json"] = message_map
	}else {
		ret1 := fmt.Sprintf("插入数据成功，数据信息：username:%s|md5_password:%s|age:%d|gender:%d|phone:%s|" +
			"addr:%s;is_active:%d",username,new_password,age,gender,phone,addr,is_active)
		logs.Info(ret1)
		message_map["code"] = 200
		message_map["msg"] = "添加成功"
		u.Data["json"] = message_map
	}

	u.ServeJSON()



}

func (u *UserController) IsActive()  {
	is_active,_ := u.GetInt("is_active_val")
	id,_ := u.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_user").Filter("id",id)

	message_map := map[string]interface{}{}
	if is_active == 1 {
		qs.Update(orm.Params{
			"is_active": 0,
		})
		ret := fmt.Sprintf("用户id:%d,停用成功",id)
		logs.Info(ret)
		message_map["msg"] = "停用成功"
	}else if is_active == 0 {
		qs.Update(orm.Params{
			"is_active": 1,
		})
		ret := fmt.Sprintf("用户id:%d,启用成功",id)
		logs.Info(ret)
		message_map["msg"] = "启用成功"
	}

	u.Data["json"] = message_map
	u.ServeJSON()
}

func (u *UserController) Delete()  {
	id,_ := u.GetInt("id")

	o := orm.NewOrm()
	o.QueryTable("sys_user").Filter("id",id).Update(orm.Params{
		"is_delete":1,
	})
	ret := fmt.Sprintf("用户id:%d,删除成功",id)
	logs.Info(ret)
	u.Redirect(beego.URLFor("UserController.List"),302)

}

func (u *UserController) ResetPassword()  {
	id,_ := u.GetInt("id")

	o := orm.NewOrm()

	new_pwd := utils.GetMd5Str("123456")
	o.QueryTable("sys_user").Filter("id",id).Update(orm.Params{
		"password":new_pwd,
	})
	ret := fmt.Sprintf("用户id:%d,重置密码成功",id)
	logs.Info(ret)
	u.Redirect(beego.URLFor("UserController.List"),302)

	
}

func (u *UserController) ToUpdate()  {
	id,_ := u.GetInt("id")
	o := orm.NewOrm()
	user_data := auth.User{}
	o.QueryTable("sys_user").Filter("id",id).One(&user_data)
	u.Data["user"] = user_data
	ret := fmt.Sprintf("用户信息修改，用户id:%d",id)
	logs.Info(ret)
	u.TplName = "user/user_edit.html"
}

func (u *UserController) DoUpdate()  {
	uid,_ := u.GetInt("uid")
	username := u.GetString("username")
	password := u.GetString("password")
	age,_ := u.GetInt("age")
	gender,_ := u.GetInt("gender")
	phone := u.GetString("phone")
	addr := u.GetString("addr")
	is_active,_ := u.GetInt("is_active")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_user").Filter("id",uid)

	message_map := map[string]interface{}{}
	new_pwd := utils.GetMd5Str(password)
	if password == ""{
		_,err := qs.Update(orm.Params{
			"username":username,
			"age":age,
			"gender":gender,
			"phone":phone,
			"addr":addr,
			"is_active":is_active,
		})
		if err != nil {
			ret := fmt.Sprintf("更新失败，用户id:%d",uid)
			logs.Error(ret)
			message_map["code"] = 10001
			message_map["msg"] = "更新失败"
		}else {
			ret := fmt.Sprintf("更新成功，用户id:%d",uid)
			logs.Info(ret)
			message_map["code"] = 200
			message_map["msg"] = "更新成功"
		}
	}else {
		_,err := qs.Update(orm.Params{
			"username":username,
			"password":new_pwd,
			"age":age,
			"gender":gender,
			"phone":phone,
			"addr":addr,
			"is_active":is_active,
		})

		if err != nil {

			ret := fmt.Sprintf("更新失败，用户id:%d",uid)
			logs.Error(ret)
			message_map["code"] = 10001
			message_map["msg"] = "更新失败"
		}else {
			ret := fmt.Sprintf("更新成功，用户id:%d",uid)
			logs.Info(ret)
			message_map["code"] = 200
			message_map["msg"] = "更新成功"
		}
	}



	u.Data["json"] = message_map
	u.ServeJSON()



}

func (u *UserController) MuliDelete()  {

	ids := u.GetString("ids")
	//"3,7,8"
	new_ids := ids[1:len(ids)-1]
	id_arr := strings.Split(new_ids,",")

	o := orm.NewOrm()
	qs := o.QueryTable("sys_user")
	for _,v :=range id_arr{
		id_int := utils.StrToInt(v)
		qs.Filter("id",id_int).Update(orm.Params{
			"is_delete":1,
		})


	}

	ret := fmt.Sprintf("批量删除成功，用户ids:%d",ids)
	logs.Info(ret)

	u.Data["json"] = map[string]interface{}{"code":200,"msg":"批量删除成功"}
	u.ServeJSON()


}