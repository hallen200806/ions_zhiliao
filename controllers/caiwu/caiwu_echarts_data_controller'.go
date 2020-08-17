package caiwu

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"math"
	"ions_zhiliao/models/caiwu"
	"ions_zhiliao/utils"
	"strconv"
	"github.com/Luxurioust/excelize"
	"github.com/astaxie/beego/logs"
)

type CaiWuEchartDataController struct {
	beego.Controller
}

func (c *CaiWuEchartDataController) Get()  {


	o := orm.NewOrm()
	qs := o.QueryTable("sys_caiwu_data")
	// 每页显示的条数
	pagePerNum := 8
	// 当前页
	currentPage,err := c.GetInt("page")
	if err != nil {   // 说明没有获取到当前页
		currentPage = 1
	}

	offsetNum := pagePerNum * (currentPage - 1)


	month := c.GetString("month")
	var count int64 = 0
	caiwu_datas := []caiwu.CaiwuData{}
	if month != ""{   // 有查询条件的
		// 总数
		count,_ = qs.Filter("caiwu_date",month).Count()
		qs.Filter("caiwu_date",month).Limit(pagePerNum).Offset(offsetNum).All(&caiwu_datas)
	}else {
		month = time.Now().Format("2006-01")

		count,_ = qs.Filter("caiwu_date",month).Count()

		qs.Filter("caiwu_date",month).Limit(pagePerNum).Offset(offsetNum).All(&caiwu_datas)

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

	c.Data["caiwu_datas"] = caiwu_datas
	c.Data["prePage"] =prePage
	c.Data["nextPage"] = nextPage
	c.Data["currentPage"] = currentPage
	c.Data["countPage"] = countPage
	c.Data["count"] = count
	c.Data["page_map"] = page_map
	c.Data["month"] = month

	c.TplName = "caiwu/echart_data_list.html"

}

func (c *CaiWuEchartDataController) ToImportExcel()  {
	c.TplName = "caiwu/echart_data_import.html"

}


func (c *CaiWuEchartDataController) DoImportExcel()  {
	f,h,err := c.GetFile("upload_file")

	message_map := map[string]interface{}{}
	err_data_arr := []string{}

	defer func() {
		f.Close()
	}()
	if err != nil {
		message_map["code"] = 10001
		message_map["msg"] = "文件上传失败"
		c.Data["json"] = message_map
		c.ServeJSON()
	}

	file_name := h.Filename

	time_unix_int := time.Now().Unix()
	time_unit_str := strconv.FormatInt(time_unix_int,10)

	file_path := "upload/echart_data_upload/"+ time_unit_str + "-" + file_name

	c.SaveToFile("upload_file",file_path)


	// 读取数据并插入数据库
	file,err1 := excelize.OpenFile(file_path)
	logs.Error(err1)
	rows,_ := file.GetRows("Sheet1")

	o := orm.NewOrm()


	i := 0
	for _,row := range rows {
		caiwu_date := row[0]
		sales_volume,_ := strconv.ParseFloat(row[1],64)
		student_incress := utils.StrToInt(row[2])
		django := utils.StrToInt(row[3])
		vue_django := utils.StrToInt(row[4])
		celery := utils.StrToInt(row[5])

		echart_data := caiwu.CaiwuData{
			CaiWuDate:caiwu_date,
			SalesVolume:sales_volume,
			StudentIncress:student_incress,
			Django:django,
			VueDjango:vue_django,
			Celery:celery,
		}

		if i == 0 {
			i ++
			continue
		}


		// 重复导入相同月份的数据：先删除已有的工资月份，再导入
		qs := o.QueryTable("sys_caiwu_data")
		is_exist := qs.Filter("caiwu_date",caiwu_date).Exist()

		if is_exist{
			qs.Filter("caiwu_date",caiwu_date).Delete()
		}



		// 精确到导入失败的数据信息提示
		_,err := o.Insert(&echart_data)

		if err != nil {  // 报错的数据
			err_data_arr = append(err_data_arr,caiwu_date)
		}
		i ++

	}

	if len(err_data_arr) <= 0 {
		message_map["code"] = 200
		message_map["msg"] = "导入成功"
	} else {
		message_map["code"] = 10002
		message_map["msg"] = "导入失败"
		message_map["err_data"] = err_data_arr

	}

	c.Data["json"] = message_map
	c.ServeJSON()

}
