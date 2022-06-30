package parse

import (
    "strings"

    "github.com/spf13/cast"
    jsoniter "github.com/json-iterator/go"

    "github.com/deatil/doak-cron/pkg/cmd"
    "github.com/deatil/doak-cron/pkg/curl"
    "github.com/deatil/doak-cron/pkg/utils"
    "github.com/deatil/doak-cron/pkg/logger"
)

// 解析文件
func MakeCron(path string, debug bool) []map[string]func() {
    if !utils.FileExists(path) {
        logger.Log().Error().Msg("配置文件不存在")

        return nil
    }

    data, err := utils.FileRead(path)
    if err != nil {
        logger.Log().Error().Msg(err.Error())

        return nil
    }

    var v []map[string]any
    err = jsoniter.Unmarshal([]byte(data), &v)
    if err != nil {
        logger.Log().Error().Msg(err.Error())

        return nil
    }

    res := make([]map[string]func(), 0)
    if len(v) > 0 {
        for _, vv := range v {
            if _, ok := vv["type"]; !ok {
                continue
            }
            if _, ok2 := vv["spec"]; !ok2 {
                continue
            }

            typ := cast.ToString(vv["type"])
            spec := cast.ToString(vv["spec"])

            switch typ {
                case "cmd":
                    res = append(res, map[string]func(){
                        spec: MakeCmd(vv, debug),
                    })
                case "request":
                    res = append(res, map[string]func(){
                        spec: MakeRequest(vv, debug),
                    })
            }
        }
    }

    return res
}

// 生成脚本
func MakeCmd(data map[string]any, debug bool) func() {
    return func() {
        cmdData := cast.ToString(data["cmd"])
        cmds := strings.Split(cmdData, " ")

        res, err := cmd.Command(cmds[0], cmds[1:]...)
        if err != nil {
            logger.Log().Error().Msg("[cmd]" + err.Error())
        }

        if debug {
            logger.Log().Info().Msg("[debug]cmd: " + cmdData + ", res: " + res)
        }
    }
}

// 生成请求
func MakeRequest(data map[string]any, debug bool) func() {
    return func() {
        if _, ok := data["url"]; !ok {
            return
        }
        if _, ok := data["method"]; !ok {
            return
        }

        opts := make([]curl.Opt, 0)

        method := cast.ToString(data["method"])
        url := cast.ToString(data["url"])

        timeout := cast.ToFloat32(data["timeout"])
        params := cast.ToStringMap(data["params"])
        headers := cast.ToStringMap(data["headers"])
        proxy := cast.ToString(data["proxy"])
        cookies := cast.ToString(data["cookies"])
        charset := cast.ToString(data["charset"])

        if timeout > 0 {
            opts = append(opts, curl.WithTimeout(timeout))
        }
        if len(params) > 0 {
            opts = append(opts, curl.WithFormParams(params))
        }
        if len(headers) > 0 {
            opts = append(opts, curl.WithHeaders(headers))
        }
        if proxy != "" {
            opts = append(opts, curl.WithProxy(proxy))
        }
        if cookies != "" {
            opts = append(opts, curl.WithCookies(cookies))
        }
        if charset != "" {
            opts = append(opts, curl.WithResCharset(charset))
        }

        resp, err := curl.CreateClient().Request(
            strings.ToUpper(method),
            url,
            opts...,
        )
        if err != nil {
            logger.Log().Error().Msg("[request]" + err.Error())
        }

        if debug {
            respData, _ := resp.GetContents()
            logger.Log().Info().Msg("[debug]url: " + url + ", method: " + method + ", res: " + respData)
        }
    }
}
