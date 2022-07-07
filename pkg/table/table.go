package table

import (
    "os"

    "github.com/spf13/cast"
    "github.com/jedib0t/go-pretty/v6/table"
)

// 显示表格
func ShowTable(settings []map[string]any) {
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

    MakeTable(newSettings)
}

// 显示表格
func MakeTable(data [][]any) {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)

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

    t.Render()
}
