package main

import (
    "fmt"
    "log"
    "os"

    "github.com/urfave/cli/v2"

    "github.com/deatil/doak-cron/pkg/cron"
    "github.com/deatil/doak-cron/pkg/logger"
)

// 版本号
var version = "1.0.1"

// go run main.go cron --conf="123.conf"
// go run main.go cron ver
func main() {
    app := cli.NewApp()
    app.EnableBashCompletion = true
    app.Commands = []*cli.Command{
        {
            Name:    "cron",
            Aliases: []string{"c"},
            Usage:   "cron",
            Flags: []cli.Flag{
                &cli.BoolFlag{Name: "debug", Aliases: []string{"d"}},
                &cli.StringFlag{Name: "conf", Aliases: []string{"c"}},
            },
            Action: func(c *cli.Context) error {
                fmt.Println("debug:", c.Bool("debug"))
                fmt.Println("conf:", c.String("conf"))

                cron.AddCrons(map[string]func(){
                    "*/5 * * * * *": func() {
                        fmt.Println("每5s执行一次cron")

                        logger.Log().Info().Msg("每5s执行一次cron")
                    },
                })

                return nil
            },
            Subcommands: []*cli.Command{
                {
                    Name:  "ver",
                    Usage: "显示计划任务版本号",
                    Action: func(c *cli.Context) error {
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
