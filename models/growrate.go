package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"github.com/ximply/myslowreport/utils"
	"github.com/shopspring/decimal"
)

type GrowRate struct {
	Id int64

	MyInsName string
	StatDate time.Time

	YesterdayTotal int64
	BeforeYesterdayTotal int64
	BasisTotal int64

	YesterdayUniq int64
	BeforeYesterdayUniq int64
	BasisUniq int64

	TotalChainRate decimal.Decimal
	TotalBasisRate decimal.Decimal

	UniqChainRate decimal.Decimal
	UniqBasisRate decimal.Decimal
}

const (
	Id = iota
	MyInsName
	StatDate

	YesterdayTotal
	BeforeYesterdayTotal
	BasisTotal

	YesterdayUniq
	BeforeYesterdayUniq
	BasisUniq

	TotalChainRate
	TotalBasisRate

	UniqChainRate
	UniqBasisRate
)

var tableGrowRateFields = []string{
	"Id",

	"MyInsName",
	"StatDate",

	"YesterdayTotal",
	"BeforeYesterdayTotal",
	"BasisTotal",

	"YesterdayUniq",
	"BeforeYesterdayUniq",
	"BasisUniq",

	"TotalChainRate",
	"TotalBasisRate",

	"UniqChainRate",
	"UniqBasisRate",

}

func tableName() string {
	return "myslow_history_grow_rate"
}

func GetByMyInsNameAndStatDate(name string, date string) (int64, GrowRate) {
	retry := DbRetry()
	sql := fmt.Sprintf("SELECT " +
		tableGrowRateFields[Id] + ", " +
		tableGrowRateFields[MyInsName] + ", " +
		tableGrowRateFields[StatDate] + ", " +
		tableGrowRateFields[YesterdayTotal] + ", " +
		tableGrowRateFields[BeforeYesterdayTotal] + ", " +
		tableGrowRateFields[BasisTotal] + ", " +
		tableGrowRateFields[YesterdayUniq] + ", " +
		tableGrowRateFields[BeforeYesterdayUniq] + ", " +
		tableGrowRateFields[BasisUniq] + ", " +
		tableGrowRateFields[TotalChainRate] + ", " +
		tableGrowRateFields[TotalBasisRate] + ", " +
		tableGrowRateFields[UniqChainRate] + ", " +
		tableGrowRateFields[UniqBasisRate] + " " +
		" FROM %s WHERE %s = '%s' AND %s = '%s'",

		tableName(), tableGrowRateFields[MyInsName], name, tableGrowRateFields[StatDate], date)

	var count int64
	var err error
	var sl []orm.ParamsList
	var gr GrowRate

	for i := 0; i < retry; i++ {
		count = 0
		err = nil
		count, err = orm.NewOrm().Raw(sql).ValuesList(&sl,
			tableGrowRateFields[Id],
			tableGrowRateFields[MyInsName],
			tableGrowRateFields[StatDate],
			tableGrowRateFields[YesterdayTotal],
			tableGrowRateFields[BeforeYesterdayTotal],
			tableGrowRateFields[BasisTotal],
			tableGrowRateFields[YesterdayUniq],
			tableGrowRateFields[BeforeYesterdayUniq],
			tableGrowRateFields[BasisUniq],
			tableGrowRateFields[TotalChainRate],
			tableGrowRateFields[TotalBasisRate],
			tableGrowRateFields[UniqChainRate],
			tableGrowRateFields[UniqBasisRate])
		if err == nil {
			break
		} else {
			beego.Info(fmt.Sprintf("Get error when query %s: %s", tableName(), err))
			continue
		}
	}

	if count > 0 {
		gr.Id = utils.InterfaceStringToInt64(sl[0][Id], defaultIntValue)
		gr.MyInsName = utils.InterfaceStringToString(sl[0][MyInsName], defaultStringValue)
		gr.StatDate = utils.InterfaceStringToTimeByFormat(sl[0][StatDate], "2006-01-02", defaultTime)

		gr.YesterdayTotal = utils.InterfaceStringToInt64(sl[0][YesterdayTotal], defaultIntValue)
		gr.BeforeYesterdayTotal = utils.InterfaceStringToInt64(sl[0][BeforeYesterdayTotal], defaultIntValue)
		gr.BasisTotal = utils.InterfaceStringToInt64(sl[0][BasisTotal], defaultIntValue)

		gr.YesterdayUniq = utils.InterfaceStringToInt64(sl[0][YesterdayUniq], defaultIntValue)
		gr.BeforeYesterdayUniq = utils.InterfaceStringToInt64(sl[0][BeforeYesterdayUniq], defaultIntValue)
		gr.BasisUniq = utils.InterfaceStringToInt64(sl[0][BasisUniq], defaultIntValue)

		gr.TotalChainRate = utils.InterfaceStringToDecimal(sl[0][TotalChainRate].(string), defaultDecimalValue)
		gr.TotalBasisRate = utils.InterfaceStringToDecimal(sl[0][TotalBasisRate], defaultDecimalValue)
		gr.UniqChainRate = utils.InterfaceStringToDecimal(sl[0][UniqChainRate].(string), defaultDecimalValue)
		gr.UniqBasisRate = utils.InterfaceStringToDecimal(sl[0][UniqBasisRate].(string), defaultDecimalValue)
	}

	return count, gr
}

func Add(gr GrowRate) int64 {
	sql := fmt.Sprintf("INSERT INTO " +
		tableName() + " " +
		"(" +
		tableGrowRateFields[MyInsName] + ", " +
		tableGrowRateFields[StatDate] + ", " +
		tableGrowRateFields[YesterdayTotal] + ", " +
		tableGrowRateFields[BeforeYesterdayTotal] + ", " +
		tableGrowRateFields[BasisTotal] + ", " +
		tableGrowRateFields[YesterdayUniq] + ", " +
		tableGrowRateFields[BeforeYesterdayUniq] + ", " +
		tableGrowRateFields[BasisUniq] + ", " +
		tableGrowRateFields[TotalChainRate] + ", " +
		tableGrowRateFields[TotalBasisRate] + ", " +
		tableGrowRateFields[UniqChainRate] + ", " +
		tableGrowRateFields[UniqBasisRate] + " " +
		") VALUES (" +
		fmt.Sprintf("'%s'", gr.MyInsName) + ", " +
		fmt.Sprintf("'%s 00:00:00'", utils.DateString(gr.StatDate.Format("2006-01-02"))) + ", " +
		fmt.Sprintf("%d", gr.YesterdayTotal) + ", " +
		fmt.Sprintf("%d", gr.BeforeYesterdayTotal) + ", " +
		fmt.Sprintf("%d", gr.BasisTotal) + ", " +
		fmt.Sprintf("%d", gr.YesterdayUniq) + ", " +
		fmt.Sprintf("%d", gr.BeforeYesterdayUniq) + ", " +
		fmt.Sprintf("%d", gr.BasisUniq) + ", " +
		fmt.Sprintf("%s", gr.TotalChainRate.String()) + ", " +
		fmt.Sprintf("%s", gr.TotalBasisRate.String()) + ", " +
		fmt.Sprintf("%s", gr.UniqChainRate.String()) + ", " +
		fmt.Sprintf("%s", gr.UniqBasisRate.String()) + " " +
		")")

	var rowsAffected int64 = 0
	retry := DbRetry()
	for i := 0; i < retry; i++ {
		p, err := orm.NewOrm().Raw(sql).Prepare()
		if err != nil {
			beego.Info(fmt.Sprintf("Add error(Prepare): %s", err.Error()))
			continue
		}
		res, err := p.Exec()
		if err != nil {
			beego.Info(fmt.Sprintf("Add error(Exec): %s", err.Error()))
			continue
		}

		i, err := res.RowsAffected()
		rowsAffected = i
		if err != nil {
			beego.Info(fmt.Sprintf("Add error(RowsAffected): %s", err.Error()))
			continue
		}
		break
	}

	return rowsAffected
}