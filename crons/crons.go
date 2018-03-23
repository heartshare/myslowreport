package crons

import (
	_ "github.com/astaxie/beego"
	"github.com/robfig/cron"
	"github.com/ximply/myslowreport/models"
	"github.com/ximply/myslowreport/rsync"
	"fmt"
	"github.com/ximply/myslowreport/perconatoolkit"
	"github.com/ximply/myslowreport/utils"
	"github.com/astaxie/beego"
	"strings"
	"github.com/ximply/myslowreport/email"
	"github.com/shopspring/decimal"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ximply/myslowreport/dingding"
)

func monthlyTableName(table string) string {
	return fmt.Sprintf("%s_monthly", table)
}

func createReport() string {
	pl := models.MyslowReportProjectsList()
	slowInfo := ""
	slowInfo += `<html>`
	slowInfo += htmlHead()
	slowInfo += `<body>`

	slowInfo += `<div id="content" class="reg">`
	slowInfo += createTitle()
	slowInfo += `<div style="LINE-HEIGHT: 20px; FONT-FAMILY: 'Microsoft Yahei'; FONT-SIZE: 12px">`
	slowInfo += `<div style="MARGIN: 0px auto; WIDTH: 1280px; BACKGROUND: #FF710C; border-top-right-radius: 8px; border-top-left-radius: 8px; border-bottom-right-radius: 8px; border-bottom-left-radius: 8px; -moz-border-radius: 8px; -webkit-border-radius: 8px">`
	slowInfo += `<div style="WIDTH: 710px; HEIGHT: 10px"></div>`
	slowInfo += `<div style="PADDING-BOTTOM: 0px; PADDING-LEFT: 0px; WIDTH: 710px; PADDING-RIGHT: 0px; BACKGROUND: #FF710C; HEIGHT: 30px; PADDING-TOP: 10px">`
	slowInfo += createSubject()
	slowInfo += `</div>`

	slowInfo += `<div style="BORDER-BOTTOM: #FF710C 1px solid; BORDER-LEFT: #FF710C 1px solid; PADDING-BOTTOM: 5px; PADDING-LEFT: 9px; PADDING-RIGHT: 9px; BACKGROUND: #f9f9f9; BORDER-TOP: #FF710C 1px solid; BORDER-RIGHT: #FF710C 1px solid; PADDING-TOP: 5px">`
	slowInfo += `<div style="PADDING-BOTTOM: 0px; MARGIN: 10px 0px 0px; PADDING-LEFT: 15px; PADDING-RIGHT: 0px; PADDING-TOP: 0px">`

	since := fmt.Sprintf("%s 23:59:59", utils.BeforeYesterdayStringByFormat("2006-01-02"))
	until := fmt.Sprintf("%s 00:00:00", utils.TodayStringByFormat("2006-01-02"))
	for i, p := range pl {
		beego.Info(fmt.Sprintf("Start GetOrderByQueryTimeMaxDesc: %s", p.SlowlogTable))
		cnt, items := models.GetOrderByQueryTimeMaxDesc(since, until, p.SlowlogTable)
		beego.Info(fmt.Sprintf("End GetOrderByQueryTimeMaxDesc: %s, Count: %d", p.SlowlogTable, cnt))

		if cnt == 0 {
			beego.Info(fmt.Sprintf("No slow logs: %s", p.SlowlogTable))
			continue
		}

		cnt, mi := models.GetMaxOrderBy(since, until, p.SlowlogTable)
		if cnt == 0 {
			beego.Info(fmt.Sprintf("No max values: %s", p.SlowlogTable))
		}

		beego.Info(fmt.Sprintf("Creating report content: %s", p.SlowlogTable))
		slowInfo += createSlowInfo(items, p, i, mi)
		slowInfo += "<br/>"
		slowInfo += `<hr>`
		slowInfo += "<br/>"
		beego.Info(fmt.Sprintf("End Create report content: %s", p.SlowlogTable))
	}

	slowInfo += `</div>`
	slowInfo += `</div>`
	slowInfo += `</div>`
	slowInfo += `</div>`
	slowInfo += `</div>`
	slowInfo += `</body>`
	slowInfo += `</html>`

	return slowInfo
}

func createMonthlyReport() string {
	pl := models.MyslowReportProjectsList()
	slowInfo := ""
	slowInfo += `<html>`
	slowInfo += htmlHead()
	slowInfo += `<body>`

	slowInfo += `<div id="content" class="reg">`
	slowInfo += createTitle()
	slowInfo += `<div style="LINE-HEIGHT: 20px; FONT-FAMILY: 'Microsoft Yahei'; FONT-SIZE: 12px">`
	slowInfo += `<div style="MARGIN: 0px auto; WIDTH: 1280px; BACKGROUND: #FF710C; border-top-right-radius: 8px; border-top-left-radius: 8px; border-bottom-right-radius: 8px; border-bottom-left-radius: 8px; -moz-border-radius: 8px; -webkit-border-radius: 8px">`
	slowInfo += `<div style="WIDTH: 710px; HEIGHT: 10px"></div>`
	slowInfo += `<div style="PADDING-BOTTOM: 0px; PADDING-LEFT: 0px; WIDTH: 710px; PADDING-RIGHT: 0px; BACKGROUND: #FF710C; HEIGHT: 30px; PADDING-TOP: 10px">`
	slowInfo += createMonthlySubject()
	slowInfo += `</div>`

	slowInfo += `<div style="BORDER-BOTTOM: #FF710C 1px solid; BORDER-LEFT: #FF710C 1px solid; PADDING-BOTTOM: 5px; PADDING-LEFT: 9px; PADDING-RIGHT: 9px; BACKGROUND: #f9f9f9; BORDER-TOP: #FF710C 1px solid; BORDER-RIGHT: #FF710C 1px solid; PADDING-TOP: 5px">`
	slowInfo += `<div style="PADDING-BOTTOM: 0px; MARGIN: 10px 0px 0px; PADDING-LEFT: 15px; PADDING-RIGHT: 0px; PADDING-TOP: 0px">`

	since := fmt.Sprintf("%s 23:59:59",
		utils.DateString(utils.BeforeYesterday().AddDate(0, -1, 0).Format("2006-01-02")))
	until := fmt.Sprintf("%s 00:00:00", utils.TodayStringByFormat("2006-01-02"))
	for i, p := range pl {
		monthlyTable := monthlyTableName(p.SlowlogTable)
		beego.Info(fmt.Sprintf("Start GetOrderByQueryTimeMaxDesc: %s", monthlyTable))
		cnt, items := models.GetOrderByQueryTimeMaxDesc(since, until, monthlyTable)
		beego.Info(fmt.Sprintf("End GetOrderByQueryTimeMaxDesc: %s, Count: %d", monthlyTable, cnt))

		if cnt == 0 {
			beego.Info(fmt.Sprintf("No slow logs: %s", monthlyTable))
			continue
		}

		cnt, mi := models.GetMaxOrderBy(since, until, monthlyTable)
		if cnt == 0 {
			beego.Info(fmt.Sprintf("No max values: %s", monthlyTable))
		}

		beego.Info(fmt.Sprintf("Creating report content: %s", monthlyTable))
		slowInfo += createMonthlySlowInfo(items, p, i, mi)
		slowInfo += "<br/>"
		slowInfo += `<hr>`
		slowInfo += "<br/>"
		beego.Info(fmt.Sprintf("End Create report content: %s", monthlyTable))
	}

	slowInfo += `</div>`
	slowInfo += `</div>`
	slowInfo += `</div>`
	slowInfo += `</div>`
	slowInfo += `</div>`
	slowInfo += `</body>`
	slowInfo += `</html>`

	return slowInfo
}

func htmlHead() string {
	head :=
		`<head>
    		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    		<style class="fox_global_style">

        		div.fox_html_content {
            		line-height: 1.5;
        		}

        		blockquote {
            		margin-top: 0px;
            		margin-bottom: 0px;
            		margin-left: 0.5em;
        		}

        		ol,
        		ul {
            		margin-top: 0px;
            		margin-bottom: 0px;
            		list-style-position: inside;
        		}

        		p {
            		margin-top: 0px;
            		margin-bottom: 0px
        		}

				hr {
					color: lightgray;
				}

				#hostP, #totalP, #uniqP {
					margin-top: 0px;
					margin-bottom: 0px;
					margin-left: 7px;
					font-size: 16px;
					text-decoration: none;
				}

				#hostSpan {
					float: left;
					color: #262626;
					font-size: 18px;
					font-weight: bold;
				}

				#hostDescSpan, #totalSpan, #uniqSpan {
					color: black;
					font-size: 15px;
				}

				#tableDispDiv {
					padding-bottom: 0px;
					padding-left: 0px;
					padding-right: 0px;
					padding-top: 0px;
				}

				#tableHeader {
					background-color: #FF710C;
					display: table-row;
					vertical-align: inherit;
				}

				#tableTr {
					display: table-row;
					vertical-align: inherit;
				}

				#tableHeader300 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: center;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					min-height: 50px;
					max-width: 300px;
					padding-right: 15px;
					height: 30px;
					color: rgb(255,255,255);
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px
				}

				#tableCell300 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: left;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					max-width: 300px;
					padding-right: 15px;
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px;
				}

				#tableHeader50 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: center;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					min-height: 50px;
					max-width: 50px;
					padding-right: 15px;
					height: 30px;
					color: rgb(255,255,255);
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px
				}

				#tableCell50 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: left;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					max-width: 50px;
					padding-right: 15px;
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px;
				}

				#tableHeader90 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: center;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					min-height: 50px;
					max-width: 90px;
					padding-right: 15px;
					height: 30px;
					color: rgb(255,255,255);
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px
				}

				#tableCell90 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: left;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					max-width: 90px;
					padding-right: 15px;
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px;
				}

				#tableHeader60 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: center;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					min-height: 50px;
					max-width: 60px;
					padding-right: 15px;
					height: 30px;
					color: rgb(255,255,255);
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px
				}

				#tableCell60 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: left;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					max-width: 60px;
					padding-right: 15px;
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px;
				}

				#tableHeader80 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: center;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					min-height: 50px;
					max-width: 80px;
					padding-right: 15px;
					height: 30px;
					color: rgb(255,255,255);
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px
				}

				#tableCell80 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: left;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					max-width: 80px;
					padding-right: 15px;
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px;
				}

				#tableHeader160 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: center;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					min-height: 50px;
					max-width: 160px;
					padding-right: 15px;
					height: 30px;
					color: rgb(255,255,255);
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px
				}

				#tableCell160 {
					border-bottom: rgb(222,222,222) 1px solid;
					text-align: left;
					border-left: rgb(222,222,222) 1px solid;
					padding-bottom: 7px;
					margin: 0px;
					padding-left: 7px;
					max-width: 160px;
					padding-right: 15px;
					font-size: 13px;
					border-right: rgb(241,241,226) 1px solid;
					padding-top: 7px;
				}
    		</style>
			<script type="text/javascript">
				function toggle(id) {
 					var tb=document.getElementById(id);
 					if (tb.style.display=='none') tb.style.display='block';
 					else tb.style.display='none';
				}

				function getTooltip(_obj, tip) {
   					var tValue = tip + _obj.innerText;
   					_obj.setAttribute("title", tValue);
				}
			</script>
		</head>`

		return head
}

func growthRate(this int64, other int64) float64 {
	if other != 0 {
		return (float64(this - other)) * 100 / float64(other)
	}
	return 0.0
}

func growtRateStrAndColor(rate float64, other int64) (string, string) {
	str := ""
	color := "black"
	if rate > 0.0 {
		str = fmt.Sprintf("+%.2f%%", rate)
		color = "red"
	} else if rate < 0.0 {
		str = fmt.Sprintf("%.2f%%", rate)
		color = "green"
	} else {
		if other == 0{
			str = "NULL"
			color = "gray"
		} else {
			str = "--"
		}
	}

	return str, color
}

func yesterdayTotal(table string) int64 {
	return models.GetSumOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.BeforeYesterdayStringByFormat("2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.TodayStringByFormat("2006-01-02")),
		table)
}

func lastMonthTotal(table string) int64 {
	since := fmt.Sprintf("%s 23:59:59",
		utils.DateString(utils.BeforeYesterday().AddDate(0, -1, 0).Format("2006-01-02")))
	until := fmt.Sprintf("%s 00:00:00", utils.TodayStringByFormat("2006-01-02"))

	return models.GetSumOfQueryCount(since, until, table)
}

func lastLastMonthTotal(table string) int64 {
	since := fmt.Sprintf("%s 23:59:59",
		utils.DateString(utils.BeforeYesterday().AddDate(0, -2, 0).Format("2006-01-02")))
	until := fmt.Sprintf("%s 00:00:00", utils.Today().AddDate(0, -1, 0).Format("2006-01-02"))

	return models.GetSumOfQueryCount(since, until, table)
}

func beforeYesterdayTotal(table string) int64 {
	return models.GetSumOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.BeforeBeforeYesterdayStringByFormat("2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.YesterdayStringByFormat("2006-01-02")),
		table)
}

func yoyBasisTotal(table string) int64 {
	return models.GetSumOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.DateStringByFormat(-9,"2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.DateStringByFormat(-7,"2006-01-02")),
		table)
}

func yesterdayUniq(table string) int64 {
	return models.GetUniqOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.BeforeYesterdayStringByFormat("2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.TodayStringByFormat("2006-01-02")),
		table)
}

func lastMonthUniq(table string) int64 {
	since := fmt.Sprintf("%s 23:59:59",
		utils.DateString(utils.BeforeYesterday().AddDate(0, -1, 0).Format("2006-01-02")))
	until := fmt.Sprintf("%s 00:00:00", utils.TodayStringByFormat("2006-01-02"))

	return models.GetUniqOfQueryCount(since, until, table)
}

func lastLastMonthUniq(table string) int64 {
	since := fmt.Sprintf("%s 23:59:59",
		utils.DateString(utils.BeforeYesterday().AddDate(0, -2, 0).Format("2006-01-02")))
	until := fmt.Sprintf("%s 00:00:00", utils.Today().AddDate(0, -1, 0).Format("2006-01-02"))

	return models.GetUniqOfQueryCount(since, until, table)
}

func beforeYesterdayUniq(table string) int64 {
	return models.GetUniqOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.BeforeBeforeYesterdayStringByFormat("2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.YesterdayStringByFormat("2006-01-02")),
		table)
}

func yoyBasisUniq(table string) int64 {
	return models.GetUniqOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.DateStringByFormat(-9,"2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.DateStringByFormat(-7,"2006-01-02")),
		table)
}

func toMyInsName(ip string, port string) string {
	return fmt.Sprintf("%s_%s",
		strings.Replace(ip, ".", "_", -1), port)
}

func createSlowInfo(items []models.Item, p models.Project, tableId int, mi models.MaxItem) string {
	var info = ""

	info += `<div style="PADDING-BOTTOM: 0px; MARGIN: 10px 0px 0px; PADDING-LEFT: 15px; PADDING-RIGHT: 0px; PADDING-TOP: 0px">`
	info += `<p id="hostP"><span id="hostSpan">{{.Host}}</span><span id="hostDescSpan">&nbsp;&nbsp;{{.Project}}</span></p>`
	info = strings.Replace(info, "{{.Host}}", p.MysqlHost, -1)
	info = strings.Replace(info, "{{.Project}}", p.Description, -1)

	// chain grow with before yesterday and on year-on-year basis (this day last week)
	yt := yesterdayTotal(p.SlowlogTable)
	byt := beforeYesterdayTotal(p.SlowlogTable)
	yoybt := yoyBasisTotal(p.SlowlogTable)

	chainGrowthRateForTotal := growthRate(yt, byt)
	chainGrowthRateForTotalStr, chainGrowthRateForTotalStrColor := growtRateStrAndColor(chainGrowthRateForTotal, byt)
	yoyBasisRateForTotal := growthRate(yt, yoybt)
	yoyBasisRateForTotalStr, yoyBasisRateForTotalStrColor := growtRateStrAndColor(yoyBasisRateForTotal, yoybt)

	yu := yesterdayUniq(p.SlowlogTable)
	byu := beforeYesterdayUniq(p.SlowlogTable)
	yoybu := yoyBasisUniq(p.SlowlogTable)

	chainGrowthRateForUniq := growthRate(yu, byu)
	chainGrowthRateForUniqStr, chainGrowthRateForUniqStrColor := growtRateStrAndColor(chainGrowthRateForUniq, byu)
	yoyBasisRateForUniq := growthRate(yu, yoybu)
	yoyBasisRateForUniqStr, yoyBasisRateForUniqStrColor := growtRateStrAndColor(yoyBasisRateForUniq, yoybu)

	myInsName := toMyInsName(p.MysqlHost, p.MysqlPort)
	statDate := fmt.Sprintf("%s 00:00:00", utils.DateString(utils.YesterdayStringByFormat("2006-01-02")))
	count, _ := models.GetByMyInsNameAndStatDate(myInsName, statDate)

	if count == 0 {
		gr := models.GrowRate{
			MyInsName: myInsName,
			StatDate: utils.Yesterday(),

			YesterdayTotal: yt,
			BeforeYesterdayTotal: byt,
			BasisTotal: yoybt,

			YesterdayUniq: yu,
			BeforeYesterdayUniq: byu,
			BasisUniq: yoybu,

			TotalChainRate: decimal.NewFromFloat(chainGrowthRateForTotal),
			TotalBasisRate: decimal.NewFromFloat(yoyBasisRateForTotal),

			UniqChainRate: decimal.NewFromFloat(chainGrowthRateForUniq),
			UniqBasisRate: decimal.NewFromFloat(yoyBasisRateForUniq),
		}
		models.AddGrowRate(gr)
	}

	info += `<p id="totalP"><span id="totalSpan">&nbsp;&nbsp; {{.TotalCount}} &nbsp;&nbsp;</span><span style="color: {{.ChainColor}}; font-size: 15px">&nbsp;&nbsp; {{.TChain}} &nbsp;&nbsp;</span></span><span style="color: {{.YOYColor}}; font-size: 15px">&nbsp;&nbsp; {{.YOYTotal}} &nbsp;&nbsp;</span></p>`
	info = strings.Replace(info, "{{.TotalCount}}", fmt.Sprintf("总语句数: %d", yt), -1)
	info = strings.Replace(info, "{{.TChain}}", fmt.Sprintf("环比前天: %s", chainGrowthRateForTotalStr), -1)
	info = strings.Replace(info, "{{.ChainColor}}", fmt.Sprintf("%s", chainGrowthRateForTotalStrColor), -1)
	info = strings.Replace(info, "{{.YOYTotal}}", fmt.Sprintf("同比上周%s: %s", utils.WeekdayCNShortString(utils.Yesterday()), yoyBasisRateForTotalStr), -1)
	info = strings.Replace(info, "{{.YOYColor}}", fmt.Sprintf("%s", yoyBasisRateForTotalStrColor), -1)
	
	info += `<p id="uniqP"><span id="uniqSpan">&nbsp;&nbsp; {{.UniqCount}} &nbsp;&nbsp;</span><span style="color: {{.ChainColor}}; font-size: 15px">&nbsp;&nbsp; {{.UChain}} &nbsp;&nbsp;</span></span><span style="color: {{.YOYColor}}; font-size: 15px">&nbsp;&nbsp; {{.YOYUniq}} &nbsp;&nbsp;</span></p>`
	info = strings.Replace(info, "{{.UniqCount}}", fmt.Sprintf("独立语句数: %d", yu), -1)
	info = strings.Replace(info, "{{.UChain}}", fmt.Sprintf("环比前天: %s", chainGrowthRateForUniqStr), -1)
	info = strings.Replace(info, "{{.ChainColor}}", fmt.Sprintf("%s", chainGrowthRateForUniqStrColor), -1)
	info = strings.Replace(info, "{{.YOYUniq}}", fmt.Sprintf("同比上周%s: %s", utils.WeekdayCNShortString(utils.Yesterday()), yoyBasisRateForUniqStr), -1)
	info = strings.Replace(info, "{{.YOYColor}}", fmt.Sprintf("%s", yoyBasisRateForUniqStrColor), -1)
	
	info += `<div style="MARGIN: 20px 0px 0px 7px">`
	info += `<div id="tableDispDiv">`
	info += `<input type="button" value="显示详情/隐藏详情" onClick="toggle('WhatTable')"/>`
	info = strings.Replace(info, "WhatTable", fmt.Sprintf("table%d", tableId), -1)
	info += `<table TableId style="display: none;TEXT-ALIGN: left; MARGIN: 2px 0px 0px; WIDTH: 100%; BORDER-COLLAPSE: collapse" rules="none" bordercolor="#c7c7c7" frame="box">`
	info = strings.Replace(info, "TableId", fmt.Sprintf("id = \"table%d\"", tableId), -1)
	info += `<tbody>`

	info += `<tr id="tableHeader">`
	// col name,col width|col name,col width|col name,col width ...
	cols := models.MyslowReportColsList()
	if len(cols) == 0 {
		return ""
	}

	tmp := `<td id="tableHeader300">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[0], -1)
	info += tmp

	tmp = `<td id="tableHeader50">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[1], -1)
	info += tmp

	tmp = `<td id="tableHeader90">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[2], -1)
	info += tmp

	tmp = `<td id="tableHeader60">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[3], -1)
	info += tmp

	tmp = `<td id="tableHeader80">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[4], -1)
	info += tmp

	tmp = `<td id="tableHeader80">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[5], -1)
	info += tmp

	tmp = `<td id="tableHeader80">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[6], -1)
	info += tmp

	tmp = `<td id="tableHeader80">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[7], -1)
	info += tmp

	tmp = `<td id="tableHeader80">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[8], -1)
	info += tmp

	tmp = `<td id="tableHeader160">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[9], -1)
	info += tmp

	info += `</tr>`

	for _, i := range items {
		info += `<tr id="tableTr">`
		info += createItem(i, cols)
		info += `</tr>`
	}

	info += `</tbody>`
	info += `</table>`
	info += `</div>`
	info += `</div>`
	info += `</div>`

	return info
}

func createMonthlySlowInfo(items []models.Item, p models.Project, tableId int, mi models.MaxItem) string {
	var info = ""

	info += `<div style="PADDING-BOTTOM: 0px; MARGIN: 10px 0px 0px; PADDING-LEFT: 15px; PADDING-RIGHT: 0px; PADDING-TOP: 0px">`
	info += `<p style="MARGIN-TOP: 0px; MARGIN-BOTTOM: 0px; MARGIN-LEFT: 7px; FONT-SIZE: 16px; TEXT-DECORATION: none"><span style="FLOAT: left; COLOR: #262626; FONT-SIZE: 18px">&nbsp;&nbsp; </span> <span style="FLOAT: left; COLOR: #262626; FONT-SIZE: 18px; font-weight:bold">192.168.10.1</span><span style="COLOR: rgb(0,0,0); FONT-SIZE: 15px">&nbsp;&nbsp; Project</span></p>`
	info = strings.Replace(info, "192.168.10.1", p.MysqlHost, -1)
	info = strings.Replace(info, "Project", p.Description, -1)

	// chain grow with before yesterday and on year-on-year basis (this day last week)
	monthlyTable := monthlyTableName(p.SlowlogTable)
	lmt := lastMonthTotal(monthlyTable)
	l2mt := lastLastMonthTotal(monthlyTable)
	chainRateForTotal := growthRate(lmt, l2mt)
	chainRateForTotalStr, chainRateForTotalStrColor := growtRateStrAndColor(chainRateForTotal, l2mt)

	lmu := lastMonthUniq(monthlyTable)
	l2mu := lastLastMonthUniq(monthlyTable)
	chainRateForUniq := growthRate(lmu, l2mu)
	chainRateForUniqStr, chainRateForUniqStrColor := growtRateStrAndColor(chainRateForUniq, l2mu)

	myInsName := toMyInsName(p.MysqlHost, p.MysqlPort)
	statDate := fmt.Sprintf("%s-01 00:00:00", utils.YearMonthStringByFormat(utils.Today(),"2006-01-02"))
	count, _ := models.GetByMyInsNameAndStatDateMonthly(myInsName, statDate)

	if count == 0 {
		grm := models.GrowRateMonthly{
			MyInsName: myInsName,
			StatDate: utils.Today(),

			LastMonthTotal: lmt,
			LastLastMonthTotal: l2mt,

			LastMonthUniq: lmu,
			LastLastMonthUniq: l2mu,

			TotalChainRate: decimal.NewFromFloat(chainRateForTotal),
			UniqChainRate: decimal.NewFromFloat(chainRateForUniq),
		}
		models.AddGrowRateMonthly(grm)
	}

	info += `<p style="MARGIN-TOP: 0px; MARGIN-BOTTOM: 0px; MARGIN-LEFT: 7px; FONT-SIZE: 16px; TEXT-DECORATION: none"><span style="COLOR: rgb(0,0,0); FONT-SIZE: 15px">&nbsp;&nbsp; TotalCount &nbsp;&nbsp;</span><span style="COLOR: ChainColor; FONT-SIZE: 15px">&nbsp;&nbsp; TChain &nbsp;&nbsp;</span></span></p>`
	info = strings.Replace(info, "TotalCount", fmt.Sprintf("总语句数: %d", lmt), -1)
	info = strings.Replace(info, "TChain", fmt.Sprintf("环比%s: %s",
		utils.YearMonthStringByFormat(utils.BeforeYesterday().AddDate(0, -2, 0),"2006-01-02"),
		chainRateForTotalStr), -1)
	info = strings.Replace(info, "ChainColor", fmt.Sprintf("%s", chainRateForTotalStrColor), -1)

	info += `<p style="MARGIN-TOP: 0px; MARGIN-BOTTOM: 0px; MARGIN-LEFT: 7px; FONT-SIZE: 16px; TEXT-DECORATION: none"><span style="COLOR: rgb(0,0,0); FONT-SIZE: 15px">&nbsp;&nbsp; UniqCount &nbsp;&nbsp;</span><span style="COLOR: ChainColor; FONT-SIZE: 15px">&nbsp;&nbsp; UChain &nbsp;&nbsp;</span></span></p>`
	info = strings.Replace(info, "UniqCount", fmt.Sprintf("独立语句数: %d", lmu), -1)
	info = strings.Replace(info, "UChain", fmt.Sprintf("环比%s: %s",
		utils.YearMonthStringByFormat(utils.BeforeYesterday().AddDate(0, -2, 0),"2006-01-02"),
		chainRateForUniqStr), -1)
	info = strings.Replace(info, "ChainColor", fmt.Sprintf("%s", chainRateForUniqStrColor), -1)

	info += `<div style="MARGIN: 20px 0px 0px 7px">`
	info += `<div style="PADDING-BOTTOM: 0px; PADDING-LEFT: 15px; PADDING-RIGHT: 0px; PADDING-TOP: 0px">`
	info += `<input type="button" value="显示详情/隐藏详情" onClick="toggle('WhatTable')"/>`
	info = strings.Replace(info, "WhatTable", fmt.Sprintf("table%d", tableId), -1)
	info += `<table TableId style="display: none;TEXT-ALIGN: left; MARGIN: 2px 0px 0px; WIDTH: 100%; BORDER-COLLAPSE: collapse" rules="none" bordercolor="#c7c7c7" frame="box">`
	info = strings.Replace(info, "TableId", fmt.Sprintf("id = \"table%d\"", tableId), -1)
	info += `<tbody>`

	info += `<tr style="BACKGROUND-COLOR: #FF710C; DISPLAY: table-row; VERTICAL-ALIGN: inherit">`
	// col name,col width|col name,col width|col name,col width ...
	cols := models.MyslowReportColsList()
	if len(cols) == 0 {
		return ""
	}

	tmp := `<td id="tableHeader300">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[0], -1)
	info += tmp

	tmp = `<td id="tableHeader50">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[1], -1)
	info += tmp

	tmp = `<td id="tableHeader90">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[2], -1)
	info += tmp

	tmp = `<td id="tableHeader60">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[3], -1)
	info += tmp

	tmp = `<td id="tableHeader80">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[4], -1)
	info += tmp

	tmp = `<td id="tableHeader80">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[5], -1)
	info += tmp

	tmp = `<td id="tableHeader80">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[6], -1)
	info += tmp

	tmp = `<td id="tableHeader80">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[7], -1)
	info += tmp

	tmp = `<td id="tableHeader80">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[8], -1)
	info += tmp

	tmp = `<td id="tableHeader160">{{.HeaderColName}}</td>`
	tmp = strings.Replace(tmp, "{{.HeaderColName}}", cols[9], -1)
	info += tmp

	info += `</tr>`

	for _, i := range items {
		info += `<tr id="tableTr">`
		info += createItem(i, cols)
		info += `</tr>`
	}

	info += `</tbody>`
	info += `</table>`
	info += `</div>`
	info += `</div>`
	info += `</div>`

	return info
}

func createTitle() string {
	return fmt.Sprintf(`<title>%s</title>`, models.MyslowReportTitle())
}

func createSubject() string {
	return fmt.Sprintf(`<p style="MARGIN: 0px; PADDING-LEFT: 30px; FONT-FAMILY: 'Microsoft Yahei'; COLOR: rgb(255,255,255); FONT-SIZE: 24px">%s %s %s</p>`,
		models.MyslowReportSubject(), utils.YesterdayStringByFormat("2006-01-02"), utils.WeekdayCNString(utils.Yesterday()))
}

func createMonthlySubject() string {
	return fmt.Sprintf(`<p style="MARGIN: 0px; PADDING-LEFT: 30px; FONT-FAMILY: 'Microsoft Yahei'; COLOR: rgb(255,255,255); FONT-SIZE: 24px">%s %s&nbsp;&nbsp;月报</p>`,
		models.MyslowReportSubject(), utils.YearMonthStringByFormat(utils.Yesterday(), "2006-01-02"))
}

func floatColString3(val float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.3f", val), "0"), ".")
}

func floatColString9(val float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.9f", val), "0"), ".")
}

func createItem(item models.Item, colsName []string) string {
	s := ""
	i := 0

	//dfMax := `<td onmouseover='getTooltip(this,"ThisColName: ")' style="box-shadow: inset 0px 0px 25px red;BORDER-BOTTOM: rgb(222,222,222) 1px solid; TEXT-ALIGN: left; BORDER-LEFT: rgb(222,222,222) 1px solid; PADDING-BOTTOM: 7px; MARGIN: 0px; PADDING-LEFT: 7px; MAX-WIDTH: 999999999px; PADDING-RIGHT: 15px; FONT-SIZE: 13px; BORDER-RIGHT: rgb(241,241,226) 1px solid; PADDING-TOP: 7px;">xxx</td>`
	//dfMax := `<td onmouseover='getTooltip(this,"ThisColName: ")' style="BORDER-BOTTOM: rgb(222,222,222) 1px solid; TEXT-ALIGN: left; BORDER-LEFT: rgb(222,222,222) 1px solid; PADDING-BOTTOM: 7px; MARGIN: 0px; PADDING-LEFT: 7px; MAX-WIDTH: 999999999px; PADDING-RIGHT: 15px; FONT-SIZE: 13px; BORDER-RIGHT: rgb(241,241,226) 1px solid; PADDING-TOP: 7px;">xxx</td>`


	tmp := `<td id="tableCell300"><textarea readonly="readonly" style="max-width:300px; min-height: 150px; max-height: 180px;">{{.xxx}}</textarea></td>`
	tmp = strings.Replace(tmp, "{{.xxx}}", item.Sample, -1)
	s += tmp
	i++

	tsCntStr := fmt.Sprintf("%d", item.TsCnt)
	tmp = `<td onmouseover='getTooltip(this,"{{.ThisColName}}: ")' id="tableCell50">{{.xxx}}</td>`
	tmp = strings.Replace(tmp, "{{.ThisColName}}", colsName[i], -1)
	tmp = strings.Replace(tmp, "{{.xxx}}", tsCntStr, -1)
	s += tmp
	i++

	tmp = `<td onmouseover='getTooltip(this,"{{.ThisColName}}: ")' id="tableCell90">{{.xxx}}</td>`
	tmp = strings.Replace(tmp, "{{.ThisColName}}", colsName[i], -1)
	tmp = strings.Replace(tmp, "{{.xxx}}", item.UserMax, -1)
	s += tmp
	i++

	f1, _ := item.QueryTimeMax.Float64()
	queryTimeMaxStr := floatColString3(f1)
	tmp = `<td onmouseover='getTooltip(this,"{{.ThisColName}}: ")' id="tableCell60">{{.xxx}}</td>`
	tmp = strings.Replace(tmp, "{{.ThisColName}}", colsName[i], -1)
	tmp = strings.Replace(tmp, "{{.xxx}}", queryTimeMaxStr, -1)
	s += tmp
	i++

	f2, _ := item.QueryTimePct95.Float64()
	tmp = `<td onmouseover='getTooltip(this,"{{.ThisColName}}: ")' id="tableCell80">{{.xxx}}</td>`
	tmp = strings.Replace(tmp, "{{.ThisColName}}", colsName[i], -1)
	tmp = strings.Replace(tmp, "{{.xxx}}", floatColString3(f2), -1)
	s += tmp
	i++

	f3, _ := item.LockTimeMax.Float64()
	lockTimeMaxStr := floatColString9(f3)
	tmp = `<td onmouseover='getTooltip(this,"{{.ThisColName}}: ")' id="tableCell80">{{.xxx}}</td>`
	tmp = strings.Replace(tmp, "{{.ThisColName}}", colsName[i], -1)
	tmp = strings.Replace(tmp, "{{.xxx}}", lockTimeMaxStr, -1)
	s += tmp
	i++

	f4, _ := item.LockTimePct95.Float64()
	tmp = `<td onmouseover='getTooltip(this,"{{.ThisColName}}: ")' id="tableCell80">{{.xxx}}</td>`
	tmp = strings.Replace(tmp, "{{.ThisColName}}", colsName[i], -1)
	tmp = strings.Replace(tmp, "{{.xxx}}", floatColString9(f4), -1)
	s += tmp
	i++

	rowsExaminedMaxStr := fmt.Sprintf("%d", item.RowsExaminedMax)
	tmp = `<td onmouseover='getTooltip(this,"{{.ThisColName}}: ")' id="tableCell80">{{.xxx}}</td>`
	tmp = strings.Replace(tmp, "{{.ThisColName}}", colsName[i], -1)
	tmp = strings.Replace(tmp, "{{.xxx}}", rowsExaminedMaxStr, -1)
	s += tmp
	i++

	tmp = `<td onmouseover='getTooltip(this,"{{.ThisColName}}: ")' id="tableCell80">{{.xxx}}</td>`
	tmp = strings.Replace(tmp, "{{.ThisColName}}", colsName[i], -1)
	tmp = strings.Replace(tmp, "{{.xxx}}", fmt.Sprintf("%d", item.RowsExaminedPct95), -1)
	s += tmp
	i++

	tmp = `<td onmouseover='getTooltip(this,"{{.ThisColName}}: ")' id="tableCell160">{{.xxx}}</td>`
	tmp = strings.Replace(tmp, "{{.ThisColName}}", colsName[i], -1)
	tmp = strings.Replace(tmp, "{{.xxx}}", item.TsMax.Format("2006-01-02 15:04:05"), -1)
	s += tmp

	return s
}

func syncMysqlDailySlowlog() {
	pl := models.MyslowReportProjectsList()
	slowlogPath := models.MyslowReportSlowlogpath()
	slowlogMonthlyPath := models.MyslowReportSlowlogMonthltPath()
	for _, p := range pl {
		beego.Info(fmt.Sprintf("Start SyncMysqlSlowlogFile: %s, %s", p.MysqlHost, p.MysqlSlowlogFileName))
		file := rsync.SyncMysqlSlowlogFile(
			p.MysqlHost,
			p.MysqlPort,
			p.RsyncModel,
			p.MysqlSlowlogFileName,
			slowlogPath)
		beego.Info(fmt.Sprintf("End SyncMysqlSlowlogFile: %s, %s", p.MysqlHost, p.MysqlSlowlogFileName))

		err := rsync.MergeMysqlSlowlogFile(p.MysqlHost, p.MysqlPort, slowlogPath, slowlogMonthlyPath)
		if err != nil {
			beego.Info(fmt.Sprintf("MergeMysqlSlowlogFile fail: %s, %s", p.SlowlogTable, err.Error()))
		}

		if !perconatoolkit.ImportMysqlSlowlogHistoryToMysql(file, p.SlowlogTable) {
			beego.Info(fmt.Sprintf("ImportMysqlSlowlogHistoryToMysql fail: %s, %s", file, p.SlowlogTable))
			continue
		}
		beego.Info(fmt.Sprintf("ImportMysqlSlowlogHistoryToMysql successfully: %s, %s", file, p.SlowlogTable))
	}
}

func importMysqlMonthlySlowlog() {
	pl := models.MyslowReportProjectsList()
	slowlogMonthlyPath := models.MyslowReportSlowlogMonthltPath()
	for _, p := range pl {
		monthFile := fmt.Sprintf("%s%s.%s.%s", slowlogMonthlyPath, p.MysqlHost, p.MysqlPort,
			utils.YearMonthStringByFormat(utils.Yesterday(), "20060102"))
		monthlyTable := monthlyTableName(p.SlowlogTable)

		if !perconatoolkit.ImportMysqlSlowlogHistoryToMysql(monthFile, monthlyTable) {
			beego.Info(fmt.Sprintf("ImportMysqlSlowlogHistoryToMysql fail: %s, %s", monthFile, monthlyTable))
			continue
		}
		beego.Info(fmt.Sprintf("ImportMysqlSlowlogHistoryToMysql successfully: %s, %s", monthFile, monthlyTable))
	}
}

func myInsNameToMyslowMonthlyTableName(myInsName string) string {
	return strings.Replace(myInsName, "192_168", "myslow_history", -1) + "_monthly"
}

func createTraceExcelReport(file string) error {
	xlsx := excelize.NewFile()
	sheetList := models.MyslowReportTraceSheetInfoList()
	colList := models.MyslowReportTraceCols()

	for _, v := range sheetList {
		xlsx.NewSheet(v.SheetName)
		for _, vv := range colList {
			xlsx.SetCellValue(v.SheetName, vv.ColCellName, vv.ColName)
		}

		since := fmt.Sprintf("%s 23:59:59",
			utils.DateString(utils.BeforeYesterday().AddDate(0, -1, 0).Format("2006-01-02")))
		until := fmt.Sprintf("%s 00:00:00", utils.TodayStringByFormat("2006-01-02"))
		step := 2
		for _, vvv := range v.MyInsChargeInfoList {
			table := myInsNameToMyslowMonthlyTableName(vvv.MyInsName)
			_, items := models.GetOrderByQueryTimeMaxDesc(since, until, table)
			count := len(items)
			if count > 0 {
				for i := 0; i < count; i++ {
					xlsx.SetCellValue(v.SheetName, fmt.Sprintf("A%d", i + step), vvv.MyInsName)
					xlsx.SetCellValue(v.SheetName, fmt.Sprintf("B%d", i + step), items[i].Sample)
					xlsx.SetCellValue(v.SheetName, fmt.Sprintf("C%d", i + step), items[i].TsCnt)
					xlsx.SetCellValue(v.SheetName, fmt.Sprintf("D%d", i + step), items[i].QueryTimePct95.String())
					xlsx.SetCellValue(v.SheetName, fmt.Sprintf("E%d", i + step), items[i].TsMax.Format("2006-01-02 15:04:05"))
					xlsx.SetCellValue(v.SheetName, fmt.Sprintf("F%d", i + step), items[i].UserMax)
					xlsx.SetCellValue(v.SheetName, fmt.Sprintf("G%d", i + step), vvv.Charge)
				}
			}
			step = step + count
		}
	}
	xlsx.DeleteSheet("Sheet1")

	return xlsx.SaveAs(file)
}

func sendMysqlSlowlogDailyReport() {
	slowInfo := createReport()
	reportFile := fmt.Sprintf("%s.html", utils.YesterdayString())
	utils.SaveReport(fmt.Sprintf("./report/%s", reportFile), slowInfo)
	sendFileList := ""
	fileListTgz := fmt.Sprintf("DailySlow.tgz")

	if utils.FirstDayOfMonth(utils.Today()) {
		slowInfoMonthly := createMonthlyReport()
		ds := utils.YearMonthStringByFormat(utils.Today(),"20060102")
		reportFileMonthly := fmt.Sprintf("%s.html", ds)
		utils.SaveReport(fmt.Sprintf("./report/%s", reportFileMonthly), slowInfoMonthly)

		traceExcel := fmt.Sprintf("%s.xlsx", ds)
		err := createTraceExcelReport(fmt.Sprintf("./report/%s", traceExcel))
		if err != nil {
			beego.Info(fmt.Sprintf("createTraceExcelReport error: %s", err.Error()))
		}

		utils.TarFile(fmt.Sprintf("%s %s %s", reportFile, reportFileMonthly, traceExcel),
			fileListTgz, "./report/")
	} else {
		utils.TarFile(fmt.Sprintf("%s", reportFile), fileListTgz, "./report/")
	}
	sendFileList = fmt.Sprintf("%s", fmt.Sprintf("./report/%s", fileListTgz))

	retry := models.MyslowReportRetry()
	body := `<html>
		<head>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
			<style>
				h2 {
					color: blue;
				}

				h3 {
					color: red;
				}

				#lineTitle {
					font-size: 18px;
					font-weight: bold;
				}

			</style>
		</head>
		<body>
			<h2 > 附件说明 </h2>
			<hr>
			<p>
				<span id="lineTitle">每日报告(每天): </span>
				<span>名如20180315.html, 该附件包括了线上每个MySQL实例的慢日志每日报告, 第二天都会报告昨天的慢日志</span>
			</p>
			<p>
				<span id="lineTitle">每月报告(每个月1号): </span>
				<span>慢日志月报告(名如201802.html)以及优化进度表(名如201802.xlsx), 其中优化进度表需要负责人进行跟进并且每个月月底汇总</span>
			</p>
			<hr>
			<h3>PS: 由于附件太大, 所以压缩为tgz包, 解压后方可查看, html文件使用浏览器打开即可查看</h3>
		</body>
		</html>`
	for i := 0; i < retry; i++ {
		ret, err := email.SendEmailWithAddition(
			models.MyslowReportMailUserName(),
			models.MyslowReportMailPassword(),
			models.MyslowReportMailHost(),
			models.MyslowReportMailPort(),
			models.MyslowReportFrom(),
			models.MyslowReportFromAlias(),
			models.MyslowReportTos(),
			models.MyslowReportCcs(),
			models.MyslowReportSubject(),
			sendFileList, body, "html")
		if ret == 0 && err == nil {
			beego.Info("Send email report successfully")
			if utils.FirstDayOfMonth(utils.Today()) {
				err := dingding.SendTextToDingding(models.MyslowReportDdTips(),
					models.MyslowReportAtlist(), false, models.MyslowReportDdWebhook())
				if err != nil {
					beego.Info(fmt.Sprintf("Send dingding tips fail: %s", err.Error()))
				} else {
					beego.Info("Send dingding tips successfully")
				}
			}
			return
		}
		beego.Info(fmt.Sprintf("Send email report fail and try again: %s,%d", err.Error(), ret))
	}

	beego.Info(fmt.Sprintf("Send email report fail with %d tries", retry))
}

func Init() {
	c := cron.New()
	c.AddFunc(models.SyncMysqlSlowlogSpec(), syncMysqlDailySlowlog)
	c.AddFunc(models.SendMysqlSlowlogReportSpec(), sendMysqlSlowlogDailyReport)
	c.AddFunc(models.ImportMysqlSlowlogMonthlySpec(), importMysqlMonthlySlowlog)
	c.Start()
}