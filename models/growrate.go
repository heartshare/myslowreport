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

	fmt.Println(sql)

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
		if sl[0][Id] != nil {
			gr.Id = utils.StringToInt64(sl[0][Id].(string), defaultIntValue)
		} else {
			gr.Id = defaultIntValue
		}

		if sl[0][MyInsName] != nil {
			gr.MyInsName = sl[0][MyInsName].(string)
		} else {
			gr.MyInsName  = defaultStringValue
		}
		if sl[0][StatDate] != nil {
			gr.StatDate = utils.StringToTimeByFormat(sl[0][StatDate].(string), "2006-01-02")
		} else {
			gr.StatDate = defaultTime
		}

		if sl[0][YesterdayTotal] != nil {
			gr.YesterdayTotal = utils.StringToInt64(sl[0][YesterdayTotal].(string), defaultIntValue)
		} else {
			gr.YesterdayTotal = defaultIntValue
		}
		if sl[0][BeforeYesterdayTotal] != nil {
			gr.BeforeYesterdayTotal = utils.StringToInt64(sl[0][BeforeYesterdayTotal].(string), defaultIntValue)
		} else {
			gr.BeforeYesterdayTotal = defaultIntValue
		}
		if sl[0][BasisTotal] != nil {
			gr.BasisTotal = utils.StringToInt64(sl[0][BasisTotal].(string), defaultIntValue)
		} else {
			gr.BasisTotal = defaultIntValue
		}

		if sl[0][YesterdayUniq] != nil {
			gr.YesterdayUniq = utils.StringToInt64(sl[0][YesterdayUniq].(string), defaultIntValue)
		} else {
			gr.YesterdayUniq = defaultIntValue
		}
		if sl[0][BeforeYesterdayUniq] != nil {
			gr.BeforeYesterdayUniq = utils.StringToInt64(sl[0][BeforeYesterdayUniq].(string), defaultIntValue)
		} else {
			gr.BeforeYesterdayUniq = defaultIntValue
		}
		if sl[0][BasisUniq] != nil {
			gr.BasisUniq = utils.StringToInt64(sl[0][BasisUniq].(string), defaultIntValue)
		} else {
			gr.BasisUniq = defaultIntValue
		}

		if sl[0][TotalChainRate] != nil {
			gr.TotalChainRate = utils.StringToDecimal(sl[0][TotalChainRate].(string), defaultDecimalValue)
		} else {
			gr.TotalChainRate = defaultDecimalValue
		}
		if sl[0][TotalBasisRate] != nil {
			gr.TotalBasisRate = utils.StringToDecimal(sl[0][TotalBasisRate].(string), defaultDecimalValue)
		} else {
			gr.TotalBasisRate = defaultDecimalValue
		}

		if sl[0][UniqChainRate] != nil {
			gr.UniqChainRate = utils.StringToDecimal(sl[0][UniqChainRate].(string), defaultDecimalValue)
		} else {
			gr.UniqChainRate = defaultDecimalValue
		}
		if sl[0][UniqBasisRate] != nil {
			gr.UniqBasisRate = utils.StringToDecimal(sl[0][UniqBasisRate].(string), defaultDecimalValue)
		} else {
			gr.UniqBasisRate = defaultDecimalValue
		}
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

	fmt.Println(sql)

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