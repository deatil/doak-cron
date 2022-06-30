package curl

import (
    "fmt"
)

// Options object
type Options struct {
    Headers    map[string]any
    BaseURI    string
    Params     map[string]any
    Body       any
    Timeout    float32
    Cookies    any
    Proxy      string
    ResCharset string

    MaxIdleConns        int
    MaxConnsPerHost     int
    MaxIdleConnsPerHost int
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

// 设置 Params
func WithParams(data map[string]any) Opt {
    return func(opt *Options) {
        opt.Params = data
    }
}

// 设置 Body
func WithBody(data any) Opt {
    return func(opt *Options) {
        opt.Body = data
    }
}

// 设置 MaxIdleConns
func WithMaxIdleConns(data int) Opt {
    return func(opt *Options) {
        opt.MaxIdleConns = data
    }
}

// 设置 MaxConnsPerHost
func WithMaxConnsPerHost(data int) Opt {
    return func(opt *Options) {
        opt.MaxConnsPerHost = data
    }
}

// 设置 MaxIdleConnsPerHost
func WithMaxIdleConnsPerHost(data int) Opt {
    return func(opt *Options) {
        opt.MaxIdleConnsPerHost = data
    }
}
