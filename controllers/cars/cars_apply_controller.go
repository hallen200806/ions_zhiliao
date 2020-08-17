package cars

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math"
	"ions_zhiliao/models/auth"
	"ions_zhiliao/utils"
	"time"
)

type CarsApplyController struct {
	beego.Controller
}

func (c *CarsApplyController) Get()  {
	o := orm.NewOrm()

	qs := o.QueryTable("sys_cars")

	cars_data := []auth.Cars{}

	// 每页显示的条数
	pagePerNum := 8
	// 当前页
	currentPage,err := c.GetInt("page")
	if err != nil {   // 说明没有获取到当前页
		currentPage = 1
	}

	offsetNum := pagePerNum * (currentPage - 1)


	kw := c.GetString("kw")
	var count int64 = 0

	ret := fmt.Sprintf("当前页;%d,查询条件：%s",currentPage,kw)
	logs.Info(ret)
	if kw != ""{   // 有查询条件的
		// 总数
		count,_ = qs.Filter("is_delete",0).Filter("name__contains",kw).Count()
		qs.Filter("is_delete",0).Filter("name__contains",kw).Limit(pagePerNum).Offset(offsetNum).RelatedSel().All(&cars_data)
	}else {
		count,_ = qs.Filter("is_delete",0).Count()
		qs.Filter("is_delete",0).Limit(pagePerNum).Offset(offsetNum).RelatedSel().All(&cars_data)

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
	c.Data["cars_data"] = cars_data
	c.Data["prePage"] =prePage
	c.Data["nextPage"] = nextPage
	c.Data["currentPage"] = currentPage
	c.Data["countPage"] = countPage
	c.Data["count"] = count
	c.Data["page_map"] = page_map
	c.Data["kw"] = kw


	c.TplName = "cars/cars_apply_list.html"

}

func (c *CarsApplyController) ToApply()  {
	id,_ := c.GetInt("id")
	c.Data["id"] = id
	c.TplName = "cars/cars_apply.html"

}

func (c *CarsApplyController) DoApply()  {
	reason := c.GetString("reason")
	destination := c.GetString("destination")
	return_date := c.GetString("return_date")
	return_date_new,_ := time.Parse("2006-01-02",return_date)
	cars_id,_ := c.GetInt("cars_id")
	uid := c.GetSession("id")

	// interface --> int
	user := auth.User{Id:uid.(int)}
	cars_date := auth.Cars{Id:cars_id}

	o := orm.NewOrm()

	// 默认：ReturnStatus=0，AuditStatus=3，IsActive=1
	cars_apply := auth.CarsApply{
		User:&user,
		Cars:&cars_date,
		Reason:reason,
		Destination:destination,
		ReturnDate:return_date_new,
		ReturnStatus:0,
		AuditStatus:3,
		IsActive:1,


	}
	_,err := o.Insert(&cars_apply)


	o.QueryTable("sys_cars").Filter("id",cars_id).Update(orm.Params{
		"status":2,
	})

	message_map := map[string]interface{}{}
	if err != nil {
		message_map["code"] = 10001
		message_map["msg"] = "申请失败"
	}

	message_map["code"] = 200
	message_map["msg"] = "申请成功"

	c.Data["json"] = message_map
	c.ServeJSON()

}

func (c *CarsApplyController) MyApply()  {
	o := orm.NewOrm()

	qs := o.QueryTable("sys_cars_apply")

	cars_apply := []auth.CarsApply{}

	// 每页显示的条数
	pagePerNum := 8
	// 当前页
	currentPage,err := c.GetInt("page")
	if err != nil {   // 说明没有获取到当前页
		currentPage = 1
	}

	offsetNum := pagePerNum * (currentPage - 1)


	kw := c.GetString("kw")
	var count int64 = 0

	uid := c.GetSession("id")

	ret := fmt.Sprintf("当前页;%d,查询条件：%s",currentPage,kw)
	logs.Info(ret)
	if kw != ""{   // 有查询条件的
		// 总数
		count,_ = qs.Filter("is_delete",0).Filter("Cars__name__contains",kw).Filter("user_id",uid.(int)).Count()
		qs.Filter("is_delete",0).Filter("Cars__name__contains",kw).Limit(pagePerNum).Offset(offsetNum).RelatedSel().Filter("user_id",uid.(int)).All(&cars_apply)
	}else {
		count,_ = qs.Filter("is_delete",0).Filter("user_id",uid.(int)).Count()
		qs.Filter("is_delete",0).Limit(pagePerNum).Filter("user_id",uid.(int)).Offset(offsetNum).RelatedSel().All(&cars_apply)

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
	c.Data["cars_apply"] = cars_apply
	c.Data["prePage"] =prePage
	c.Data["nextPage"] = nextPage
	c.Data["currentPage"] = currentPage
	c.Data["countPage"] = countPage
	c.Data["count"] = count
	c.Data["page_map"] = page_map
	c.Data["kw"] = kw
	c.TplName = "cars/my_apply_list.html"
}

func (c *CarsApplyController) AuditApply()  {
	o := orm.NewOrm()

	qs := o.QueryTable("sys_cars_apply")

	cars_apply := []auth.CarsApply{}

	// 每页显示的条数
	pagePerNum := 8
	// 当前页
	currentPage,err := c.GetInt("page")
	if err != nil {   // 说明没有获取到当前页
		currentPage = 1
	}

	offsetNum := pagePerNum * (currentPage - 1)


	kw := c.GetString("kw")
	var count int64 = 0

	ret := fmt.Sprintf("当前页;%d,查询条件：%s",currentPage,kw)
	logs.Info(ret)
	if kw != ""{   // 有查询条件的
		// 总数
		count,_ = qs.Filter("is_delete",0).Filter("Cars__name__contains",kw).Count()
		qs.Filter("is_delete",0).Filter("Cars__name__contains",kw).Limit(pagePerNum).Offset(offsetNum).RelatedSel().All(&cars_apply)
	}else {
		count,_ = qs.Filter("is_delete",0).Count()
		qs.Filter("is_delete",0).Limit(pagePerNum).Offset(offsetNum).RelatedSel().All(&cars_apply)

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
	c.Data["cars_apply"] = cars_apply
	c.Data["prePage"] =prePage
	c.Data["nextPage"] = nextPage
	c.Data["currentPage"] = currentPage
	c.Data["countPage"] = countPage
	c.Data["count"] = count
	c.Data["page_map"] = page_map
	c.Data["kw"] = kw
	c.TplName = "cars/audit_apply_list.html"

}

func (c *CarsApplyController) ToAuditApply()  {
	id,_ := c.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_cars_apply")
	cars_apply := auth.CarsApply{}
	qs.Filter("id",id).One(&cars_apply)
	c.Data["cars_apply"] = cars_apply
	c.TplName = "cars/audit_apply.html"
	
}

func (c *CarsApplyController) DoAuditApply()  {

	option := c.GetString("option")
	audit_status,_ := c.GetInt("audit_status")

	id,_ := c.GetInt("id")

	o := orm.NewOrm()

	qs := o.QueryTable("sys_cars_apply")
	_,err := qs.Filter("id",id).Update(orm.Params{
		"audit_option":option,
		"audit_status":audit_status,
	})

	message_map := map[string]interface{}{}
	if err != nil {
		message_map["code"] = 10001
		message_map["msg"] = "审核失败"
	}

	message_map["code"] = 200
	message_map["msg"] = "审核成功"

	c.Data["json"] = message_map
	c.ServeJSON()
	

}

func (c *CarsApplyController) DoReturn()  {

	id,_ := c.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_cars_apply")
	qs.Filter("id",id).Update(orm.Params{
		"return_status":1,
	})

	cars_apply := auth.CarsApply{}
	qs.Filter("id",id).One(&cars_apply)

	o.QueryTable("sys_cars").Filter("id",cars_apply.Cars.Id).Update(orm.Params{
		"status":1,
	})


	c.Redirect(beego.URLFor("CarsApplyController.MyApply"),302)

	
}