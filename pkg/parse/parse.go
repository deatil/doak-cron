package parse

import (
    "time"
    "errors"
    "strings"
    "net/http"

    "github.com/spf13/cast"
    "github.com/go-resty/resty/v2"
    jsoniter "github.com/json-iterator/go"

    "github.com/deatil/doak-cron/pkg/cmd"
    "github.com/deatil/doak-cron/pkg/utils"
    "github.com/deatil/doak-cron/pkg/logger"
)

// 解析文件
func MakeCron(path string, debug bool) ([]map[string]func(), []map[string]any) {
    if !utils.FileExists(path) {
        logger.Log().Error().Msg("配置文件不存在")

        return nil, nil
    }

    data, err := utils.FileRead(path)
    if err != nil {
        logger.Log().Error().Msg(err.Error())

        return nil, nil
    }

    var v []map[string]any
    err = jsoniter.Unmarshal([]byte(data), &v)
    if err != nil {
        logger.Log().Error().Msg(err.Error())

        return nil, nil
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

    return res, v
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

        method := cast.ToString(data["method"])
        url := cast.ToString(data["url"])

        timeout := cast.ToDuration(data["timeout"])
        params := cast.ToStringMapString(data["params"])
        queryString := cast.ToString(data["query_string"])
        headers := cast.ToStringMapString(data["headers"])
        proxyData := cast.ToString(data["proxy"])
        cookies := cast.ToStringMapString(data["cookies"])
        files := cast.ToStringMapString(data["files"])
        formDatas := cast.ToStringMapString(data["form_data"])

        body, bodyOk := data["body"]

        // 创建客户端
        client := resty.New()

        // 错误处理
        client.OnError(func(req *resty.Request, err error) {
            msg := ""
            if v, ok := err.(*resty.ResponseError); ok {
                msg = v.Err.Error()
            } else {
                msg = err.Error()
            }

            logger.Log().Error().Msg("[request]" + msg)
        })

        // 过期时间
        if timeout > 0 {
            client.SetTimeout(time.Duration(timeout * time.Minute))
        }

        // 代理
        if proxyData != "" {
            // 设置代理
            // proxy = "127.0.0.1:9150"
            client.SetProxy(proxyData)

            // 移除代理
            // client.RemoveProxy()
        }

        // 请求前
        client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
            return nil
        })

        // 相应后
        client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
            return nil
        })

        // 重试
        client.
            SetRetryCount(3).
            SetRetryWaitTime(5 * time.Second).
            SetRetryMaxWaitTime(20 * time.Second).
            SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
                return 0, errors.New("quota exceeded")
            })

        // 调试
        if debug {
            // client.SetDebug(true)
        }

        // 创建请求
        r := client.R()

        // 设置内容长度
        r.SetContentLength(true)

        if len(params) > 0 {
            // r.SetQueryParam("size", "large")
            r.SetQueryParams(params)
        }
        if queryString != "" {
            r.SetQueryString(queryString)
        }
        if bodyOk {
            r.SetBody(body)
        }
        if len(headers) > 0 {
            r.SetHeaders(headers)
        }
        if len(cookies) > 0 {
            for k, v := range cookies{
                r.SetCookie(&http.Cookie{
                    Name: k,
                    Value: v,
                })
            }
            // r.SetCookies(cookies)
        }

        // 文件上传
        if len(files) > 0 {
            // notesBytes, _ := ioutil.ReadFile("/Users/deatil/text-file.txt")
            // r.SetFileReader("notes", "text-file.txt", bytes.NewReader(notesBytes))

            // r.SetFile("profile_img", "/Users/deatil/test-img.png")
            r.SetFiles(files)
        }

        // 表单
        if len(formDatas) > 0 {
            r.SetFormData(formDatas)
        }

        // 方法名大写
        method = strings.ToUpper(method)

        // 请求
        resp, err := r.Execute(method, url)
        if err != nil {
            logger.Log().Error().Msg("[request]" + err.Error())
        }

        if debug {
            respData := resp.Body()
            logger.Log().Info().Msg("[debug]url: " + url + ", method: " + method + ", res: " + string(respData))
        }
    }
}
