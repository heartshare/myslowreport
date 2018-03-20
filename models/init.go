package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"strings"
	"fmt"
	"time"
	"github.com/shopspring/decimal"
)

var defaultStringValue = "Unknown"
var defaultFloatValue = 0.0
var defaultIntValue int64 = 0
var defaultTime = time.Now()
var defaultDecimalValue, _ = decimal.NewFromString("00.00")

type Project struct {
	MysqlHost string
	MysqlPort string
	SlowlogTable string
	Description string
	RsyncModel string
	MysqlSlowlogFileName string
}

type MyInsChargeInfo struct {
	MyInsName string
	Desc string
	Charge string
}

type TraceSheetInfo struct {
	SheetName string
	MyInsChargeInfoList []MyInsChargeInfo
}

type TraceColInfo struct {
	ColName string
	ColCellName string
	ColCellAlapha string
}

func Init() {
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn)

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}

func DbRetry() int {
	retry, err := beego.AppConfig.Int("db.retry")
	if err != nil {
		retry = 2
	}
	return retry
}

func MyslowReportDbHost() string {
	return beego.AppConfig.String("db.host")
}

func MyslowReportDbPort() string {
	return beego.AppConfig.String("db.port")
}

func MyslowReportDbUser() string {
	return beego.AppConfig.String("db.user")
}

func MyslowReportDbPassword() string {
	return beego.AppConfig.String("db.password")
}

func MyslowReportDbName() string {
	return beego.AppConfig.String("db.name")
}

func MyslowReportTimeout() int {
	t, err := beego.AppConfig.Int("myslowreport.timeout")
	if err != nil {
		t = 10
	}
	return t
}

func SyncMysqlSlowlogSpec() string {
	return beego.AppConfig.String("myslowreport.syncmysqlslowlogspec")
}

func SendMysqlSlowlogReportSpec() string {
	return beego.AppConfig.String("myslowreport.sendreportemailspec")
}

func ImportMysqlSlowlogMonthlySpec() string {
	return beego.AppConfig.String("myslowreport.importmonthlyspec")
}

func MyslowReportRetry() int {
	retry, err := beego.AppConfig.Int("myslowreport.retry")
	if err != nil {
		retry = 2
	}
	return retry
}

func MyslowReportTitle() string {
	return beego.AppConfig.String("myslowreport.title")
}

func MyslowReportSubject() string {
	return beego.AppConfig.String("myslowreport.subject")
}

func MyslowReportFrom() string {
	return beego.AppConfig.String("myslowreport.from")
}

func MyslowReportFromAlias() string {
	return beego.AppConfig.String("myslowreport.fromalias")
}

func MyslowReportMailUserName() string {
	return beego.AppConfig.String("myslowreport.mailusername")
}

func MyslowReportMailPassword() string {
	return beego.AppConfig.String("myslowreport.mailpassword")
}

func MyslowReportMailHost() string {
	return beego.AppConfig.String("myslowreport.mailhost")
}

func MyslowReportMailPort() string {
	return beego.AppConfig.String("myslowreport.mailport")
}

func MyslowReportMailTemplate() string {
	return beego.AppConfig.String("myslowreport.mailtemplate")
}

func myslowReportToops() string {
	return beego.AppConfig.String("myslowreport.toops")
}

func myslowReportTodev() string {
	return beego.AppConfig.String("myslowreport.todev")
}

func myslowReportCcleader() string {
	return beego.AppConfig.String("myslowreport.ccleader")
}

func myslowReportTotest() string {
	return beego.AppConfig.String("myslowreport.totest")
}

func MyslowReportTos() string {
	tos := ""
	todev := strings.TrimRight(strings.TrimLeft(myslowReportTodev(), ";"), ";")
	if len(todev) > 1 {
		tos += todev
		tos += ";"
	}

	totest := strings.TrimRight(strings.TrimLeft(myslowReportTotest(), ";"), ";")
	if len(totest) > 1 {
		tos += totest
		tos += ";"
	}

	toops := strings.TrimRight(strings.TrimLeft(myslowReportToops(), ";"), ";")
	if len(toops) > 1 {
		tos += toops
	}

	tos = strings.TrimLeft(tos, ";")
	tos = strings.TrimRight(tos, ";")
	return fmt.Sprintf("%s", tos)
}

func MyslowReportCcs() string {
	ccs := ""
	ccleader := strings.TrimRight(strings.TrimLeft(myslowReportCcleader(), ";"), ";")
	if len(ccleader) > 1 {
		ccs += ccleader
	}

	ccs = strings.TrimLeft(ccs, ";")
	ccs = strings.TrimRight(ccs, ";")
	return fmt.Sprintf("%s", ccs)
}

func myslowReportCols() string {
	return beego.AppConfig.String("myslowreport.cols")
}

func MyslowReportColsList() []string {
	return strings.Split(myslowReportCols(), "|")
}

func myslowReportProjects() string {
	return beego.AppConfig.String("myslowreport.projects")
}

func MyslowReportProjectsList() []Project {
	var pl []Project
	l := strings.Split(myslowReportProjects(), ";")
	for _, p := range l {
		if len(p) < 2 {
			continue
		}
		s := strings.Split(p, "|")
		if len(s) < 6 {
			continue
		}
		project := Project {
			MysqlHost: s[0],
			MysqlPort: s[1],
			SlowlogTable: s[2],
			Description: s[3],
			RsyncModel: s[4],
			MysqlSlowlogFileName: s[5],
		}

		pl = append(pl, project)
	}

	return pl
}

func MyslowReportSlowlogpath() string {
	return beego.AppConfig.String("myslowreport.slowlogpath")
}

func MyslowReportSlowlogMonthltPath() string {
	return beego.AppConfig.String("myslowreport.slowlogmonthlypath")
}

func MyslowReportPtQueryDigest() string {
	return beego.AppConfig.String("myslowreport.ptquerydigest")
}

func MyslowReportTraceSheetInfoList() []TraceSheetInfo {
	l := []TraceSheetInfo{}
	s := beego.AppConfig.String("myslowreport.tracesheets")
	sl1 := strings.Split(s, "|")
	for _, s1 := range sl1 {
		items := strings.Split(s1, ";")
		// the first one is sheet name
		tsi := TraceSheetInfo{}
		myl := []MyInsChargeInfo{}

		tsi.SheetName = items[0]
		for i, v := range items {
			if i == 0 {
				continue
			}

			mylItems := strings.Split(v, ",")
			if len(mylItems) != 3 {
				continue
			}
			myi := MyInsChargeInfo{
				MyInsName: mylItems[0],
				Desc: mylItems[1],
				Charge: strings.Replace(mylItems[2], "#", " ", -1),
			}
			myl = append(myl, myi)
		}
		tsi.MyInsChargeInfoList = myl
		l = append(l, tsi)
	}

	return l
}

func MyslowReportTraceCols() []TraceColInfo {
	l := []TraceColInfo{}
	s := beego.AppConfig.String("myslowreport.tracecols")
	tl := strings.Split(s, "|")
	intForASCIIA := 65
	for i, v := range tl {
		item := TraceColInfo{
			ColName: v,
			ColCellName: fmt.Sprintf("%s1", string(i + intForASCIIA)),
			ColCellAlapha: fmt.Sprintf("%s", string(i + intForASCIIA)),
		}
		l = append(l, item)
	}
	return l
}