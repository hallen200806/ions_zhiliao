package echarts

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type EchartsCourseController struct {
	beego.Controller
}

func (e *EchartsCourseController) Get()  {
	e.TplName = "echarts/echarts_course.html"

}

func (e *EchartsCourseController) GetCourseChart()  {

	var caiwu_date orm.ParamsList
	o := orm.NewOrm()
	o.Raw("select caiwu_date from sys_caiwu_data").ValuesFlat(&caiwu_date)

	var course_type orm.ParamsList = orm.ParamsList{"django","vue_django","celery"}

	var django_list orm.ParamsList
	o.Raw("select django from sys_caiwu_data").ValuesFlat(&django_list)

	var vue_django_list orm.ParamsList
	o.Raw("select vue_django from sys_caiwu_data").ValuesFlat(&vue_django_list)

	var celery_list orm.ParamsList
	o.Raw("select celery from sys_caiwu_data").ValuesFlat(&celery_list)


	map_data := map[string]interface{}{}

	//map_series := map[string]interface{}{}

	series_data := []map[string]interface{}{}


	course_mapping := map[string]interface{}{
		"django":django_list,
		"vue_django":vue_django_list,
		"celery":celery_list,
	}
	for i:=0;i<len(course_type) ;i++  {
		map_series := map[string]interface{}{}
		map_series["name"] =  course_type[i]
		map_series["type"] =  "line"
		map_series["stack"] =  "总量"
		data_list := course_mapping[course_type[i].(string)]
		map_series["data"] = data_list

		series_data = append(series_data,map_series)
	}

	map_data["caiwu_date"]  = caiwu_date
	map_data["course_type"] = course_type
	map_data["series_data"] = series_data

	e.Data["json"] = map_data
	e.ServeJSON()

}
