package rsync

import (
	"fmt"
	"os/exec"
	"github.com/ximply/myslowreport/utils"
)

func SyncMysqlSlowlogFile(ip string, port string, model string, remoteFile string, logPath string) string {
	yestoday := utils.YesterdayString()
	localFile := fmt.Sprintf("%s%s.%s.%s", logPath, ip, port, yestoday)
	sourceFile := fmt.Sprintf("%s.%s", remoteFile, yestoday)
	
	cmd := exec.Command("rsync",
		"-av",
			fmt.Sprintf("root@%s::%s/%s", ip, model, sourceFile),
			fmt.Sprintf("%s", localFile))
	cmd.Start()
	cmd.Run()
	cmd.Wait()

	return localFile
}
