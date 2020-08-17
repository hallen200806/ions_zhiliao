package caiwu

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"ions_zhiliao/models/my_center"
	"math"
	"ions_zhiliao/utils"
	"strconv"
	"github.com/Luxurioust/excelize"
	"github.com/astaxie/beego/logs"
)

type CaiWuSalarySlipController struct {
	beego.Controller
}

func (c *CaiWuSalarySlipController) Get()  {

	o := orm.NewOrm()
	qs := o.QueryTable("sys_salary_slip")
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
	salary_slips := []my_center.SalarySlip{}
	if month != ""{   // 有查询条件的
		// 总数
		count,_ = qs.Filter("pay_date",month).Count()
		qs.Filter("pay_date",month).Limit(pagePerNum).Offset(offsetNum).All(&salary_slips)
	}else {
		month = time.Now().Format("2006-01")
		count,_ = qs.Filter("pay_date",month).Count()
		qs.Filter("pay_date",month).Limit(pagePerNum).Offset(offsetNum).All(&salary_slips)

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

	c.Data["salary_slips"] = salary_slips
	c.Data["prePage"] =prePage
	c.Data["nextPage"] = nextPage
	c.Data["currentPage"] = currentPage
	c.Data["countPage"] = countPage
	c.Data["count"] = count
	c.Data["page_map"] = page_map
	c.Data["month"] = month
	c.TplName = "caiwu/salary_slip_list.html"

}

func (c *CaiWuSalarySlipController) ToImportExcel()  {
	c.TplName = "caiwu/salary_slip_import.html"
	
}

func (c *CaiWuSalarySlipController) DoImportExcel()  {
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

	file_path := "upload/salary_slip_upload/"+ time_unit_str + "-" + file_name

	c.SaveToFile("upload_file",file_path)


	// 读取数据并插入数据库
	file,err1 := excelize.OpenFile(file_path)
	logs.Error(err1)
	rows,_ := file.GetRows("Sheet1")

	o := orm.NewOrm()

	new_salary_data := []my_center.SalarySlip{}
	i := 0
	for _,row := range rows {



		card_id := row[2]
		base_pay,_ := strconv.ParseFloat(row[3],64)
		working_days,_ := strconv.ParseFloat(row[4],64)
		days_off,_ := strconv.ParseFloat(row[5],64)
		days_off_no,_ := strconv.ParseFloat(row[6],64)
		reward,_ := strconv.ParseFloat(row[7],64)
		rent_subsidy,_ := strconv.ParseFloat(row[8],64)
		trans_subsidy,_ := strconv.ParseFloat(row[9],64)
		social_security,_ := strconv.ParseFloat(row[10],64)
		hous_provident_fund,_ := strconv.ParseFloat(row[11],64)
		personal_pncome_taxm,_ := strconv.ParseFloat(row[12],64)
		fine,_ := strconv.ParseFloat(row[13],64)
		net_salary,_ := strconv.ParseFloat(row[14],64)
		pay_date := row[15]

		salary_slip := my_center.SalarySlip{
			CardId:card_id,
			BasePay:base_pay,
			WorkingDays:working_days,
			DaysOff:days_off,
			DaysOffNo:days_off_no,
			Reward:reward,
			RentSubsidy:rent_subsidy,
			TransSubsidy:trans_subsidy,
			SocialSecurity:social_security,
			HousProvidentFund:hous_provident_fund,
			PersonalPncomeTax:personal_pncome_taxm,
			Fine:fine,
			NetSalary:net_salary,
			PayDate:pay_date,


		}

		if i == 0 {

			i ++
			continue
		}

		// 重复导入相同月份的数据：先删除已有的工资月份，再导入
		qs := o.QueryTable("sys_salary_slip")
		is_exist := qs.Filter("pay_date",pay_date).Exist()
		if is_exist{
			qs.Filter("pay_date",pay_date).Delete()
		}

		// 精确到导入失败的数据信息提示
		_,err := o.Insert(&salary_slip)

		if err != nil {  // 报错的数据
			err_data_arr = append(err_data_arr,card_id)
		}
		new_salary_data = append(new_salary_data,salary_slip)

		i ++


	}

	if len(err_data_arr) <= 0 {
		message_map["code"] = 200
		message_map["msg"] = "导入成功"
	} else {
		o.InsertMulti(100,&new_salary_data)
		message_map["code"] = 10002
		message_map["msg"] = "导入失败"
		message_map["err_data"] = err_data_arr

	}


	c.Data["json"] = message_map
	c.ServeJSON()


}
