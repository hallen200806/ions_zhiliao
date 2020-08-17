package news

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math"
	"ions_zhiliao/utils"
	"ions_zhiliao/models/news"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get()  {

	o := orm.NewOrm()

	qs := o.QueryTable("sys_category")

	categrories := []news.Category{}
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
		qs.Filter("is_delete",0).Filter("name__contains",kw).Limit(pagePerNum).Offset(offsetNum).All(&categrories)
	}else {
		count,_ = qs.Filter("is_delete",0).Count()
		qs.Filter("is_delete",0).Limit(pagePerNum).Offset(offsetNum).All(&categrories)

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

	c.Data["categrories"] = categrories
	c.Data["prePage"] =prePage
	c.Data["nextPage"] = nextPage
	c.Data["currentPage"] = currentPage
	c.Data["countPage"] = countPage
	c.Data["count"] = count
	c.Data["page_map"] = page_map
	c.Data["kw"] = kw

	c.TplName = "news/category_list.html"

}

func (c *CategoryController) ToAdd()  {
	c.TplName = "news/category_add.html"

}

func (c *CategoryController) DoAdd()  {

	name := c.GetString("name")
	desc := c.GetString("desc")
	is_active,_ := c.GetInt("is_active")


	o := orm.NewOrm()
	category := news.Category{
		Name:name,
		Desc:desc,
		IsActive:is_active,
	}
	_,err := o.Insert(&category)


	message_map := map[string]interface{}{}
	if err != nil{
		message_map["code"] = 10001
		message_map["msg"] = "添加栏目失败"

	}

	message_map["code"] = 200
	message_map["msg"] = "添加成功"

	c.Data["json"] = message_map
	c.ServeJSON()
	
}
