package table

import (
    "os"
    "time"

    "github.com/spf13/cast"
    "github.com/jedib0t/go-pretty/v6/table"
)

// 显示表格
func ShowTable(title string, settings []map[string]any) {
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

            typ := cast.ToString(v["type"])
            spec := cast.ToString(v["spec"])

            newSettings = append(newSettings, []any{k+1, typ, spec, status})
        }
    } else {
        newSettings = append(newSettings, []any{1, "none", "-", "stop"})
    }

    MakeTable(title, newSettings)
}

// 显示表格
func MakeTable(title string, data [][]any) {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)

    t.SetTitle(title)

    t.AppendHeader(table.Row{
        "#", "Type", "Spec", "Status",
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

    loc, _ := time.LoadLocation("Asia/Shanghai")
    nowTime := time.Now().
        In(loc).
        Format("2006-01-02 15:04:05")

    caption := "Start At " + nowTime + ".\n"
    t.SetCaption(caption)

    t.Render()
}
