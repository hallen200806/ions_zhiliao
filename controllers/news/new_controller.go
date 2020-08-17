package news

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ions_zhiliao/models/news"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math"
	"ions_zhiliao/utils"
	"time"
	"strconv"
)

type NewsController struct {
	beego.Controller
}

func (n *NewsController) Get()  {

	o := orm.NewOrm()

	qs := o.QueryTable("sys_news")

	news_data := []news.News{}
	// 每页显示的条数
	pagePerNum := 8
	// 当前页
	currentPage,err := n.GetInt("page")
	if err != nil {   // 说明没有获取到当前页
		currentPage = 1
	}

	offsetNum := pagePerNum * (currentPage - 1)


	kw := n.GetString("kw")

	var count int64 = 0

	ret := fmt.Sprintf("当前页;%d,查询条件：%s",currentPage,kw)
	logs.Info(ret)
	if kw != ""{   // 有查询条件的
		// 总数
		count,_ = qs.Filter("is_delete",0).Filter("title__contains",kw).Count()
		qs.Filter("is_delete",0).Filter("title__contains",kw).Limit(pagePerNum).Offset(offsetNum).RelatedSel().All(&news_data)
	}else {
		count,_ = qs.Filter("is_delete",0).Count()
		qs.Filter("is_delete",0).Limit(pagePerNum).Offset(offsetNum).RelatedSel().All(&news_data)

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

	n.Data["news_data"] = news_data
	n.Data["prePage"] =prePage
	n.Data["nextPage"] = nextPage
	n.Data["currentPage"] = currentPage
	n.Data["countPage"] = countPage
	n.Data["count"] = count
	n.Data["page_map"] = page_map
	n.Data["kw"] = kw


	n.TplName = "news/news_list.html"

}

func (n *NewsController) ToAdd()  {
	o := orm.NewOrm()
	qs := o.QueryTable("sys_category")
	categories := []news.Category{}
	qs.Filter("is_delete",0).All(&categories)
	n.Data["categories"] = categories
	n.TplName = "news/news_add.html"

}

func (n *NewsController) DoAdd()  {
	content := n.GetString("content")
	title := n.GetString("title")
	category_id,_ := n.GetInt("category_id")
	is_active,_ := n.GetInt("is_active")

	category := news.Category{Id:category_id}
	o := orm.NewOrm()
	news_data := news.News{
		Content:content,
		Title:title,
		Category:&category,
		IsActive:is_active,
	}
	_,err := o.Insert(&news_data)

	message_map := map[string]interface{}{}
	if err != nil {
		message_map["code"] = 10001
		message_map["msg"] = "添加失败"
	}
	message_map["code"] = 200
	message_map["msg"] = "添加成功"

	n.Data["json"] = message_map
	n.ServeJSON()
}

func (n *NewsController) UploadImg()  {

	f,h,err := n.GetFile("file")

	message_map := map[string]interface{}{}

	defer func() {
		f.Close()
	}()

	file_name := h.Filename

	time_unix_int := time.Now().Unix()
	time_unit_str := strconv.FormatInt(time_unix_int,10)

	file_path := "upload/news_img/"+ time_unit_str + "-" + file_name

	img_link := "http://127.0.0.1:8080/" + file_path

	if err != nil {
		message_map["code"] = 10001
		message_map["msg"] = "文件上传失败"
		message_map["link"] = img_link

	}

	n.SaveToFile("file",file_path)

	message_map["code"] = 200
	message_map["msg"] = "文件上传成功"
	message_map["link"] = img_link

	n.Data["json"] = message_map
	n.ServeJSON()


}

func (n *NewsController) ToEdit()  {
	n_id,_ := n.GetInt("id")

	o := orm.NewOrm()
	qs := o.QueryTable("sys_news")

	news_data := news.News{}
	qs.Filter("id",n_id).RelatedSel().One(&news_data)
	n.Data["news_data"] = news_data

	categories := []news.Category{}
	o.QueryTable("sys_category").Exclude("id",news_data.Category.Id).All(&categories)

	n.Data["news_data"] = news_data
	n.Data["categories"] = categories
	n.TplName = "news/news_edit.html"

}

func (n *NewsController) DoEdit()  {

	news_id,_ := n.GetInt("news_id")
	content := n.GetString("content")
	title := n.GetString("title")
	category_id,_ := n.GetInt("category_id")
	is_active,_ := n.GetInt("is_active")

	o := orm.NewOrm()
	qs := o.QueryTable("sys_news")
	_,err := qs.Filter("id",news_id).Update(orm.Params{
		"title":title,
		"content":content,
		"category_id":category_id,
		"is_active":is_active,
	})
	message_map := map[string]interface{}{}
	if err != nil {
		message_map["code"] = 10001
		message_map["msg"] = "更新失败"
	}
	message_map["code"] = 200
	message_map["msg"] = "更新成功"

	n.Data["json"] = message_map
	n.ServeJSON()







}
