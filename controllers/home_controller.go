package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ions_zhiliao/models/auth"
	"time"
	"github.com/astaxie/beego/logs"
	"math"
	"ions_zhiliao/utils"
	"fmt"
)

type HomeController struct {
	beego.Controller
}

func (h *HomeController) Get()  {
	// 后端首页
	o := orm.NewOrm()

	user_id := h.GetSession("id")
	// interface --> int
	user := auth.User{Id:user_id.(int)}

	o.LoadRelated(&user,"Role")

	auth_arr := []int{}
	for _,role := range user.Role{
		role_data := auth.Role{Id:role.Id}
		o.LoadRelated(&role_data,"Auth")
		for _,auth_date := range role_data.Auth{
			auth_arr = append(auth_arr,auth_date.Id)
		}

	}




	qs := o.QueryTable("sys_auth")

	auths := []auth.Auth{}
	qs.Filter("pid",0).Filter("id__in",auth_arr).OrderBy("-weight").All(&auths)
	//"select * from sys_user where id in (1,2,3,1)"

	trees := []auth.Tree{}
	for _,auth_data := range auths{   // 一级菜单

		pid := auth_data.Id   // 根据pid获取所有的子解点
		tree_data := auth.Tree{Id:auth_data.Id,AuthName:auth_data.AuthName,UrlFor:auth_data.UrlFor,Weight:auth_data.Weight,Children:[]*auth.Tree{}}
		GetChildNode(pid,&tree_data)
		trees = append(trees,tree_data)

	}

	//for _,tree_data := range trees{
	//	for _,tree_data2 := range tree_data.Children{
	//		fmt.Println(tree_data2)
	//	}
	//}

	o.QueryTable("sys_user").Filter("id",user_id).One(&user)



	// 消息通知,发送消息，使用定时任务优化
	qs1 := o.QueryTable("sys_cars_apply")
	cars_apply := []auth.CarsApply{}
	qs1.Filter("user_id",user_id.(int)).Filter("return_status",0).Filter("notify_tag",0).All(&cars_apply)

	cur_time,_ := time.Parse("2006-01-02",time.Now().Format("2006-01-02"))

	for _,apply := range cars_apply{
		return_date := apply.ReturnDate
		ret := cur_time.Sub(return_date)
		content := fmt.Sprintf("%s用户，你借的车辆归还时间为%v,已经预期，请尽快归还!!",user.UserName,return_date.Format("2006-01-02"))
		if ret > 0 {  // 已经逾期
			message_notify := auth.MessageNotify{
				Flag:1,
				Title:"车辆归还逾期",
				Content:content,
				User:&user,
				ReadTag:0,

			}
			o.Insert(&message_notify)
		}

		apply.NotifyTag = 1

		o.Update(&apply)

	}

	// 展示消息,使用websocket优化

	qs2 := o.QueryTable("sys_message_notify")
	notify_count,_ := qs2.Filter("read_tag",0).Count()


	h.Data["notify_count"] = notify_count
	h.Data["trees"] = trees
	h.Data["user"] = user
	h.TplName = "index.html"

}

func (h *HomeController) Welcome()  {
	h.TplName = "welcome.html"
}


// 递归
func GetChildNode(pid int, treenode *auth.Tree)  {

	o := orm.NewOrm()

	qs := o.QueryTable("sys_auth")
	auths := []auth.Auth{}
	_,err := qs.Filter("pid",pid).OrderBy("-weight").All(&auths)

	if err != nil {
		return
	}


	// 查询三级及以上的菜单
	for i:= 0; i<len(auths);i++{
		pid := auths[i].Id   // 根据pid获取所有的子解点
		tree_data := auth.Tree{Id:auths[i].Id,AuthName:auths[i].AuthName,UrlFor:auths[i].UrlFor,Weight:auths[i].Weight,Children:[]*auth.Tree{}}
		treenode.Children = append(treenode.Children,&tree_data)
		GetChildNode(pid,&tree_data)
	}

	return


}

func (h *HomeController) NotifyList()  {
	o := orm.NewOrm()

	qs := o.QueryTable("sys_message_notify")

	nofities := []auth.MessageNotify{}
	// 每页显示的条数
	pagePerNum := 8
	// 当前页
	currentPage,err := h.GetInt("page")

	offsetNum := pagePerNum * (currentPage - 1)


	kw := h.GetString("kw")
	var count int64 = 0

	ret := fmt.Sprintf("当前页;%d,查询条件：%s",currentPage,kw)
	logs.Info(ret)
	if kw != ""{   // 有查询条件的
		// 总数
		count,_ = qs.Filter("title__contains",kw).Count()
		qs.Filter("title__contains",kw).Limit(pagePerNum).Offset(offsetNum).All(&nofities)
	}else {
		count,_ = qs.Count()
		qs.Limit(pagePerNum).Offset(offsetNum).All(&nofities)

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

	h.Data["nofities"] = nofities
	h.Data["prePage"] =prePage
	h.Data["nextPage"] = nextPage
	h.Data["currentPage"] = currentPage
	h.Data["countPage"] = countPage
	h.Data["count"] = count
	h.Data["page_map"] = page_map
	h.Data["kw"] = kw

	h.TplName = "notify_list.html"

}

func (h *HomeController) ReadNotify()  {
	id,_ := h.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_message_notify")
	qs.Filter("id",id).Update(orm.Params{
		"read_tag":1,
	})
	h.Redirect(beego.URLFor("HomeController.NotifyList"),302)


}
