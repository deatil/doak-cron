package curl

import (
    "sync"
    "net/http"
    "net/http/cookiejar"
)

var curSiteCookiesJar, _ = cookiejar.New(nil)
var httpClient = sync.Pool{
    New: func() any {
        return &http.Client{
            Jar: curSiteCookiesJar,
        }
    },
}

// 默认headers头信息，尽可能伪装成为真实的浏览器
var defaultHeader = map[string]any{
    "User-Agent":   "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.81 Safari/537.36 SE 2.X MetaSr 1.0",
    "Accept":       "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
    "Content-Type": "application/x-www-form-urlencoded;charset=utf-8",
    // 特别提醒：真实的浏览器该值为 Accept-Encoding: gzip, deflate，表示浏览器接受压缩后的二进制，浏览器端再解析为html展示，
    // 但是HttpClient解析就麻烦了，所以必须为空或者不设置该值，接受原始数据。否则很容易出现乱码
    "Accept-Encoding":           "",
    "Accept-Language":           "zh-CN,zh;q=0.9",
    "Upgrade-Insecure-Requests": "1",
    "Connection":                "keep-alive",
    "Cache-Control":             "max-age=0",
    "Host":                      "",
}

// 创建一个 HttpClient 客户端用于发送请求
func CreateClient(opts ...Opt) *Request {
    var client = httpClient.Get().(*http.Client)
    defer httpClient.Put(client)

    req := &Request{
        cli:  client,
        opts: &Options{
            Headers:    make(map[string]any, 0),
            FormParams: make(map[string]any, 0),
        },
    }

    // 默认
    newOpts := []Opt{
        WithHeaders(defaultHeader),
    }

    // 合并默认
    opts = append(newOpts, opts...)

    // 设置
    for _, opt := range opts {
        opt(req.opts)
    }

    req.cookiesJar = curSiteCookiesJar

    return req
}
