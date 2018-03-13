package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"github.com/ximply/myslowreport/utils"
)

type Item struct {
	DbMax string
	UserMax string
	Sample string

	TsMin time.Time
	TsMax time.Time
	TsCnt float64

	QueryTimeSum float64
	QueryTimeMin float64
	QueryTimeMax float64
	QueryTimePct95 float64
	QueryTimeStddev float64
	QueryTimeMedian float64

	LockTimeSum float64
	LockTimeMin float64
	LockTimeMax float64
	LockTimePct95 float64
	LockTimeStddev float64
	LockTimeMedian float64

	RowsSentSum float64
	RowsSentMin float64
	RowsSentMax float64
	RowsSentPct95 float64
	RowsSentStddev float64
	RowsSentMedian float64

	RowsExaminedSum float64
	RowsExaminedMin float64
	RowsExaminedMax float64
	RowsExaminedPct95 float64
	RowsExaminedStddev float64
	RowsExaminedMedian float64
}

const (
	DbMax = iota
	UserMax
	Sample

	TsMin
	TsMax
	TsCnt

	QueryTimeSum
	QueryTimeMin
	QueryTimeMax
	QueryTimePct95
	QueryTimeStddev
	QueryTimeMedian

	LockTimeSum
	LockTimeMin
	LockTimeMax
	LockTimePct95
	LockTimeStddev
	LockTimeMedian

	RowsSentSum
	RowsSentMin
	RowsSentMax
	RowsSentPct95
	RowsSentStddev
	RowsSentMedian

	RowsExaminedSum
	RowsExaminedMin
	RowsExaminedMax
	RowsExaminedPct95
	RowsExaminedStddev
	RowsExaminedMedian

)

var tableFields = []string{
	"db_max",
	"user_max",
	"sample",

	"ts_min",
	"ts_max",
	"ts_cnt",

	"Query_time_sum",
	"Query_time_min",
	"Query_time_max",
	"Query_time_pct_95",
	"Query_time_stddev",
	"Query_time_median",

	"Lock_time_sum",
	"Lock_time_min",
	"Lock_time_max",
	"Lock_time_pct_95",
	"Lock_time_stddev",
	"Lock_time_median",

	"Rows_sent_sum",
	"Rows_sent_min",
	"Rows_sent_max",
	"Rows_sent_pct_95",
	"Rows_sent_stddev",
	"Rows_sent_median",

	"Rows_examined_sum",
	"Rows_examined_min",
	"Rows_examined_max",
	"Rows_examined_pct_95",
	"Rows_examined_stddev",
	"Rows_examined_median",

}

var defaultStringValue = "Unknown"
var defaultFloatValue = 0.0
var defaultTime = time.Now()

func GetSumOfQueryCount(since string, until string, table string) int64 {
	retry := DbRetry()
	sql := fmt.Sprintf("SELECT SUM(%s) FROM %s WHERE %s > '%s' AND %s < '%s'",
	tableFields[TsCnt], table, tableFields[TsMin], since, tableFields[TsMax], until)
	var sum int64
	var err error

	for i := 0; i < retry; i++ {
		sum = 0
		err = nil
		err = orm.NewOrm().Raw(sql).QueryRow(&sum)

		if err == nil {
			break
		} else {
			beego.Info(fmt.Sprintf("Get error when GetSumOfQueryCount %s: %s", table, err))
			continue
		}
	}
	
	return int64(sum)
}

func GetUniqOfQueryCount(since string, until string, table string) int64 {
	retry := DbRetry()
	sql := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE %s > '%s' AND %s < '%s'",
		table, tableFields[TsMin], since, tableFields[TsMax], until)
	var count int64
	var err error

	for i := 0; i < retry; i++ {
		count = 0
		err = nil
		err = orm.NewOrm().Raw(sql).QueryRow(&count)

		if err == nil {
			break
		} else {
			beego.Info(fmt.Sprintf("Get error when GetUniqOfQueryCount %s: %s", table, err))
			continue
		}
	}
	
	return int64(count)
}

func GetOrderByQueryTimeMaxDesc(since string, until string, table string) (int64, []Item) {
	retry := DbRetry()
	sql := fmt.Sprintf("SELECT " +

			tableFields[DbMax] + "," +
			tableFields[UserMax] + "," +
			tableFields[Sample]	+ "," +

			tableFields[TsMin] + "," +
			tableFields[TsMax] + "," +
			tableFields[TsCnt] + "," +

			tableFields[QueryTimeSum] + "," +
			tableFields[QueryTimeMin] + "," +
			tableFields[QueryTimeMax] + "," +
			tableFields[QueryTimePct95] + "," +
			tableFields[QueryTimeStddev] + "," +
			tableFields[QueryTimeMedian] + "," +

			tableFields[LockTimeSum] + "," +
			tableFields[LockTimeMin] + "," +
			tableFields[LockTimeMax] + "," +
			tableFields[LockTimePct95] + "," +
			tableFields[LockTimeStddev] + "," +
			tableFields[LockTimeMedian] + "," +

			tableFields[RowsSentSum] + "," +
			tableFields[RowsSentMin] + "," +
			tableFields[RowsSentMax] + "," +
			tableFields[RowsSentPct95] + "," +
			tableFields[RowsSentStddev] + "," +
			tableFields[RowsSentMedian] + "," +

			tableFields[RowsExaminedSum] + "," +
			tableFields[RowsExaminedMin] + "," +
			tableFields[RowsExaminedMax] + "," +
			tableFields[RowsExaminedPct95] + "," +
			tableFields[RowsExaminedStddev] + "," +
			tableFields[RowsExaminedMedian] + " " +

			" FROM %s WHERE ts_min > '%s' AND ts_max < '%s' ORDER BY Query_time_max DESC",

			table, since, until)

	var count int64
	var err error
	var sl []orm.ParamsList
	var il []Item

	for i := 0; i < retry; i++ {
		count = 0
		err = nil
		count, err = orm.NewOrm().Raw(sql).ValuesList(&sl,
			tableFields[DbMax],
			tableFields[UserMax],
			tableFields[Sample],

			tableFields[TsMin],
			tableFields[TsMax],
			tableFields[TsCnt],

			tableFields[QueryTimeSum],
			tableFields[QueryTimeMin],
			tableFields[QueryTimeMax],
			tableFields[QueryTimePct95],
			tableFields[QueryTimeStddev],
			tableFields[QueryTimeMedian],

			tableFields[LockTimeSum],
			tableFields[LockTimeMin],
			tableFields[LockTimeMax],
			tableFields[LockTimePct95],
			tableFields[LockTimeStddev],
			tableFields[LockTimeMedian],

			tableFields[RowsSentSum],
			tableFields[RowsSentMin],
			tableFields[RowsSentMax],
			tableFields[RowsSentPct95],
			tableFields[RowsSentStddev],
			tableFields[RowsSentMedian],

			tableFields[RowsExaminedSum],
			tableFields[RowsExaminedMin],
			tableFields[RowsExaminedMax],
			tableFields[RowsExaminedPct95],
			tableFields[RowsExaminedStddev],
			tableFields[RowsExaminedMedian])
		if err == nil {
			break
		} else {
			beego.Info(fmt.Sprintf("Get error when query %s: %s", table, err))
			continue
		}
	}

	if count > 0 {
		for _, s := range sl {
			var item Item

			if s[DbMax] != nil {
				item.DbMax = s[DbMax].(string)
			} else {
				item.DbMax = defaultStringValue
			}

			if s[UserMax] != nil {
				item.UserMax = s[UserMax].(string)
			} else {
				item.UserMax = defaultStringValue
			}

			if s[Sample] != nil {
				item.Sample = s[Sample].(string)
			} else {
				item.Sample = defaultStringValue
			}


			if s[TsMin] != nil {
				item.TsMin = utils.StringToTimeByFormat(s[TsMin].(string), "2006-01-02 15:04:05")
			} else {
				item.TsMin = defaultTime
			}

			if s[TsMax] != nil {
				item.TsMax = utils.StringToTimeByFormat(s[TsMax].(string), "2006-01-02 15:04:05")
			} else {
				item.TsMax = defaultTime
			}

			if s[TsCnt] != nil {
				item.TsCnt = utils.StringToFloat64(s[TsCnt].(string), defaultFloatValue)
			} else {
				item.TsCnt = defaultFloatValue
			}


			if s[QueryTimeSum] != nil {
				item.QueryTimeSum = utils.StringToFloat64(s[QueryTimeSum].(string), defaultFloatValue)
			} else {
				item.QueryTimeSum = defaultFloatValue
			}

			if s[QueryTimeMin] != nil {
				item.QueryTimeMin = utils.StringToFloat64(s[QueryTimeMin].(string), defaultFloatValue)
			} else {
				item.QueryTimeMin = defaultFloatValue
			}

			if s[QueryTimeMax] != nil {
				item.QueryTimeMax = utils.StringToFloat64(s[QueryTimeMax].(string), defaultFloatValue)
			} else {
				item.QueryTimeMax = defaultFloatValue
			}

			if s[QueryTimePct95] != nil {
				item.QueryTimePct95 = utils.StringToFloat64(s[QueryTimePct95].(string), defaultFloatValue)
			} else {
				item.QueryTimePct95 = defaultFloatValue
			}

			if s[QueryTimeStddev] != nil {
				item.QueryTimeStddev = utils.StringToFloat64(s[QueryTimeStddev].(string), defaultFloatValue)
			} else {
				item.QueryTimeStddev = defaultFloatValue
			}

			if s[QueryTimeMedian] != nil {
				item.QueryTimeMedian = utils.StringToFloat64(s[QueryTimeMedian].(string), defaultFloatValue)
			} else {
				item.QueryTimeMedian = defaultFloatValue
			}


			if s[LockTimeSum] != nil {
				item.LockTimeSum = utils.StringToFloat64(s[LockTimeSum].(string), defaultFloatValue)
			} else {
				item.LockTimeSum = defaultFloatValue
			}

			if s[LockTimeMin] != nil {
				item.LockTimeMin = utils.StringToFloat64(s[LockTimeMin].(string), defaultFloatValue)
			} else {
				item.LockTimeMin = defaultFloatValue
			}

			if s[LockTimeMax] != nil {
				item.LockTimeMax = utils.StringToFloat64(s[LockTimeMax].(string), defaultFloatValue)
			} else {
				item.LockTimeMax = defaultFloatValue
			}

			if s[LockTimePct95] != nil {
				item.LockTimePct95 = utils.StringToFloat64(s[LockTimePct95].(string), defaultFloatValue)
			} else {
				item.LockTimePct95 = defaultFloatValue
			}

			if s[LockTimeStddev] != nil {
				item.LockTimeStddev = utils.StringToFloat64(s[LockTimeStddev].(string), defaultFloatValue)
			} else {
				item.LockTimeStddev = defaultFloatValue
			}

			if s[LockTimeMedian] != nil {
				item.LockTimeMedian = utils.StringToFloat64(s[LockTimeMedian].(string), defaultFloatValue)
			} else {
				item.LockTimeMedian = defaultFloatValue
			}


			if s[RowsSentSum] != nil {
				item.RowsSentSum = utils.StringToFloat64(s[RowsSentSum].(string), defaultFloatValue)
			} else {
				item.RowsSentSum = defaultFloatValue
			}

			if s[RowsSentMin] != nil {
				item.RowsSentMin = utils.StringToFloat64(s[RowsSentMin].(string), defaultFloatValue)
			} else {
				item.RowsSentMin = defaultFloatValue
			}

			if s[RowsSentMax] != nil {
				item.RowsSentMax = utils.StringToFloat64(s[RowsSentMax].(string), defaultFloatValue)
			} else {
				item.RowsSentMax = defaultFloatValue
			}

			if s[RowsSentPct95] != nil {
				item.RowsSentPct95 = utils.StringToFloat64(s[RowsSentPct95].(string), defaultFloatValue)
			} else {
				item.RowsSentPct95 = defaultFloatValue
			}

			if s[RowsSentStddev] != nil {
				item.RowsSentStddev = utils.StringToFloat64(s[RowsSentStddev].(string), defaultFloatValue)
			} else {
				item.RowsSentStddev = defaultFloatValue
			}

			if s[RowsSentMedian] != nil {
				item.RowsSentMedian = utils.StringToFloat64(s[RowsSentMedian].(string), defaultFloatValue)
			} else {
				item.RowsSentMedian = defaultFloatValue
			}


			if s[RowsExaminedSum] != nil {
				item.RowsExaminedSum = utils.StringToFloat64(s[RowsExaminedSum].(string), defaultFloatValue)
			} else {
				item.RowsExaminedSum = defaultFloatValue
			}

			if s[RowsExaminedMin] != nil {
				item.RowsExaminedMin = utils.StringToFloat64(s[RowsExaminedMin].(string), defaultFloatValue)
			} else {
				item.RowsExaminedMin = defaultFloatValue
			}

			if s[RowsExaminedMax] != nil {
				item.RowsExaminedMax = utils.StringToFloat64(s[RowsExaminedMax].(string), defaultFloatValue)
			} else {
				item.RowsExaminedMax = defaultFloatValue
			}

			if s[RowsExaminedPct95] != nil {
				item.RowsExaminedPct95 = utils.StringToFloat64(s[RowsExaminedPct95].(string), defaultFloatValue)
			} else {
				item.RowsExaminedPct95 = defaultFloatValue
			}

			if s[RowsExaminedStddev] != nil {
				item.RowsExaminedStddev = utils.StringToFloat64(s[RowsExaminedStddev].(string), defaultFloatValue)
			} else {
				item.RowsExaminedStddev = defaultFloatValue
			}

			if s[RowsExaminedMedian] != nil {
				item.RowsExaminedMedian = utils.StringToFloat64(s[RowsExaminedMedian].(string), defaultFloatValue)
			} else {
				item.RowsExaminedMedian = defaultFloatValue
			}

			il = append(il, item)
		}
	}

	return count, il
}

func GetOrderByQueryCountDesc(since string, until string, table string) (int64, []Item) {
	retry := DbRetry()
	sql := fmt.Sprintf("SELECT " +

		tableFields[DbMax] + "," +
		tableFields[UserMax] + "," +
		tableFields[Sample]	+ "," +

		tableFields[TsMin] + "," +
		tableFields[TsMax] + "," +
		tableFields[TsCnt] + "," +

		tableFields[QueryTimeSum] + "," +
		tableFields[QueryTimeMin] + "," +
		tableFields[QueryTimeMax] + "," +
		tableFields[QueryTimePct95] + "," +
		tableFields[QueryTimeStddev] + "," +
		tableFields[QueryTimeMedian] + "," +

		tableFields[LockTimeSum] + "," +
		tableFields[LockTimeMin] + "," +
		tableFields[LockTimeMax] + "," +
		tableFields[LockTimePct95] + "," +
		tableFields[LockTimeStddev] + "," +
		tableFields[LockTimeMedian] + "," +

		tableFields[RowsSentSum] + "," +
		tableFields[RowsSentMin] + "," +
		tableFields[RowsSentMax] + "," +
		tableFields[RowsSentPct95] + "," +
		tableFields[RowsSentStddev] + "," +
		tableFields[RowsSentMedian] + "," +

		tableFields[RowsExaminedSum] + "," +
		tableFields[RowsExaminedMin] + "," +
		tableFields[RowsExaminedMax] + "," +
		tableFields[RowsExaminedPct95] + "," +
		tableFields[RowsExaminedStddev] + "," +
		tableFields[RowsExaminedMedian] + " " +

		" FROM %s WHERE ts_min > '%s' AND ts_max < '%s' ORDER BY ts_cnt DESC",

		table, since, until)

	var count int64
	var err error
	var sl []orm.ParamsList
	var il []Item

	for i := 0; i < retry; i++ {
		count = 0
		err = nil
		count, err = orm.NewOrm().Raw(sql).ValuesList(&sl,
			tableFields[DbMax],
			tableFields[UserMax],
			tableFields[Sample],

			tableFields[TsMin],
			tableFields[TsMax],
			tableFields[TsCnt],

			tableFields[QueryTimeSum],
			tableFields[QueryTimeMin],
			tableFields[QueryTimeMax],
			tableFields[QueryTimePct95],
			tableFields[QueryTimeStddev],
			tableFields[QueryTimeMedian],

			tableFields[LockTimeSum],
			tableFields[LockTimeMin],
			tableFields[LockTimeMax],
			tableFields[LockTimePct95],
			tableFields[LockTimeStddev],
			tableFields[LockTimeMedian],

			tableFields[RowsSentSum],
			tableFields[RowsSentMin],
			tableFields[RowsSentMax],
			tableFields[RowsSentPct95],
			tableFields[RowsSentStddev],
			tableFields[RowsSentMedian],

			tableFields[RowsExaminedSum],
			tableFields[RowsExaminedMin],
			tableFields[RowsExaminedMax],
			tableFields[RowsExaminedPct95],
			tableFields[RowsExaminedStddev],
			tableFields[RowsExaminedMedian])
		if err == nil {
			break
		} else {
			beego.Info(fmt.Sprintf("Get error when query %s: %s", table, err))
			continue
		}
	}

	if count > 0 {
		for _, s := range sl {
			var item Item

			if s[DbMax] != nil {
				item.DbMax = s[DbMax].(string)
			} else {
				item.DbMax = defaultStringValue
			}

			if s[UserMax] != nil {
				item.UserMax = s[UserMax].(string)
			} else {
				item.UserMax = defaultStringValue
			}

			if s[Sample] != nil {
				item.Sample = s[Sample].(string)
			} else {
				item.Sample = defaultStringValue
			}


			if s[TsMin] != nil {
				item.TsMin = utils.StringToTimeByFormat(s[TsMin].(string), "2006-01-02 15:04:05")
			} else {
				item.TsMin = defaultTime
			}

			if s[TsMax] != nil {
				item.TsMax = utils.StringToTimeByFormat(s[TsMax].(string), "2006-01-02 15:04:05")
			} else {
				item.TsMax = defaultTime
			}

			if s[TsCnt] != nil {
				item.TsCnt = utils.StringToFloat64(s[TsCnt].(string), defaultFloatValue)
			} else {
				item.TsCnt = defaultFloatValue
			}


			if s[QueryTimeSum] != nil {
				item.QueryTimeSum = utils.StringToFloat64(s[QueryTimeSum].(string), defaultFloatValue)
			} else {
				item.QueryTimeSum = defaultFloatValue
			}

			if s[QueryTimeMin] != nil {
				item.QueryTimeMin = utils.StringToFloat64(s[QueryTimeMin].(string), defaultFloatValue)
			} else {
				item.QueryTimeMin = defaultFloatValue
			}

			if s[QueryTimeMax] != nil {
				item.QueryTimeMax = utils.StringToFloat64(s[QueryTimeMax].(string), defaultFloatValue)
			} else {
				item.QueryTimeMax = defaultFloatValue
			}

			if s[QueryTimePct95] != nil {
				item.QueryTimePct95 = utils.StringToFloat64(s[QueryTimePct95].(string), defaultFloatValue)
			} else {
				item.QueryTimePct95 = defaultFloatValue
			}

			if s[QueryTimeStddev] != nil {
				item.QueryTimeStddev = utils.StringToFloat64(s[QueryTimeStddev].(string), defaultFloatValue)
			} else {
				item.QueryTimeStddev = defaultFloatValue
			}

			if s[QueryTimeMedian] != nil {
				item.QueryTimeMedian = utils.StringToFloat64(s[QueryTimeMedian].(string), defaultFloatValue)
			} else {
				item.QueryTimeMedian = defaultFloatValue
			}


			if s[LockTimeSum] != nil {
				item.LockTimeSum = utils.StringToFloat64(s[LockTimeSum].(string), defaultFloatValue)
			} else {
				item.LockTimeSum = defaultFloatValue
			}

			if s[LockTimeMin] != nil {
				item.LockTimeMin = utils.StringToFloat64(s[LockTimeMin].(string), defaultFloatValue)
			} else {
				item.LockTimeMin = defaultFloatValue
			}

			if s[LockTimeMax] != nil {
				item.LockTimeMax = utils.StringToFloat64(s[LockTimeMax].(string), defaultFloatValue)
			} else {
				item.LockTimeMax = defaultFloatValue
			}

			if s[LockTimePct95] != nil {
				item.LockTimePct95 = utils.StringToFloat64(s[LockTimePct95].(string), defaultFloatValue)
			} else {
				item.LockTimePct95 = defaultFloatValue
			}

			if s[LockTimeStddev] != nil {
				item.LockTimeStddev = utils.StringToFloat64(s[LockTimeStddev].(string), defaultFloatValue)
			} else {
				item.LockTimeStddev = defaultFloatValue
			}

			if s[LockTimeMedian] != nil {
				item.LockTimeMedian = utils.StringToFloat64(s[LockTimeMedian].(string), defaultFloatValue)
			} else {
				item.LockTimeMedian = defaultFloatValue
			}


			if s[RowsSentSum] != nil {
				item.RowsSentSum = utils.StringToFloat64(s[RowsSentSum].(string), defaultFloatValue)
			} else {
				item.RowsSentSum = defaultFloatValue
			}

			if s[RowsSentMin] != nil {
				item.RowsSentMin = utils.StringToFloat64(s[RowsSentMin].(string), defaultFloatValue)
			} else {
				item.RowsSentMin = defaultFloatValue
			}

			if s[RowsSentMax] != nil {
				item.RowsSentMax = utils.StringToFloat64(s[RowsSentMax].(string), defaultFloatValue)
			} else {
				item.RowsSentMax = defaultFloatValue
			}

			if s[RowsSentPct95] != nil {
				item.RowsSentPct95 = utils.StringToFloat64(s[RowsSentPct95].(string), defaultFloatValue)
			} else {
				item.RowsSentPct95 = defaultFloatValue
			}

			if s[RowsSentStddev] != nil {
				item.RowsSentStddev = utils.StringToFloat64(s[RowsSentStddev].(string), defaultFloatValue)
			} else {
				item.RowsSentStddev = defaultFloatValue
			}

			if s[RowsSentMedian] != nil {
				item.RowsSentMedian = utils.StringToFloat64(s[RowsSentMedian].(string), defaultFloatValue)
			} else {
				item.RowsSentMedian = defaultFloatValue
			}


			if s[RowsExaminedSum] != nil {
				item.RowsExaminedSum = utils.StringToFloat64(s[RowsExaminedSum].(string), defaultFloatValue)
			} else {
				item.RowsExaminedSum = defaultFloatValue
			}

			if s[RowsExaminedMin] != nil {
				item.RowsExaminedMin = utils.StringToFloat64(s[RowsExaminedMin].(string), defaultFloatValue)
			} else {
				item.RowsExaminedMin = defaultFloatValue
			}

			if s[RowsExaminedMax] != nil {
				item.RowsExaminedMax = utils.StringToFloat64(s[RowsExaminedMax].(string), defaultFloatValue)
			} else {
				item.RowsExaminedMax = defaultFloatValue
			}

			if s[RowsExaminedPct95] != nil {
				item.RowsExaminedPct95 = utils.StringToFloat64(s[RowsExaminedPct95].(string), defaultFloatValue)
			} else {
				item.RowsExaminedPct95 = defaultFloatValue
			}

			if s[RowsExaminedStddev] != nil {
				item.RowsExaminedStddev = utils.StringToFloat64(s[RowsExaminedStddev].(string), defaultFloatValue)
			} else {
				item.RowsExaminedStddev = defaultFloatValue
			}

			if s[RowsExaminedMedian] != nil {
				item.RowsExaminedMedian = utils.StringToFloat64(s[RowsExaminedMedian].(string), defaultFloatValue)
			} else {
				item.RowsExaminedMedian = defaultFloatValue
			}

			il = append(il, item)
		}
	}

	return count, il
}
