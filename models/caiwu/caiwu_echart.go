package caiwu

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type CaiwuData struct {
	Id int `orm:"pk;auto"`
	CaiWuDate string `orm:"description(财务月份);size(32);column(caiwu_date)"`
	SalesVolume float64 `orm:"digits(10);decimals(2);description(本月销售额)"`
	StudentIncress int `orm:"description(学员增加数)"`
	Django int `orm:"description(django课程卖出数量)"`
	VueDjango int `orm:"description(vue+django课程卖出数量)"`
	Celery int `orm:"description(celery课程卖出数量)"`
	CreateDate time.Time `orm:"type(datetime);auto_now"`
}

func (c *CaiwuData) TableName() string {
	return "sys_caiwu_data"

}

func init()  {
	orm.RegisterModel(new(CaiwuData))

}

