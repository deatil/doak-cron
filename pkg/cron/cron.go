package cron

import (
    "github.com/robfig/cron/v3"
)

type (
    // 计划任务
    Cron = cron.Cron
)

// 添加计划任务
func AddCron(fn func(*cron.Cron)) {
    logger := cron.PrintfLogger(NewLogger())
    // logger := VerbosePrintfLogger(NewLogger())

    // 创建一个cron实例
    newCron := cron.New(
        cron.WithSeconds(),
        cron.WithLogger(logger),
        cron.WithChain(cron.Recover(logger)),
    )

    // 添加计划任务
    fn(newCron)

    // 启动/关闭
    newCron.Start()
    defer newCron.Stop()

    // 查询语句，保持程序运行，在这里等同于for{}
    select {}
}
