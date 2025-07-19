package initialize

import (
	"fmt"
	"ygpt/server/task"

	"github.com/robfig/cron/v3"

	"ygpt/server/global"
)

func Timer() {
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFunc("PushResInfo", "*/15 * * * *", func() {
			task.PushResInfo() // 定时任务方法定在task文件包中
		}, "定时推送资源信息", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFunc("PushResInfo", "0 */5 * * * *", func() {
			task.QueryMonitorData() // 定时任务方法定在task文件包中
		}, "定时查询实例监控数据", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}
		// 其他定时任务定在这里 参考上方使用方法

		//_, err := global.GVA_Timer.AddTaskByFunc("定时任务标识", "corn表达式", func() {
		//	具体执行内容...
		//  ......
		//}, option...)
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}
	}()
}
