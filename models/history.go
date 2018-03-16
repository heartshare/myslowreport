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
	TsCnt int64

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

	RowsSentSum int64
	RowsSentMin int64
	RowsSentMax int64
	RowsSentPct95 int64
	RowsSentStddev int64
	RowsSentMedian int64

	RowsExaminedSum int64
	RowsExaminedMin int64
	RowsExaminedMax int64
	RowsExaminedPct95 int64
	RowsExaminedStddev int64
	RowsExaminedMedian int64
}

type MaxItem struct {
	MaxTsCnt int64
	MaxQueryTimeMax float64
	MaxLockTimeMax float64
	MaxRowsSentMax int64
	MaxRowsExaminedMax int64
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

const (
	MaxTsCnt = iota
	MaxQueryTimeMax
	MaxLockTimeMax
	MaxRowsSentMax
	MaxRowsExaminedMax
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

var maxTableFields = []string{
	"max_ts_cnt",
	"max_Query_time_max",
	"max_Lock_time_max",
	"max_Rows_sent_max",
	"max_Rows_examined_max",
}

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

func GetMax(since string, until string, table string) (int64, MaxItem) {
	retry := DbRetry()
	sql := fmt.Sprintf("SELECT " +

		"MAX(" + tableFields[TsCnt] + ") AS " + maxTableFields[MaxTsCnt] + "," +
		"MAX(" + tableFields[QueryTimeMax] + ") AS " + maxTableFields[MaxQueryTimeMax] + "," +
		"MAX(" + tableFields[LockTimeMax] + ") AS " + maxTableFields[MaxLockTimeMax] + "," +
		"MAX(" + tableFields[RowsSentMax] + ") AS " + maxTableFields[MaxRowsSentMax] + "," +
		"MAX(" + tableFields[RowsExaminedMax] + ") AS " + maxTableFields[MaxRowsExaminedMax] +

		" FROM %s WHERE ts_min > '%s' AND ts_max < '%s'",

		table, since, until)

	var count int64
	var err error
	var sl []orm.ParamsList
	var mi MaxItem

	for i := 0; i < retry; i++ {
		count = 0
		err = nil
		count, err = orm.NewOrm().Raw(sql).ValuesList(&sl,
			maxTableFields[MaxTsCnt],
			maxTableFields[MaxQueryTimeMax],
			maxTableFields[MaxLockTimeMax],
			maxTableFields[MaxRowsSentMax],
			maxTableFields[MaxRowsExaminedMax])
		if err == nil {
			break
		} else {
			beego.Info(fmt.Sprintf("Get error when query %s: %s", table, err))
			continue
		}
	}

	if count > 0 {
		if sl[0][MaxTsCnt] != nil {
			mi.MaxTsCnt = utils.StringToInt64(sl[0][MaxTsCnt].(string), defaultIntValue)
		} else {
			mi.MaxTsCnt = defaultIntValue
		}

		if sl[0][MaxQueryTimeMax] != nil {
			mi.MaxQueryTimeMax = utils.StringToFloat64(sl[0][MaxQueryTimeMax].(string), defaultFloatValue)
		} else {
			mi.MaxQueryTimeMax = defaultFloatValue
		}

		if sl[0][MaxLockTimeMax] != nil {
			mi.MaxLockTimeMax = utils.StringToFloat64(sl[0][MaxLockTimeMax].(string), defaultFloatValue)
		} else {
			mi.MaxLockTimeMax = defaultFloatValue
		}

		if sl[0][MaxRowsSentMax] != nil {
			mi.MaxRowsSentMax = utils.StringToInt64(sl[0][MaxRowsSentMax].(string), defaultIntValue)
		} else {
			mi.MaxRowsSentMax = defaultIntValue
		}

		if sl[0][MaxRowsExaminedMax] != nil {
			mi.MaxRowsExaminedMax = utils.StringToInt64(sl[0][MaxRowsExaminedMax].(string), defaultIntValue)
		} else {
			mi.MaxRowsExaminedMax = defaultIntValue
		}
	}

	return count, mi
}

func GetMaxOrderBy(since string, until string, table string) (int64, MaxItem) {
	retry := DbRetry()

	sql := fmt.Sprintf("SELECT * FROM (SELECT %s FROM %s WHERE ts_min > '%s' AND ts_max < '%s' ORDER BY %s DESC LIMIT 1) AS t1",
		tableFields[TsCnt], table, since, until, tableFields[TsCnt])
	sql += " UNION "
	sql += fmt.Sprintf("SELECT * FROM (SELECT %s FROM %s WHERE ts_min > '%s' AND ts_max < '%s' ORDER BY %s DESC LIMIT 1) AS t2",
		tableFields[QueryTimeMax], table, since, until, tableFields[QueryTimeMax])
	sql += " UNION "
	sql += fmt.Sprintf("SELECT * FROM (SELECT %s FROM %s WHERE ts_min > '%s' AND ts_max < '%s' ORDER BY %s DESC LIMIT 1) AS t3",
		tableFields[LockTimeMax], table, since, until, tableFields[LockTimeMax])
	sql += " UNION "
	sql += fmt.Sprintf("SELECT * FROM (SELECT %s FROM %s WHERE ts_min > '%s' AND ts_max < '%s' ORDER BY %s DESC LIMIT 1) AS t4",
		tableFields[RowsSentMax], table, since, until, tableFields[RowsSentMax])
	sql += " UNION "
	sql += fmt.Sprintf("SELECT * FROM (SELECT %s FROM %s WHERE ts_min > '%s' AND ts_max < '%s' ORDER BY %s DESC LIMIT 1) AS t5",
		tableFields[MaxRowsExaminedMax], table, since, until, tableFields[MaxRowsExaminedMax])

	var count int64
	var err error
	var sl []orm.ParamsList
	var mi MaxItem

	for i := 0; i < retry; i++ {
		count = 0
		err = nil
		count, err = orm.NewOrm().Raw(sql).ValuesList(&sl,
			"ts_cnt")
		if err == nil {
			break
		} else {
			beego.Info(fmt.Sprintf("Get error when query %s: %s", table, err))
			continue
		}
	}

	if count == 5 {
		if sl[MaxTsCnt][MaxTsCnt] != nil {
			mi.MaxTsCnt = utils.StringToInt64(sl[MaxTsCnt][MaxTsCnt].(string), defaultIntValue)
		} else {
			mi.MaxTsCnt = defaultIntValue
		}

		if sl[MaxQueryTimeMax][MaxTsCnt] != nil {
			mi.MaxQueryTimeMax = utils.StringToFloat64(sl[MaxQueryTimeMax][MaxTsCnt].(string), defaultFloatValue)
		} else {
			mi.MaxQueryTimeMax = defaultFloatValue
		}

		if sl[MaxLockTimeMax][MaxTsCnt] != nil {
			mi.MaxLockTimeMax = utils.StringToFloat64(sl[MaxLockTimeMax][MaxTsCnt].(string), defaultFloatValue)
		} else {
			mi.MaxLockTimeMax = defaultFloatValue
		}

		if sl[MaxRowsSentMax][MaxTsCnt] != nil {
			mi.MaxRowsSentMax = utils.StringToInt64(sl[MaxRowsSentMax][MaxTsCnt].(string), defaultIntValue)
		} else {
			mi.MaxRowsSentMax = defaultIntValue
		}

		if sl[MaxRowsExaminedMax][MaxTsCnt] != nil {
			mi.MaxRowsExaminedMax = utils.StringToInt64(sl[MaxRowsExaminedMax][MaxTsCnt].(string), defaultIntValue)
		} else {
			mi.MaxRowsExaminedMax = defaultIntValue
		}
	}

	return count, mi
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
				item.TsCnt = utils.StringToInt64(s[TsCnt].(string), defaultIntValue)
			} else {
				item.TsCnt = defaultIntValue
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
				item.RowsSentSum = utils.StringToInt64(s[RowsSentSum].(string), defaultIntValue)
			} else {
				item.RowsSentSum = defaultIntValue
			}

			if s[RowsSentMin] != nil {
				item.RowsSentMin = utils.StringToInt64(s[RowsSentMin].(string), defaultIntValue)
			} else {
				item.RowsSentMin = defaultIntValue
			}

			if s[RowsSentMax] != nil {
				item.RowsSentMax = utils.StringToInt64(s[RowsSentMax].(string), defaultIntValue)
			} else {
				item.RowsSentMax = defaultIntValue
			}

			if s[RowsSentPct95] != nil {
				item.RowsSentPct95 = utils.StringToInt64(s[RowsSentPct95].(string), defaultIntValue)
			} else {
				item.RowsSentPct95 = defaultIntValue
			}

			if s[RowsSentStddev] != nil {
				item.RowsSentStddev = utils.StringToInt64(s[RowsSentStddev].(string), defaultIntValue)
			} else {
				item.RowsSentStddev = defaultIntValue
			}

			if s[RowsSentMedian] != nil {
				item.RowsSentMedian = utils.StringToInt64(s[RowsSentMedian].(string), defaultIntValue)
			} else {
				item.RowsSentMedian = defaultIntValue
			}


			if s[RowsExaminedSum] != nil {
				item.RowsExaminedSum = utils.StringToInt64(s[RowsExaminedSum].(string), defaultIntValue)
			} else {
				item.RowsExaminedSum = defaultIntValue
			}

			if s[RowsExaminedMin] != nil {
				item.RowsExaminedMin = utils.StringToInt64(s[RowsExaminedMin].(string), defaultIntValue)
			} else {
				item.RowsExaminedMin = defaultIntValue
			}

			if s[RowsExaminedMax] != nil {
				item.RowsExaminedMax = utils.StringToInt64(s[RowsExaminedMax].(string), defaultIntValue)
			} else {
				item.RowsExaminedMax = defaultIntValue
			}

			if s[RowsExaminedPct95] != nil {
				item.RowsExaminedPct95 = utils.StringToInt64(s[RowsExaminedPct95].(string), defaultIntValue)
			} else {
				item.RowsExaminedPct95 = defaultIntValue
			}

			if s[RowsExaminedStddev] != nil {
				item.RowsExaminedStddev = utils.StringToInt64(s[RowsExaminedStddev].(string), defaultIntValue)
			} else {
				item.RowsExaminedStddev = defaultIntValue
			}

			if s[RowsExaminedMedian] != nil {
				item.RowsExaminedMedian = utils.StringToInt64(s[RowsExaminedMedian].(string), defaultIntValue)
			} else {
				item.RowsExaminedMedian = defaultIntValue
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
				item.TsCnt = utils.StringToInt64(s[TsCnt].(string), defaultIntValue)
			} else {
				item.TsCnt = defaultIntValue
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
				item.RowsSentSum = utils.StringToInt64(s[RowsSentSum].(string), defaultIntValue)
			} else {
				item.RowsSentSum = defaultIntValue
			}

			if s[RowsSentMin] != nil {
				item.RowsSentMin = utils.StringToInt64(s[RowsSentMin].(string), defaultIntValue)
			} else {
				item.RowsSentMin = defaultIntValue
			}

			if s[RowsSentMax] != nil {
				item.RowsSentMax = utils.StringToInt64(s[RowsSentMax].(string), defaultIntValue)
			} else {
				item.RowsSentMax = defaultIntValue
			}

			if s[RowsSentPct95] != nil {
				item.RowsSentPct95 = utils.StringToInt64(s[RowsSentPct95].(string), defaultIntValue)
			} else {
				item.RowsSentPct95 = defaultIntValue
			}

			if s[RowsSentStddev] != nil {
				item.RowsSentStddev = utils.StringToInt64(s[RowsSentStddev].(string), defaultIntValue)
			} else {
				item.RowsSentStddev = defaultIntValue
			}

			if s[RowsSentMedian] != nil {
				item.RowsSentMedian = utils.StringToInt64(s[RowsSentMedian].(string), defaultIntValue)
			} else {
				item.RowsSentMedian = defaultIntValue
			}


			if s[RowsExaminedSum] != nil {
				item.RowsExaminedSum = utils.StringToInt64(s[RowsExaminedSum].(string), defaultIntValue)
			} else {
				item.RowsExaminedSum = defaultIntValue
			}

			if s[RowsExaminedMin] != nil {
				item.RowsExaminedMin = utils.StringToInt64(s[RowsExaminedMin].(string), defaultIntValue)
			} else {
				item.RowsExaminedMin = defaultIntValue
			}

			if s[RowsExaminedMax] != nil {
				item.RowsExaminedMax = utils.StringToInt64(s[RowsExaminedMax].(string), defaultIntValue)
			} else {
				item.RowsExaminedMax = defaultIntValue
			}

			if s[RowsExaminedPct95] != nil {
				item.RowsExaminedPct95 = utils.StringToInt64(s[RowsExaminedPct95].(string), defaultIntValue)
			} else {
				item.RowsExaminedPct95 = defaultIntValue
			}

			if s[RowsExaminedStddev] != nil {
				item.RowsExaminedStddev = utils.StringToInt64(s[RowsExaminedStddev].(string), defaultIntValue)
			} else {
				item.RowsExaminedStddev = defaultIntValue
			}

			if s[RowsExaminedMedian] != nil {
				item.RowsExaminedMedian = utils.StringToInt64(s[RowsExaminedMedian].(string), defaultIntValue)
			} else {
				item.RowsExaminedMedian = defaultIntValue
			}

			il = append(il, item)
		}
	}

	return count, il
}

