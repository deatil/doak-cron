package cron

import (
    "fmt"

    "github.com/robfig/cron/v3"
)

type Option struct {
    // 时间
    Spec string

    // 脚本
    Cmd func()
}

// 添加计划任务
func AddCron(fn func(*cron.Cron)) {
    // 创建一个cron实例
    newCron := cron.New(cron.WithSeconds())

    // 添加计划任务
    fn(newCron)

    // 启动/关闭
    newCron.Start()
    defer newCron.Stop()

    // 查询语句，保持程序运行，在这里等同于for{}
    select {}
}

// 添加计划任务
func AddCrons(opts []Option) {
    AddCron(func(croner *cron.Cron) {
        if len(opts) > 0 {
            for _, opt := range opts {
                _, err := croner.AddFunc(opt.Spec, opt.Cmd)
                if err != nil{
                    fmt.Println(err)
                }
            }
        }
    })
}
