package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)

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

func MyslowReportTimeout() int {
	t, err := beego.AppConfig.Int("myslowreport.timeout")
	if err != nil {
		t = 10
	}
	return t
}

func SyncMysqlSlowlogSpec() string {
	return beego.AppConfig.String("myslowreport.syncmysqlslowlog")
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
