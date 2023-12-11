package system

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

// cronCmd represents the cron command
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "crontab 表达式",
	Long: `详细使用:

1. linux cron description:
# .---------------- minute (0 - 59)
# | .------------- hour (0 - 23)
# | | .---------- day of month (1 - 31)
# | | | .------- month (1 - 12) OR jan,feb,mar,apr ...
# | | | | .---- day of week (0 - 6) (Sunday=0 or 7) OR sun,mon,tue,wed,thu,fri,sat
# | | | | |
# * * * * * user-name command to be executed

2. 常用的cron表达式

实例1：每1分钟执行一次myCommand
* * * * * myCommand

实例2：每小时的第3和第15分钟执行
3,15 * * * * myCommand

实例3：在上午8点到11点的第3和第15分钟执行
3,15 8-11 * * * myCommand

实例4：每隔两天的上午8点到11点的第3和第15分钟执行
3,15 8-11 */2  *  * myCommand

实例5：每周一上午8点到11点的第3和第15分钟执行
3,15 8-11 * * 1 myCommand

实例6：每晚的21:30重启smb
30 21 * * * /etc/init.d/smb restart

实例7：每月1、10、22日的4 : 45重启smb
45 4 1,10,22 * * /etc/init.d/smb restart

实例8：每周六、周日的1 : 10重启smb
10 1 * * 6,0 /etc/init.d/smb restart

实例9：每天18 : 00至23 : 00之间每隔30分钟重启smb
0,30 18-23 * * * /etc/init.d/smb restart

实例10：每星期六的晚上11 : 00 pm重启smb
0 23 * * 6 /etc/init.d/smb restart

实例11：每一小时重启smb
0 */1 * * * /etc/init.d/smb restart

实例12：晚上11点到早上7点之间，每隔一小时重启smb
0 23-7/1 * * * /etc/init.d/smb restart.`,

	Run: func(_ *cobra.Command, arg []string) {
		for _, value := range arg {
			withLocation := fmt.Sprintf("CRON_TZ=%s %s", "Asia/Shanghai", value)
			schedule, err := cron.ParseStandard(withLocation)
			if err != nil {
				fmt.Println("parse error:", err)
				return
			}
			record := time.Now()
			fmt.Println("cron: ", value)
			for i := 0; i < 10; i++ {
				record = schedule.Next(record)
				fmt.Println(record)
			}
		}
	},
}

func init() {
	cronCmd.Flags().IntP("time", "t", 10, "cron执行次数")
}

func NewCronCmd() *cobra.Command {
	return cronCmd
}
