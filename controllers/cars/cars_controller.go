package cars

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math"
	"ions_zhiliao/models/auth"
	"ions_zhiliao/utils"
)

type CarsController struct {
	beego.Controller
}

func (c *CarsController) Get()  {
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
	c.TplName = "cars/cars_list.html"

}

func (c *CarsController) ToAdd()  {
	o := orm.NewOrm()
	qs := o.QueryTable("sys_cars_brand")
	cars_brand := []auth.CarBrand{}
	qs.Filter("is_delete",0).All(&cars_brand)
	c.Data["cars_brand"] = cars_brand
 	c.TplName = "cars/cars_add.html"

}

func (c *CarsController) DoAdd()  {
	cars_brand_id,_ := c.GetInt("cars_brand_id")
	name := c.GetString("name")
	is_active,_ := c.GetInt("is_active")
	o := orm.NewOrm()

	cars_brand := auth.CarBrand{Id:cars_brand_id}
	cars_data := auth.Cars{
		Name:name,
		CarBrand:&cars_brand,
		IsActive:is_active,

	}
	_,err := o.Insert(&cars_data)

	message_map := map[string]interface{}{}
	if err != nil {
		message_map["code"] = 10001
		message_map["msg"] = "添加失败"
	}

	message_map["code"] = 200
	message_map["msg"] = "添加成功"
	c.Data["json"] = message_map
	c.ServeJSON()


}



