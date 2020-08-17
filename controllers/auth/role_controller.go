package auth

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ions_zhiliao/models/auth"
	"math"
	"ions_zhiliao/utils"
	"time"
	"strings"
)

type RoleController struct {
	beego.Controller
}

func (r *RoleController) List()  {

	roles := []auth.Role{}
	o := orm.NewOrm()


	// 每页显示的条数
	pagePerNum := 8
	// 当前页
	currentPage,err := r.GetInt("page")
	if err != nil {   // 说明没有获取到当前页
		currentPage = 1
	}

	offsetNum := pagePerNum * (currentPage - 1)

	qs := o.QueryTable("sys_role")
	qs.Filter("is_delete",0).All(&roles)

	count,_ := qs.Filter("is_delete",0).Count()
	qs.Filter("is_delete",0).Limit(pagePerNum).Offset(offsetNum).All(&roles)

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



	r.Data["roles"] = roles
	r.Data["prePage"] =prePage
	r.Data["nextPage"] = nextPage
	r.Data["currentPage"] = currentPage
	r.Data["countPage"] = countPage
	r.Data["count"] = count
	r.Data["page_map"] = page_map
	r.TplName = "auth/role_list.html"
}

func (r *RoleController) ToAdd()  {
	r.TplName = "auth/role_add.html"

}

func (r *RoleController) DoAdd()  {

	role_name := r.GetString("role_name")
	desc := r.GetString("desc")
	is_active,_ := r.GetInt("is_active")

	role := auth.Role{RoleName:role_name,Desc:desc,IsActive:is_active,CreateTime:time.Now()}
	o := orm.NewOrm()
	_,err := o.Insert(&role)

	message_map := map[string]interface{}{}
	if err != nil { // 发生错误
		message_map["code"] = 10001
		message_map["msg"] = "添加数据错误，请重新添加"

	}else {
		message_map["code"] = 200
		message_map["msg"] = "添加成功"
	}

	r.Data["json"] = message_map
	r.ServeJSON()

}

// 角色--一用户配置
func (r *RoleController) ToRoleUser()  {
	id,_ := r.GetInt("role_id")

	o := orm.NewOrm()
	role := auth.Role{}
	o.QueryTable("sys_role").Filter("id",id).One(&role)

	// 已绑定的用户
	o.LoadRelated(&role,"User")


	// 未绑定的用户
	users := []auth.User{}
	if len(role.User) > 0 {
		o.QueryTable("sys_user").Filter("is_delete",0).Filter("is_active",1).Exclude("id__in",role.User).All(&users)

	}else {   // 没有绑定的数据
		o.QueryTable("sys_user").Filter("is_delete",0).Filter("is_active",1).All(&users)

	}

	r.Data["role"] = role
	r.Data["users"] = users
	r.TplName = "auth/role-user-add.html"
	
}


// 角色--一用户配置
func (r *RoleController) DoRoleUser()  {
	role_id,_ := r.GetInt("role_id")
	user_ids := r.GetString("user_ids")

	//new_user_ids := user_ids[1:len(user_ids)-1]
	user_id_arr := strings.Split(user_ids,",")

	// "10,12,13"

	o := orm.NewOrm()
	role := auth.Role{Id:role_id}

	// 查询出已绑定的数据
	m2m := o.QueryM2M(&role,"User")
	m2m.Clear()

	for _,user_id := range user_id_arr{
		user := auth.User{Id:utils.StrToInt(user_id)}
		m2m := o.QueryM2M(&role,"User")
		m2m.Add(user)

	}

	r.Data["json"] = map[string]interface{}{"code":200,"msg":"添加成功"}
	r.ServeJSON()
}

// 角色--权限配置
func (r *RoleController) ToRoleAuth()  {
	role_id,_ := r.GetInt("role_id")

	o := orm.NewOrm()
	qs := o.QueryTable("sys_role")
	role := auth.Role{}
	qs.Filter("id",role_id).One(&role)
	r.Data["role"] = role
	r.TplName = "auth/role-auth-add.html"

}

func (r *RoleController) GetAuthJson()  {
	role_id,_ := r.GetInt("role_id")


	o := orm.NewOrm()
	qs := o.QueryTable("sys_auth")

	// 已绑定的权限
	role := auth.Role{Id:role_id}
	o.LoadRelated(&role,"Auth")

	//[11,14,16]
	auth_ids_has := []int{}
	for _,auth_data := range role.Auth{
		auth_ids_has = append(auth_ids_has,auth_data.Id)
	}




	// 所有的权限
	auths := []auth.Auth{}
	qs.Filter("is_delete",0).All(&auths)

	auth_arr_map := []map[string]interface{}{}  // map数组

	for _,auth_data := range auths{
		id := auth_data.Id
		pId := auth_data.Pid
		name := auth_data.AuthName
		if pId == 0 {
			auth_map := map[string]interface{}{"id":id,"pId":pId,"name":name,"open":false}
			auth_arr_map = append(auth_arr_map,auth_map)
		}else {
			auth_map := map[string]interface{}{"id":id,"pId":pId,"name":name}
			auth_arr_map = append(auth_arr_map,auth_map)
		}

	}

	auth_maps := map[string]interface{}{}
	auth_maps["auth_arr_map"] = auth_arr_map
	auth_maps["auth_ids_has"] = auth_ids_has
	r.Data["json"] = auth_maps
	r.ServeJSON()

}

func (r *RoleController) DoRoleAuth()  {

	role_id,_ := r.GetInt("role_id")
	auth_ids := r.GetString("auth_ids")
	//"13,15,16"       "13  15    16"
	//new_auth_ids := auth_ids[1:len(auth_ids)-1]
	id_arr := strings.Split(auth_ids,",")


	o := orm.NewOrm()
	role := auth.Role{Id:role_id}
	m2m := o.QueryM2M(&role,"Auth")
	m2m.Clear()

	for _,auth_id := range id_arr{
		auth_id_int := utils.StrToInt(auth_id)
		if auth_id_int !=0 {
			auth_data := auth.Auth{Id:auth_id_int}
			m2m := o.QueryM2M(&role,"Auth")
			m2m.Add(auth_data)
		}

	}

	r.Data["json"] = map[string]interface{}{"code":200,"msg":"添加成功"}
	r.ServeJSON()
	
	
}