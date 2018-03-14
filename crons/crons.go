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
	"strconv"
	"github.com/ximply/myslowreport/email"
)

func syncMysqlSlowlog() {
	pl := models.MyslowReportProjectsList()
	for _, p := range pl {
		beego.Info(fmt.Sprintf("Start SyncMysqlSlowlogFile: %s, %s", p.MysqlHost, p.MysqlSlowlogFileName))
		file := rsync.SyncMysqlSlowlogFile(
			p.MysqlHost,
			p.MysqlPort,
			p.RsyncModel,
			p.MysqlSlowlogFileName,
			models.MyslowReportSlowlogpath())
		beego.Info(fmt.Sprintf("End SyncMysqlSlowlogFile: %s, %s", p.MysqlHost, p.MysqlSlowlogFileName))

		if !perconatoolkit.ImportMysqlSlowlogHistoryToMysql(file, p.SlowlogTable) {
			beego.Info(fmt.Sprintf("ImportMysqlSlowlogHistoryToMysql fail: %s, %s", file, p.SlowlogTable))
			continue
		}
		beego.Info(fmt.Sprintf("ImportMysqlSlowlogHistoryToMysql successfully: %s, %s", file, p.SlowlogTable))
	}
}

func createReport() string {
	pl := models.MyslowReportProjectsList()
	slowInfo := ""
	slowInfo += `<html>`
	slowInfo +=
		`<head>
    		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    		<style class="fox_global_style">
        		div.fox_html_content {
            		line-height: 1.5;
        		}

        		blockquote {
            		margin-Top: 0px;
            		margin-Bottom: 0px;
            		margin-Left: 0.5em
        		}

        		ol,
        		ul {
            		margin-Top: 0px;
            		margin-Bottom: 0px;
            		list-style-position: inside;
        		}

        		p {
            		margin-Top: 0px;
            		margin-Bottom: 0px
        		}
    		</style>
			<script type="text/javascript">
				function toggle(id) {
 					var tb=document.getElementById(id);
 					if (tb.style.display=='none') tb.style.display='block';
 					else tb.style.display='none';
				}
			</script>
		</head>`

	slowInfo += `<body>`

	slowInfo += `<div id="content" class="reg">`
	slowInfo += fmt.Sprintf(`<title>%s</title>`, models.MyslowReportTitle())
	slowInfo += `<div style="LINE-HEIGHT: 20px; FONT-FAMILY: 'Microsoft Yahei'; FONT-SIZE: 12px">`
	slowInfo += `<div style="MARGIN: 0px auto; WIDTH: 1280px; BACKGROUND: #FF710C; border-top-right-radius: 8px; border-top-left-radius: 8px; border-bottom-right-radius: 8px; border-bottom-left-radius: 8px; -moz-border-radius: 8px; -webkit-border-radius: 8px">`
	slowInfo += `<div style="WIDTH: 710px; HEIGHT: 10px"></div>`
	slowInfo += `<div style="PADDING-BOTTOM: 0px; PADDING-LEFT: 0px; WIDTH: 710px; PADDING-RIGHT: 0px; BACKGROUND: #FF710C; HEIGHT: 30px; PADDING-TOP: 10px">`
	slowInfo += fmt.Sprintf(`<p style="MARGIN: 0px; PADDING-LEFT: 30px; FONT-FAMILY: 'Microsoft Yahei'; COLOR: rgb(255,255,255); FONT-SIZE: 24px">%s %s %s</p>`,
		models.MyslowReportSubject(), utils.YesterdayStringByFormat("2006-01-02"), utils.WeekdayCNString(utils.Yesterday()))
	slowInfo += `</div>`

	slowInfo += `<div style="BORDER-BOTTOM: #FF710C 1px solid; BORDER-LEFT: #FF710C 1px solid; PADDING-BOTTOM: 5px; PADDING-LEFT: 9px; PADDING-RIGHT: 9px; BACKGROUND: #f9f9f9; BORDER-TOP: #FF710C 1px solid; BORDER-RIGHT: #FF710C 1px solid; PADDING-TOP: 5px">`
	slowInfo += `<div style="PADDING-BOTTOM: 0px; MARGIN: 10px 0px 0px; PADDING-LEFT: 15px; PADDING-RIGHT: 0px; PADDING-TOP: 0px">`

	for i, p := range pl {
		beego.Info(fmt.Sprintf("Start GetOrderByQueryTimeMaxDesc: %s", p.SlowlogTable))
		cnt, items := models.GetOrderByQueryTimeMaxDesc(
			fmt.Sprintf("%s 23:59:59", utils.BeforeYesterdayStringByFormat("2006-01-02")),
			fmt.Sprintf("%s 00:00:00", utils.TodayStringByFormat("2006-01-02")),
			p.SlowlogTable)
		beego.Info(fmt.Sprintf("End GetOrderByQueryTimeMaxDesc: %s, Count: %d", p.SlowlogTable, cnt))

		if cnt == 0 {
			beego.Info(fmt.Sprintf("No slow logs: %s", p.SlowlogTable))
			continue
		}

		beego.Info(fmt.Sprintf("Creating report content: %s", p.SlowlogTable))
		slowInfo += createSlowInfo(items, p, i)
		slowInfo += "<br/>"
		slowInfo += `<hr style="#FF710C;" />`
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

func createSlowInfo(items []models.Item, p models.Project, tableId int) string {
	var info = ""

	info += `<div style="PADDING-BOTTOM: 0px; MARGIN: 10px 0px 0px; PADDING-LEFT: 15px; PADDING-RIGHT: 0px; PADDING-TOP: 0px">`
	info += `<p style="MARGIN-TOP: 0px; MARGIN-BOTTOM: 0px; MARGIN-LEFT: 7px; FONT-SIZE: 16px; TEXT-DECORATION: none"><span style="FLOAT: left; COLOR: #262626; FONT-SIZE: 18px">&nbsp;&nbsp; </span> <span style="FLOAT: left; COLOR: #262626; FONT-SIZE: 18px; font-weight:bold">192.168.10.1</span><span style="COLOR: rgb(0,0,0); FONT-SIZE: 15px">&nbsp;&nbsp; Project</span></p>`
	info = strings.Replace(info, "192.168.10.1", p.MysqlHost, -1)
	info = strings.Replace(info, "Project", p.Description, -1)

	yesterdayTotal := models.GetSumOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.BeforeYesterdayStringByFormat("2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.TodayStringByFormat("2006-01-02")),
		p.SlowlogTable)
	beforeYesterdayTotal := models.GetSumOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.BeforeBeforeYesterdayStringByFormat("2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.YesterdayStringByFormat("2006-01-02")),
		p.SlowlogTable)
	yoyBasisTotal := models.GetSumOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.DateStringByFormat(-9,"2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.DateStringByFormat(-7,"2006-01-02")),
		p.SlowlogTable)

	yesterdayUniq := models.GetUniqOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.BeforeYesterdayStringByFormat("2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.TodayStringByFormat("2006-01-02")),
		p.SlowlogTable)
	beforeYesterdayUniq := models.GetUniqOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.BeforeBeforeYesterdayStringByFormat("2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.YesterdayStringByFormat("2006-01-02")),
		p.SlowlogTable)
	yoyBasisUniq := models.GetUniqOfQueryCount(fmt.Sprintf("%s 23:59:59", utils.DateStringByFormat(-9,"2006-01-02")),
		fmt.Sprintf("%s 00:00:00", utils.DateStringByFormat(-7,"2006-01-02")),
		p.SlowlogTable)

	// chain grow with before yesterday and on year-on-year basis (this day last week)
	chainGrowthRateForTotal := 0.0
	if  beforeYesterdayTotal != 0 {
		chainGrowthRateForTotal = (float64(yesterdayTotal - beforeYesterdayTotal)) * 100 / float64(beforeYesterdayTotal)
	}
	chainGrowthRateForTotalStr := ""
	chainGrowthRateForTotalStrColor := "black"
	if chainGrowthRateForTotal > 0.0 {
		chainGrowthRateForTotalStr = fmt.Sprintf("+%.2f%%", chainGrowthRateForTotal)
		chainGrowthRateForTotalStrColor = "red"
	} else if chainGrowthRateForTotal < 0.0 {
		chainGrowthRateForTotalStr = fmt.Sprintf("%.2f%%", chainGrowthRateForTotal)
		chainGrowthRateForTotalStrColor = "green"
	} else {
		if beforeYesterdayTotal == 0{
			chainGrowthRateForTotalStr = "NULL"
			chainGrowthRateForTotalStrColor = "gray"
		} else {
			chainGrowthRateForTotalStr = "--"
		}
	}

	yoyBasisRateForTotal := 0.0
	if  yoyBasisTotal != 0 {
		yoyBasisRateForTotal = (float64(yesterdayTotal - yoyBasisTotal)) * 100 / float64(yoyBasisTotal)
	}
	yoyBasisRateForTotalStr := ""
	yoyBasisRateForTotalStrColor := "black"
	if yoyBasisRateForTotal > 0.0 {
		yoyBasisRateForTotalStr = fmt.Sprintf("+%.2f%%", yoyBasisRateForTotal)
		yoyBasisRateForTotalStrColor = "red"
	} else if yoyBasisRateForTotal < 0.0 {
		yoyBasisRateForTotalStr = fmt.Sprintf("%.2f%%", yoyBasisRateForTotal)
		yoyBasisRateForTotalStrColor = "green"
	} else {
		if yoyBasisTotal == 0 {
			yoyBasisRateForTotalStr = "NULL"
			yoyBasisRateForTotalStrColor = "gray"
		} else {
			yoyBasisRateForTotalStr = "--"
		}
	}

	chainGrowthRateForUniq := 0.0
	if beforeYesterdayUniq != 0 {
		chainGrowthRateForUniq = (float64(yesterdayUniq - beforeYesterdayUniq)) * 100 / float64(beforeYesterdayUniq)
	}
	chainGrowthRateForUniqStr := ""
	chainGrowthRateForUniqStrColor := "black"
	if chainGrowthRateForUniq > 0.0 {
		chainGrowthRateForUniqStr = fmt.Sprintf("+%.2f%%", chainGrowthRateForUniq)
		chainGrowthRateForUniqStrColor = "red"
	} else if chainGrowthRateForUniq < 0.0 {
		chainGrowthRateForUniqStr = fmt.Sprintf("%.2f%%", chainGrowthRateForUniq)
		chainGrowthRateForUniqStrColor = "green"
	} else {
		if beforeYesterdayUniq == 0 {
			chainGrowthRateForUniqStr = "NULL"
			chainGrowthRateForUniqStrColor = "gray"
		} else {
			chainGrowthRateForUniqStr = "--"
		}
	}

	yoyBasisRateForUniq := 0.0
	if yoyBasisUniq != 0 {
		yoyBasisRateForUniq = (float64(yesterdayUniq - yoyBasisUniq)) * 100 / float64(yoyBasisUniq)
	}
	yoyBasisRateForUniqStr := ""
	yoyBasisRateForUniqStrColor := "black"
	if yoyBasisRateForUniq > 0.0 {
		yoyBasisRateForUniqStr = fmt.Sprintf("+%.2f%%", yoyBasisRateForUniq)
		yoyBasisRateForUniqStrColor = "red"
	} else if yoyBasisRateForUniq < 0.0 {
		yoyBasisRateForUniqStr = fmt.Sprintf("%.2f%%", yoyBasisRateForUniq)
		yoyBasisRateForUniqStrColor = "green"
	} else {
		if yoyBasisUniq == 0 {
			yoyBasisRateForUniqStr = "NULL"
			yoyBasisRateForUniqStrColor = "gray"
		} else {
			yoyBasisRateForUniqStr = "--"
		}
	}


	info += `<p style="MARGIN-TOP: 0px; MARGIN-BOTTOM: 0px; MARGIN-LEFT: 7px; FONT-SIZE: 16px; TEXT-DECORATION: none"><span style="COLOR: rgb(0,0,0); FONT-SIZE: 15px">&nbsp;&nbsp; TotalCount &nbsp;&nbsp;</span><span style="COLOR: ChainColor; FONT-SIZE: 15px">&nbsp;&nbsp; TChain &nbsp;&nbsp;</span></span><span style="COLOR: YOYColor; FONT-SIZE: 15px">&nbsp;&nbsp; YOYTotal &nbsp;&nbsp;</span></p>`
	info = strings.Replace(info, "TotalCount", fmt.Sprintf("总语句数: %d", yesterdayTotal), -1)
	info = strings.Replace(info, "TChain", fmt.Sprintf("环比前天: %s", chainGrowthRateForTotalStr), -1)
	info = strings.Replace(info, "ChainColor", fmt.Sprintf("%s", chainGrowthRateForTotalStrColor), -1)
	info = strings.Replace(info, "YOYTotal", fmt.Sprintf("同比上周%s: %s", utils.WeekdayCNShortString(utils.Yesterday()), yoyBasisRateForTotalStr), -1)
	info = strings.Replace(info, "YOYColor", fmt.Sprintf("%s", yoyBasisRateForTotalStrColor), -1)
	
	info += `<p style="MARGIN-TOP: 0px; MARGIN-BOTTOM: 0px; MARGIN-LEFT: 7px; FONT-SIZE: 16px; TEXT-DECORATION: none"><span style="COLOR: rgb(0,0,0); FONT-SIZE: 15px">&nbsp;&nbsp; UniqCount &nbsp;&nbsp;</span><span style="COLOR: ChainColor; FONT-SIZE: 15px">&nbsp;&nbsp; UChain &nbsp;&nbsp;</span></span><span style="COLOR: YOYColor; FONT-SIZE: 15px">&nbsp;&nbsp; YOYUniq &nbsp;&nbsp;</span></p>`
	info = strings.Replace(info, "UniqCount", fmt.Sprintf("独立语句数: %d", yesterdayUniq), -1)
	info = strings.Replace(info, "UChain", fmt.Sprintf("环比前天: %s", chainGrowthRateForUniqStr), -1)
	info = strings.Replace(info, "ChainColor", fmt.Sprintf("%s", chainGrowthRateForUniqStrColor), -1)
	info = strings.Replace(info, "YOYUniq", fmt.Sprintf("同比上周%s: %s", utils.WeekdayCNShortString(utils.Yesterday()), yoyBasisRateForUniqStr), -1)
	info = strings.Replace(info, "YOYColor", fmt.Sprintf("%s", yoyBasisRateForUniqStrColor), -1)
	
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
	// col name,col width
	var colsWidth = []string{}
	for _, col := range cols {
		l := strings.Split(col, ",")
		if len(l) != 2 {
			return ""
		}
		colName := l[0]
		width, _ := strconv.ParseInt(l[1], 10,32)
		if width < 0 || width > 10000 {
			width = 100
		}
		colWidth := fmt.Sprintf("%d", width)
		colsWidth = append(colsWidth, colWidth)
		tmp := `<td style="BORDER-BOTTOM: rgb(222,222,222) 1px solid; TEXT-ALIGN: center; BORDER-LEFT: rgb(222,222,222) 1px solid; PADDING-BOTTOM: 7px; MARGIN: 0px; PADDING-LEFT: 7px; min-height: 88888888px; MAX-WIDTH: 999999999px; PADDING-RIGHT: 15px; HEIGHT: 30px; COLOR: rgb(255,255,255); FONT-SIZE: 13px; BORDER-RIGHT: rgb(241,241,226) 1px solid; PADDING-TOP: 7px">SlowQuerySample</td>`		
		tmp = strings.Replace(tmp, "999999999", colWidth, -1)		
		tmp = strings.Replace(tmp, "88888888", "50", -1)		
		tmp = strings.Replace(tmp, "SlowQuerySample", colName, -1)
		info += tmp
	}
	info += `</tr>`

	for _, i := range items {
		info += `<tr style="DISPLAY: table-row; VERTICAL-ALIGN: inherit">`
		info += createItem(i, colsWidth)
		info += `</tr>`
	}

	info += `</tbody>`
	info += `</table>`
	info += `</div>`
	info += `</div>`
	info += `</div>`

	return info
}

func createItem(item models.Item, colsWidth []string) string {
	defaultItem := `<td style="BORDER-BOTTOM: rgb(222,222,222) 1px solid; TEXT-ALIGN: left; BORDER-LEFT: rgb(222,222,222) 1px solid; PADDING-BOTTOM: 7px; MARGIN: 0px; PADDING-LEFT: 7px; MAX-WIDTH: 999999999px; PADDING-RIGHT: 15px; FONT-SIZE: 13px; BORDER-RIGHT: rgb(241,241,226) 1px solid; PADDING-TOP: 7px;">xxx</td>`
	defaultItemSample :=  `<td style="BORDER-BOTTOM: rgb(222,222,222) 1px solid; TEXT-ALIGN: left; BORDER-LEFT: rgb(222,222,222) 1px solid; PADDING-BOTTOM: 7px; MARGIN: 0px; PADDING-LEFT: 7px; min-height: 88888888px; MAX-WIDTH: 999999999px; PADDING-RIGHT: 15px; FONT-SIZE: 13px; BORDER-RIGHT: rgb(241,241,226) 1px solid; PADDING-TOP: 7px;"><textarea readonly="readonly" style="max-width:200px; min-height: 150px; max-height:180px;">xxx</textarea></td>`	
	
	s := ""
	i := 0

	tmp := defaultItemSample
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp = strings.Replace(tmp, "xxx", item.Sample, -1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp = strings.Replace(tmp, "xxx",
		strings.TrimRight(fmt.Sprintf("%.1f", item.TsCnt), ".0"), -1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp = strings.Replace(tmp, "xxx", item.UserMax, -1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp  = strings.Replace(tmp, "xxx",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", item.QueryTimeMin), "0"), "."), -1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp  = strings.Replace(tmp, "xxx",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", item.QueryTimeMax),"0"),"."),-1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp = strings.Replace(tmp, "xxx",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", item.QueryTimePct95),"0"),"."),-1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp = strings.Replace(tmp, "xxx",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", item.LockTimeMin),"0"),"."),-1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp = strings.Replace(tmp, "xxx",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", item.QueryTimeMax),"0"),"."),-1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp = strings.Replace(tmp, "xxx",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", item.QueryTimePct95),"0"),"."),-1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp = strings.Replace(tmp, "xxx",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", item.RowsExaminedMin),"0"),"."),-1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp = strings.Replace(tmp, "xxx",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", item.RowsExaminedMax),"0"),"."),-1)
	s += tmp
	i++

	tmp = defaultItem
	tmp = strings.Replace(tmp, "999999999", colsWidth[i], -1)
	tmp = strings.Replace(tmp, "88888888", "50", -1)
	tmp = strings.Replace(tmp, "xxx",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", item.RowsExaminedPct95), "0"),"."),-1)
	s += tmp
	i++

	return s
}

func sendMysqlSlowlogReport() {
	slowInfo := createReport()
	reportFile := fmt.Sprintf("./report/%s.html", utils.YesterdayString())
	utils.SaveReport(reportFile, slowInfo)
	retry := models.MyslowReportRetry()
	for i := 0; i < retry; i++ {
		ret, err := email.SendEmailWithAddition(models.MyslowReportFrom(),
			models.MyslowReportMailUserName(),
			models.MyslowReportMailPassword(),
			models.MyslowReportMailHost(),
			models.MyslowReportMailPort(),
			models.MyslowReportTos(),
			models.MyslowReportSubject(),
			reportFile, "查看详情,请双击打开或者使用浏览器打开附件(推荐使用浏览器打开,Foxmail有些浏览器功能不支持)", "plain")
		if ret == 0 && err == nil {
			beego.Info("Send email report successfully")
			return
		}
		beego.Info(fmt.Sprintf("Send email report fail and try again: %s,%d", err.Error(), ret))
	}

	beego.Info(fmt.Sprintf("Send email report fail with %d tries", retry))
}

func Init() {
	c := cron.New()
	c.AddFunc(models.SyncMysqlSlowlogSpec(), syncMysqlSlowlog)
	c.AddFunc(models.SendMysqlSlowlogReportSpec(), sendMysqlSlowlogReport)
	c.Start()
}

