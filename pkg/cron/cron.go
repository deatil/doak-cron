package cron

import (
    "fmt"

    "github.com/robfig/cron/v3"
)

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
func AddCrons(fns map[string]func()) {
    AddCron(func(croner *cron.Cron) {
        if len(fns) > 0 {
            for spec, cmd := range fns {
                enterId, err := croner.AddFunc(spec, cmd)
                if err != nil{
                    fmt.Println(err)
                }

                fmt.Printf("任务id为: %d \n", enterId)
            }
        }
    })
}
