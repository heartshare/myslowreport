package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"github.com/ximply/myslowreport/utils"
	"github.com/shopspring/decimal"
)

type Item struct {
	DbMax string
	UserMax string
	Sample string

	TsMin time.Time
	TsMax time.Time
	TsCnt int64

	QueryTimeSum decimal.Decimal
	QueryTimeMin decimal.Decimal
	QueryTimeMax decimal.Decimal
	QueryTimePct95 decimal.Decimal
	QueryTimeStddev decimal.Decimal
	QueryTimeMedian decimal.Decimal

	LockTimeSum decimal.Decimal
	LockTimeMin decimal.Decimal
	LockTimeMax decimal.Decimal
	LockTimePct95 decimal.Decimal
	LockTimeStddev decimal.Decimal
	LockTimeMedian decimal.Decimal

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
	MaxQueryTimeMax decimal.Decimal
	MaxLockTimeMax decimal.Decimal
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
		mi.MaxTsCnt = utils.InterfaceStringToInt64(sl[0][MaxTsCnt], defaultIntValue)
		mi.MaxQueryTimeMax = utils.InterfaceStringToDecimal(sl[0][MaxQueryTimeMax], defaultDecimalValue)
		mi.MaxLockTimeMax = utils.InterfaceStringToDecimal(sl[0][MaxLockTimeMax], defaultDecimalValue)
		mi.MaxRowsSentMax = utils.InterfaceStringToInt64(sl[0][MaxRowsSentMax], defaultIntValue)
		mi.MaxRowsExaminedMax = utils.InterfaceStringToInt64(sl[0][MaxRowsExaminedMax], defaultIntValue)
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
		mi.MaxTsCnt = utils.InterfaceStringToInt64(sl[MaxTsCnt][MaxTsCnt], defaultIntValue)
		mi.MaxQueryTimeMax = utils.InterfaceStringToDecimal(sl[MaxQueryTimeMax][MaxTsCnt], defaultDecimalValue)
		utils.InterfaceStringToDecimal(sl[MaxLockTimeMax][MaxTsCnt], defaultDecimalValue)
		mi.MaxRowsSentMax = utils.InterfaceStringToInt64(sl[MaxRowsSentMax][MaxTsCnt], defaultIntValue)
		mi.MaxRowsExaminedMax = utils.InterfaceStringToInt64(sl[MaxRowsExaminedMax][MaxTsCnt], defaultIntValue)
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
			item.DbMax = utils.InterfaceStringToString(s[DbMax], defaultStringValue)
			item.UserMax = utils.InterfaceStringToString(s[UserMax], defaultStringValue)
			item.Sample = utils.InterfaceStringToString(s[Sample], defaultStringValue)

			item.TsMin = utils.InterfaceStringToTimeByFormat(s[TsMin],"2006-01-02 15:04:05", defaultTime)
			item.TsMax = utils.InterfaceStringToTimeByFormat(s[TsMax],"2006-01-02 15:04:05", defaultTime)
			item.TsCnt = utils.InterfaceStringToInt64(s[TsCnt], defaultIntValue)

			item.QueryTimeSum = utils.InterfaceStringToDecimal(s[QueryTimeSum], defaultDecimalValue)
			item.QueryTimeMin = utils.InterfaceStringToDecimal(s[QueryTimeMin], defaultDecimalValue)
			item.QueryTimeMax = utils.InterfaceStringToDecimal(s[QueryTimeMax], defaultDecimalValue)
			item.QueryTimePct95 = utils.InterfaceStringToDecimal(s[QueryTimePct95], defaultDecimalValue)
			item.QueryTimeStddev = utils.InterfaceStringToDecimal(s[QueryTimeStddev], defaultDecimalValue)
			item.QueryTimeMedian = utils.InterfaceStringToDecimal(s[QueryTimeMedian], defaultDecimalValue)

			item.LockTimeSum = utils.InterfaceStringToDecimal(s[LockTimeSum], defaultDecimalValue)
			item.LockTimeMin = utils.InterfaceStringToDecimal(s[LockTimeMin], defaultDecimalValue)
			item.LockTimeMax = utils.InterfaceStringToDecimal(s[LockTimeMax], defaultDecimalValue)
			item.LockTimePct95 = utils.InterfaceStringToDecimal(s[LockTimePct95], defaultDecimalValue)
			item.LockTimeStddev = utils.InterfaceStringToDecimal(s[LockTimeStddev], defaultDecimalValue)
			item.LockTimeMedian = utils.InterfaceStringToDecimal(s[LockTimeMedian], defaultDecimalValue)

			item.RowsSentSum = utils.InterfaceStringToInt64(s[RowsSentSum], defaultIntValue)
			item.RowsSentMin = utils.InterfaceStringToInt64(s[RowsSentMin], defaultIntValue)
			item.RowsSentMax = utils.InterfaceStringToInt64(s[RowsSentMax], defaultIntValue)
			item.RowsSentPct95 = utils.InterfaceStringToInt64(s[RowsSentPct95], defaultIntValue)
			item.RowsSentStddev = utils.InterfaceStringToInt64(s[RowsSentStddev], defaultIntValue)
			item.RowsSentMedian = utils.InterfaceStringToInt64(s[RowsSentMedian], defaultIntValue)

			item.RowsExaminedSum = utils.InterfaceStringToInt64(s[RowsExaminedSum], defaultIntValue)
			item.RowsExaminedMin = utils.InterfaceStringToInt64(s[RowsExaminedMin], defaultIntValue)
			item.RowsExaminedMax = utils.InterfaceStringToInt64(s[RowsExaminedMax], defaultIntValue)
			item.RowsExaminedPct95 = utils.InterfaceStringToInt64(s[RowsExaminedPct95], defaultIntValue)
			item.RowsExaminedStddev = utils.InterfaceStringToInt64(s[RowsExaminedStddev], defaultIntValue)
			item.RowsExaminedMedian = utils.InterfaceStringToInt64(s[RowsExaminedMedian], defaultIntValue)

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
			item.DbMax = utils.InterfaceStringToString(s[DbMax], defaultStringValue)
			item.UserMax = utils.InterfaceStringToString(s[UserMax], defaultStringValue)
			item.Sample = utils.InterfaceStringToString(s[Sample], defaultStringValue)

			item.TsMin = utils.InterfaceStringToTimeByFormat(s[TsMin],"2006-01-02 15:04:05", defaultTime)
			item.TsMax = utils.InterfaceStringToTimeByFormat(s[TsMax],"2006-01-02 15:04:05", defaultTime)
			item.TsCnt = utils.InterfaceStringToInt64(s[TsCnt], defaultIntValue)

			item.QueryTimeSum = utils.InterfaceStringToDecimal(s[QueryTimeSum], defaultDecimalValue)
			item.QueryTimeMin = utils.InterfaceStringToDecimal(s[QueryTimeMin], defaultDecimalValue)
			item.QueryTimeMax = utils.InterfaceStringToDecimal(s[QueryTimeMax], defaultDecimalValue)
			item.QueryTimePct95 = utils.InterfaceStringToDecimal(s[QueryTimePct95], defaultDecimalValue)
			item.QueryTimeStddev = utils.InterfaceStringToDecimal(s[QueryTimeStddev], defaultDecimalValue)
			item.QueryTimeMedian = utils.InterfaceStringToDecimal(s[QueryTimeMedian], defaultDecimalValue)

			item.LockTimeSum = utils.InterfaceStringToDecimal(s[LockTimeSum], defaultDecimalValue)
			item.LockTimeMin = utils.InterfaceStringToDecimal(s[LockTimeMin], defaultDecimalValue)
			item.LockTimeMax = utils.InterfaceStringToDecimal(s[LockTimeMax], defaultDecimalValue)
			item.LockTimePct95 = utils.InterfaceStringToDecimal(s[LockTimePct95], defaultDecimalValue)
			item.LockTimeStddev = utils.InterfaceStringToDecimal(s[LockTimeStddev], defaultDecimalValue)
			item.LockTimeMedian = utils.InterfaceStringToDecimal(s[LockTimeMedian], defaultDecimalValue)

			item.RowsSentSum = utils.InterfaceStringToInt64(s[RowsSentSum], defaultIntValue)
			item.RowsSentMin = utils.InterfaceStringToInt64(s[RowsSentMin], defaultIntValue)
			item.RowsSentMax = utils.InterfaceStringToInt64(s[RowsSentMax], defaultIntValue)
			item.RowsSentPct95 = utils.InterfaceStringToInt64(s[RowsSentPct95], defaultIntValue)
			item.RowsSentStddev = utils.InterfaceStringToInt64(s[RowsSentStddev], defaultIntValue)
			item.RowsSentMedian = utils.InterfaceStringToInt64(s[RowsSentMedian], defaultIntValue)

			item.RowsExaminedSum = utils.InterfaceStringToInt64(s[RowsExaminedSum], defaultIntValue)
			item.RowsExaminedMin = utils.InterfaceStringToInt64(s[RowsExaminedMin], defaultIntValue)
			item.RowsExaminedMax = utils.InterfaceStringToInt64(s[RowsExaminedMax], defaultIntValue)
			item.RowsExaminedPct95 = utils.InterfaceStringToInt64(s[RowsExaminedPct95], defaultIntValue)
			item.RowsExaminedStddev = utils.InterfaceStringToInt64(s[RowsExaminedStddev], defaultIntValue)
			item.RowsExaminedMedian = utils.InterfaceStringToInt64(s[RowsExaminedMedian], defaultIntValue)

			il = append(il, item)
		}
	}

	return count, il
}

