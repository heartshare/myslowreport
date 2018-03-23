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

type GrowRateMonthly struct {
	Id int64

	MyInsName string
	StatDate time.Time

	LastMonthTotal int64
	LastLastMonthTotal int64

	LastMonthUniq int64
	LastLastMonthUniq int64

	TotalChainRate decimal.Decimal
	UniqChainRate decimal.Decimal
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

const (
	IdM = iota

	MyInsNameM
	StatDateM

	LastMonthTotalM
	LastLastMonthTotalM

	LastMonthUniqM
	LastLastMonthUniqM

	TotalChainRateM
	UniqChainRateM
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

var tableGrowRateMonthlyFields = []string{
	"Id",

	"MyInsName",
	"StatDate",

	"LastMonthTotal",
	"LastLastMonthTotal",

	"LastMonthUniq",
	"LastLastMonthUniq",

	"TotalChainRate",
	"UniqChainRate",
}

func tableName() string {
	return "myslow_history_grow_rate"
}

func tableMonthlyName() string {
	return "myslow_history_grow_rate_monthly"
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

func GetByMyInsNameAndStatDateMonthly(name string, date string) (int64, GrowRateMonthly) {
	retry := DbRetry()
	sql := fmt.Sprintf("SELECT " +
		tableGrowRateMonthlyFields[IdM] + ", " +
		tableGrowRateMonthlyFields[MyInsNameM] + ", " +
		tableGrowRateMonthlyFields[StatDateM] + ", " +
		tableGrowRateMonthlyFields[LastMonthTotalM] + ", " +
		tableGrowRateMonthlyFields[LastLastMonthTotalM] + ", " +
		tableGrowRateMonthlyFields[LastMonthUniqM] + ", " +
		tableGrowRateMonthlyFields[LastLastMonthUniqM] + ", " +
		tableGrowRateMonthlyFields[TotalChainRateM] + ", " +
		tableGrowRateMonthlyFields[UniqChainRateM] + " " +
		" FROM %s WHERE %s = '%s' AND %s = '%s'",

		tableMonthlyName(), tableGrowRateMonthlyFields[MyInsNameM], name, tableGrowRateMonthlyFields[StatDateM], date)

	var count int64
	var err error
	var sl []orm.ParamsList
	var grm GrowRateMonthly
	fmt.Println(sql)
	for i := 0; i < retry; i++ {
		count = 0
		err = nil
		count, err = orm.NewOrm().Raw(sql).ValuesList(&sl,
			tableGrowRateMonthlyFields[IdM],
			tableGrowRateMonthlyFields[MyInsNameM],
			tableGrowRateMonthlyFields[StatDateM],
			tableGrowRateMonthlyFields[LastMonthTotalM],
			tableGrowRateMonthlyFields[LastLastMonthTotalM],
			tableGrowRateMonthlyFields[LastMonthUniqM],
			tableGrowRateMonthlyFields[LastLastMonthUniqM],
			tableGrowRateMonthlyFields[TotalChainRateM],
			tableGrowRateMonthlyFields[UniqChainRateM])
		if err == nil {
			break
		} else {
			beego.Info(fmt.Sprintf("Get error when query %s: %s", tableMonthlyName(), err))
			continue
		}
	}

	if count > 0 {
		grm.Id = utils.InterfaceStringToInt64(sl[0][IdM], defaultIntValue)
		grm.MyInsName = utils.InterfaceStringToString(sl[0][MyInsNameM], defaultStringValue)
		grm.StatDate = utils.InterfaceStringToTimeByFormat(sl[0][StatDateM], "2006-01-02", defaultTime)

		grm.LastMonthTotal = utils.InterfaceStringToInt64(sl[0][LastMonthTotalM], defaultIntValue)
		grm.LastLastMonthTotal = utils.InterfaceStringToInt64(sl[0][LastLastMonthTotalM], defaultIntValue)

		grm.LastMonthUniq = utils.InterfaceStringToInt64(sl[0][LastMonthUniqM], defaultIntValue)
		grm.LastLastMonthUniq = utils.InterfaceStringToInt64(sl[0][LastLastMonthUniqM], defaultIntValue)

		grm.TotalChainRate = utils.InterfaceStringToDecimal(sl[0][TotalChainRateM], defaultDecimalValue)
		grm.UniqChainRate = utils.InterfaceStringToDecimal(sl[0][UniqChainRateM].(string), defaultDecimalValue)
	}

	return count, grm
}

func AddGrowRate(gr GrowRate) int64 {
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

func AddGrowRateMonthly(grm GrowRateMonthly) int64 {
	sql := fmt.Sprintf("INSERT INTO " +
		tableMonthlyName() + " " +
		"(" +
		tableGrowRateMonthlyFields[MyInsNameM] + ", " +
		tableGrowRateMonthlyFields[StatDateM] + ", " +
		tableGrowRateMonthlyFields[LastMonthTotalM] + ", " +
		tableGrowRateMonthlyFields[LastLastMonthTotalM] + ", " +
		tableGrowRateMonthlyFields[LastMonthUniqM] + ", " +
		tableGrowRateMonthlyFields[LastLastMonthUniqM] + ", " +
		tableGrowRateMonthlyFields[TotalChainRateM] + ", " +
		tableGrowRateMonthlyFields[UniqChainRateM] + " " +
		") VALUES (" +
		fmt.Sprintf("'%s'", grm.MyInsName) + ", " +
		fmt.Sprintf("'%s-01 00:00:00'", utils.YearMonthStringByFormat(grm.StatDate, "2006-01-02")) + ", " +
		fmt.Sprintf("%d", grm.LastMonthTotal) + ", " +
		fmt.Sprintf("%d", grm.LastLastMonthTotal) + ", " +
		fmt.Sprintf("%d", grm.LastMonthUniq) + ", " +
		fmt.Sprintf("%d", grm.LastLastMonthUniq) + ", " +
		fmt.Sprintf("%s", grm.TotalChainRate.String()) + ", " +
		fmt.Sprintf("%s", grm.UniqChainRate.String()) + " " +
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

func GetLast30DaysCountInfo(name string) (int64, map[string]string) {
	m := make(map[string]string)
	from := time.Now().AddDate(0, -1, -3)
	retry := DbRetry()
	sql := fmt.Sprintf("SELECT " +
		tableGrowRateFields[StatDate] + ", " +
		tableGrowRateFields[YesterdayTotal] + ", " +
		tableGrowRateFields[YesterdayUniq] + " " +
		" FROM %s WHERE %s = '%s' AND %s > '%s'",

		tableName(), tableGrowRateFields[MyInsName],
		name, tableGrowRateFields[StatDate], from.Format("2006-01-02"))

	var count int64
	var err error
	var sl []orm.ParamsList

	for i := 0; i < retry; i++ {
		count = 0
		err = nil
		count, err = orm.NewOrm().Raw(sql).ValuesList(&sl,
			tableGrowRateFields[StatDate],
			tableGrowRateFields[YesterdayTotal],
			tableGrowRateFields[YesterdayUniq])
		if err == nil {
			break
		} else {
			beego.Info(fmt.Sprintf("Get error when query %s: %s", tableName(), err))
			continue
		}
	}

	for _, s := range sl {
		date := utils.DateString(utils.InterfaceStringToString(s[0], defaultStringValue))
		m[date] = fmt.Sprintf("%d,%d", utils.InterfaceStringToInt64(s[1], defaultIntValue),
			utils.InterfaceStringToInt64(s[2], defaultIntValue))
	}

	return count, m
}