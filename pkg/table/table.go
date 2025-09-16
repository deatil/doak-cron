package table

import (
    "os"
    "time"

    "github.com/spf13/cast"
    "github.com/jedib0t/go-pretty/v6/table"
)

// 显示表格
func ShowTable(title string, loc string, settings []map[string]any) {
    // 显示计划任务信息
    newSettings := make([][]any, 0)
    if len(settings) > 0 {
        for k, v := range settings {
            status := "run"
            if _, ok := v["type"]; !ok {
                status = "stop"
            }
            if _, ok2 := v["spec"]; !ok2 {
                status = "stop"
            }

            name := cast.ToString(v["name"])
            typ := cast.ToString(v["type"])
            spec := cast.ToString(v["spec"])
            cronId := v["cron_id"]

            if name == "" {
                name = typ
            }

            newSettings = append(newSettings, []any{k+1, name, cronId, typ, spec, status})
        }
    } else {
        newSettings = append(newSettings, []any{1, "none", "0", "none", "-", "stop"})
    }

    makeTable(title, loc, newSettings)
}

// 生成表格
func makeTable(title string, setLoc string, data [][]any) {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)

    t.SetTitle(title)

    t.AppendHeader(table.Row{
        "#", "Name", "Cron_Id", "Type", "Spec", "Status",
    })

    num := len(data)

    var i = 1
    if num > 0 {
        for _, v := range data {
            if i > 1 {
                t.AppendSeparator()
            }

            t.AppendRow(v)

            i++
        }
    }

    if setLoc == "" {
        setLoc = "Asia/Shanghai"
    }

    loc, err := time.LoadLocation(setLoc)
    if err != nil {
        loc = time.UTC
    }

    nowTime := time.Now().
        In(loc).
        Format("2006-01-02 15:04:05 ")

    caption := "Start At " + nowTime + loc.String() + ".\n"
    t.SetCaption(caption)

    t.Render()
}
