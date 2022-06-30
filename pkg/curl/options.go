package curl

import (
    "fmt"
)

// Options object
type Options struct {
    Headers    map[string]any
    BaseURI    string
    FormParams map[string]any
    JSON       any
    XML        string
    Timeout    float32
    Cookies    any
    Proxy      string
    ResCharset string
}

// 设置
type Opt func(*Options)

// 设置 BaseURI
func WithBaseURI(data string) Opt {
    return func(opt *Options) {
        opt.BaseURI = data
    }
}

// 设置 Headers
func WithHeaders(data map[string]any) Opt {
    return func(opt *Options) {
        for k, v := range data {
            opt.Headers[k] = fmt.Sprintf("%v", v)
        }
    }
}

// 设置 Timeout
func WithTimeout(data float32) Opt {
    return func(opt *Options) {
        opt.Timeout = data
    }
}

// 设置 Proxy
func WithProxy(data string) Opt {
    return func(opt *Options) {
        opt.Proxy = data
    }
}

// 设置 Cookies
func WithCookies(data any) Opt {
    return func(opt *Options) {
        opt.Cookies = data
    }
}

// 设置 ResCharset
func WithResCharset(data string) Opt {
    return func(opt *Options) {
        opt.ResCharset = data
    }
}

// 设置 FormParams
func WithFormParams(data map[string]any) Opt {
    return func(opt *Options) {
        opt.FormParams = data
    }
}

// 设置 JSON
func WithJSON(data any) Opt {
    return func(opt *Options) {
        opt.JSON = data
    }
}

// 设置 XML
func WithXML(data string) Opt {
    return func(opt *Options) {
        opt.XML = data
    }
}
