package perconatoolkit

import (
	"github.com/ximply/myslowreport/utils"
	"github.com/ximply/myslowreport/models"
	"os/exec"
	"fmt"
)

func ImportMysqlSlowlogHistoryToMysql(file string, table string) bool {
	if len(table) < 1 {
		return false
	}

	if !utils.FileExists(file) {
		return false
	}

	ptQueryDigest := models.MyslowReportPtQueryDigest()

	cmd := exec.Command(ptQueryDigest,
		fmt.Sprintf("--user=%s", models.MyslowReportDbUser()),
		fmt.Sprintf("--password=%s", models.MyslowReportDbPassword()),
		fmt.Sprintf("--port=%s", models.MyslowReportDbPort()),
		"--history",
		fmt.Sprintf("h=%s,D=%s,t=%s", models.MyslowReportDbHost(), models.MyslowReportDbName(), table),
		"--no-report",
		"--limit=100%",
		file)
	cmd.Start()
	cmd.Run()
	cmd.Wait()

	return true
}