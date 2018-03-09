package rsync

import (
	"fmt"
	"os/exec"
	"time"
)

func SyncMysqlSlowlogFile(ip string, model string, remoteFile string, logPath string) string {
	yestoday := time.Now().AddDate(0, 0, -1).Format("20180101")
	localFile := fmt.Sprintf("%s%s.%s", logPath, ip, yestoday)
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
