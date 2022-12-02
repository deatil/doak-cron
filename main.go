package main

import (
    "os"
    "fmt"
    "log"

    "github.com/urfave/cli/v2"

    "github.com/deatil/doak-cron/pkg/cron"
    "github.com/deatil/doak-cron/pkg/parse"
    "github.com/deatil/doak-cron/pkg/table"
    "github.com/deatil/doak-cron/pkg/logger"
)

// 版本号
var version = "1.0.7"

/**
 * go版本的通用计划任务
 *
 * > go run main.go cron --conf="./cron.json" --debug
 * > go run main.go cron --conf="./cron.json" --log="./cron.log" --debug
 * > go run main.go cron ver
 *
 * > main.exe cron --conf="./cron.json" --debug
 * > main.exe cron --conf="./cron.json" --log="./cron.log" --debug
 * > main.exe cron ver
 *
 * @create 2022-6-29
 * @author deatil
 */
func main() {
    app := cli.NewApp()
    app.EnableBashCompletion = true
    app.Commands = []*cli.Command{
        {
            Name:    "cron",
            Aliases: []string{"c"},
            Usage:   "go版本的通用计划任务",
            Flags: []cli.Flag{
                &cli.BoolFlag{Name: "debug", Aliases: []string{"d"}},
                &cli.StringFlag{Name: "conf", Aliases: []string{"c"}},
                &cli.StringFlag{Name: "log", Aliases: []string{"l"}},
            },
            Action: func(ctx *cli.Context) error {
                // 设置日志存储文件
                log := ctx.String("log")
                if log != "" {
                    logger.WithLogFile(log)
                }

                conf := ctx.String("conf")
                debug := ctx.Bool("debug")

                crons, settings := parse.MakeCron(conf, debug)
                if crons == nil {
                    fmt.Println("配置文件错误")
                    return nil
                }

                // 执行计划任务
                cron.AddCron(func(croner *cron.Cron) {
                    if len(crons) > 0 {
                        for k, c := range crons {
                            cronId, err := croner.AddFunc(c.Spec, c.Cmd)
                            if err != nil{
                                 logger.Log().Error().Msg("[cron]" + err.Error())
                            }

                            settings[k]["cron_id"] = cronId
                        }
                    }

                    fmt.Println("")

                    // 显示详情
                    title := "Doak Cron v" + version
                    table.ShowTable(title, settings)
                })

                return nil
            },
            Subcommands: []*cli.Command{
                {
                    Name:  "ver",
                    Usage: "显示计划任务版本号",
                    Action: func(ctx *cli.Context) error {
                        fmt.Println("计划任务当前版本号为: ", version)

                        return nil
                    },
                },
            },
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
