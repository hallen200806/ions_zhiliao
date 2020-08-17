package auth

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"ions_zhiliao/utils"
	"ions_zhiliao/models/auth"
	"time"
)

type AuthController struct {
	beego.Controller
}

func (a *AuthController) List()  {

	o := orm.NewOrm()
	qs := o.QueryTable("sys_auth")

	auths := []auth.Auth{}

	// 每页显示的条数
	pagePerNum := 8

	// 当前页
	currentPage,err := a.GetInt("page")
	if err != nil {
		currentPage = 1
	}
	// 总数
	count,_ := qs.Filter("is_delete",0).Count()

	// 总页数
	countPage := int(math.Ceil(float64(count) / float64(pagePerNum)))

	qs.Filter("is_delete",0).All(&auths)

	page_map := utils.Paginator(currentPage,pagePerNum,count)

	// 当前页码小于总页数，下一页可以+1
	nextPage := 1
	if currentPage < countPage{
		nextPage = currentPage + 1
	}else if currentPage >= countPage{   // 不能+1
		nextPage = currentPage
	}

	// 当前页面=1，不能-1

	prePage := 1
	if currentPage == 1{
		prePage = currentPage
	}else if currentPage >1{
		prePage = currentPage - 1
	}


	a.Data["page_map"] = page_map
	a.Data["countPage"] = countPage
	a.Data["count"] = count
	a.Data["auths"] = auths
	a.Data["currentPage"] = currentPage
	a.Data["prePage"] = prePage
	a.Data["nextPage"] = nextPage
	a.TplName = "auth/auth-list.html"

}

func (a *AuthController) ToAdd()  {

	auths := []auth.Auth{}

	o := orm.NewOrm()
	qs := o.QueryTable("sys_auth")
	qs.Filter("is_delete",0).All(&auths)
	a.Data["auths"] = auths
	a.TplName = "auth/auth-add.html"

}

func (a *AuthController) DoAdd()  {
	auth_parent_id,_ := a.GetInt("auth_parent_id")
	auth_name := a.GetString("auth_name")
	auth_url := a.GetString("auth_url")
	auth_desc := a.GetString("auth_desc")
	is_active,_ := a.GetInt("is_active")
	auth_weight,_ := a.GetInt("auth_weight")


	o := orm.NewOrm()

	auth_data := auth.Auth{AuthName:auth_name,UrlFor:auth_url,Pid:auth_parent_id,Desc:auth_desc,CreateTime:time.Now(),IsActive:is_active,Weight:auth_weight}
	o.Insert(&auth_data)


	a.Data["json"] = map[string]interface{}{"code":200,"msg":"添加成功"}
	a.ServeJSON()

}

