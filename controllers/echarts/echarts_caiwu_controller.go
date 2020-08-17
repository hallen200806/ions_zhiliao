package echarts

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type EchartsCaiwuController struct {
	beego.Controller
}

func (e *EchartsCaiwuController) Get()  {
	e.TplName = "echarts/echarts_caiwu.html"
}

func (e *EchartsCaiwuController) GetCaiWuChart()  {
	var caiwu_date orm.ParamsList
	var sales_volume orm.ParamsList
	o := orm.NewOrm()
	o.Raw("select caiwu_date from sys_caiwu_data").ValuesFlat(&caiwu_date)
	o.Raw("select sales_volume from sys_caiwu_data").ValuesFlat(&sales_volume)

	map_data := map[string]interface{}{}

	map_data["caiwu_date"]  = caiwu_date
	map_data["sales_volume"] = sales_volume

	e.Data["json"] = map_data
	e.ServeJSON()



}
