package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"strings"
	"fmt"
)

type Project struct {
	MysqlHost string
	MysqlPort string
	SlowlogTable string
	Description string
	RsyncModel string
	MysqlSlowlogFileName string
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

func myslowReportToleader() string {
	return beego.AppConfig.String("myslowreport.toleader")
}

func myslowReportTotest() string {
	return beego.AppConfig.String("myslowreport.totest")
}

func MyslowReportTos() string {
	todev := strings.TrimRight(strings.TrimLeft(myslowReportTodev(), ";"), ";")
	totest := strings.TrimRight(strings.TrimLeft(myslowReportTotest(), ";"), ";")
	toops := strings.TrimRight(strings.TrimLeft(myslowReportToops(), ";"), ";")
	toleader := strings.TrimRight(strings.TrimLeft(myslowReportToleader(), ";"), ";")

	return fmt.Sprintf("%s;%s;%s;%s", todev, totest, toops, toleader)
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

func MyslowReportPtQueryDigest() string {
	return beego.AppConfig.String("myslowreport.ptquerydigest")
}
