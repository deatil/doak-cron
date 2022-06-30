package curl

import (
    "io"
    "os"
    "fmt"
    "path"
    "time"
    "bufio"
    "bytes"
    "errors"
    "crypto/tls"
    "encoding/json"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "strings"

    "github.com/axgle/mahonia"
)

// 请求
type Request struct {
    opts              *Options
    cli               *http.Client
    req               *http.Request
    body              io.Reader
    subFormDataParams string
    cookiesJar        *cookiejar.Jar
}

// Get send get request
func (this *Request) Get(uri string, opts ...Opt) (*Response, error) {
    return this.Request("GET", uri, opts...)
}

// Post send post request
func (this *Request) Post(uri string, opts ...Opt) (*Response, error) {
    return this.Request("POST", uri, opts...)
}

// Put send put request
func (this *Request) Put(uri string, opts ...Opt) (*Response, error) {
    return this.Request("PUT", uri, opts...)
}

// Patch send patch request
func (this *Request) Patch(uri string, opts ...Opt) (*Response, error) {
    return this.Request("PATCH", uri, opts...)
}

// Delete send delete request
func (this *Request) Delete(uri string, opts ...Opt) (*Response, error) {
    return this.Request("DELETE", uri, opts...)
}

// Options send options request
func (this *Request) Options(uri string, opts ...Opt) (*Response, error) {
    return this.Request("OPTIONS", uri, opts...)
}

// Get method download files
func (this *Request) Download(resourceUrl string, savePath, saveName string, opts ...Opt) (bool, error) {
    var vError error
    var vResponse *Response

    uri, err := url.ParseRequestURI(resourceUrl)
    if err != nil {
        return false, err
    }

    if vResponse, vError = this.Request("GET", resourceUrl, opts...); vError == nil {
        filename := path.Base(uri.Path)
        if len(saveName) > 0 {
            filename = saveName
        }

        if vResponse.GetStatusCode() == 200 || vResponse.GetContentLength() > 0 {
            body := vResponse.GetBody()

            return this.saveFile(body, savePath+filename)
        } else {
            return false, errors.New(downloadFileIsEmpty)
        }
    }

    return false, vError
}

func (this *Request) saveFile(body io.ReadCloser, fileName string) (bool, error) {
    var isOccurError bool
    var OccurError error

    file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)

    defer func() {
        body.Close()
        file.Close()
    }()

    reader := bufio.NewReaderSize(body, 1024*50) //相当于一个临时缓冲区(设置为可以单次存储5M的文件)，每次读取以后就把原始数据重新加载一份，等待下一次读取
    if err != nil {
        return false, err
    }

    writer := bufio.NewWriter(file)
    buff := make([]byte, 50*1024)

    for {
        currReadSize, readerErr := reader.Read(buff)
        if currReadSize > 0 {
            _, OccurError = writer.Write(buff[0:currReadSize])

            if OccurError != nil {
                isOccurError = true
                break
            }
        }

        // 读取结束
        if readerErr == io.EOF {
            _ = writer.Flush()
            break
        }
    }

    if isOccurError == false {
        return true, nil
    } else {
        return false, OccurError
    }
}

// Request send request
func (this *Request) Request(method, uri string, opts ...Opt) (*Response, error) {
    if len(opts) > 0 {
        for _, opt := range opts{
            opt(this.opts)
        }
    }

    switch method {
        case http.MethodGet, http.MethodDelete:
            // 解析链接
            url := this.parseUrl(uri, this.parseFormData())

            // 请求
            req, err := http.NewRequest(method, url, nil)
            if err != nil {
                return nil, err
            }

            this.req = req
        case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions:
            // 解析内容
            this.parseBody()

            // 解析链接
            url := this.parseUrl(uri, "")

            // 请求
            req, err := http.NewRequest(method, url, this.body)
            if err != nil {
                return nil, err
            }

            this.req = req
        default:
            return nil, errors.New("invalid request method")
    }

    this.opts.Headers["Host"] = fmt.Sprintf("%v", this.req.Host)

    // parseTimeout
    this.parseTimeout()

    // parseClient
    this.parseClient()

    // parse headers
    this.parseHeaders()

    // parse cookies
    this.parseCookies()

    httpResp, err := this.cli.Do(this.req)

    resp := &Response{
        resp:       httpResp,
        req:        this.req,
        cookiesJar: this.cookiesJar,
        err:        err,
        charset:    this.opts.ResCharset,
    }

    if err != nil {
        return resp, err
    }

    return resp, nil
}

// 格式化链接
func (this *Request) parseUrl(uri string, params string) string {
    url := ""

    if this.opts.BaseURI != "" {
        url = strings.TrimSuffix(this.opts.BaseURI, "/")
        url = url + "/" + strings.TrimPrefix(uri, "/")
    } else {
        url = uri
    }

    if strings.Contains(url, "?") {
        url = url + "&"
    } else {
        url = url + "?"
    }

    url = url + params

    return url
}

func (this *Request) parseTimeout() {
    if this.opts.Timeout == 0 {
        this.opts.Timeout = 30
    }
}

func (this *Request) parseClient() {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    if this.opts.Proxy != "" {
        proxy, err := url.Parse(this.opts.Proxy)
        if err == nil {
            tr.Proxy = http.ProxyURL(proxy)
        } else {
            fmt.Println(this.opts.Proxy + proxyError, err.Error())
        }
    }

    // 过期时间
    timeout := time.Duration(this.opts.Timeout*1000) * time.Millisecond

    this.cli = &http.Client{
        Timeout:   timeout,
        Transport: tr,
        Jar:       this.cookiesJar,
    }
}

func (this *Request) parseCookies() {
    switch this.opts.Cookies.(type) {
        case string:
            cookies := this.opts.Cookies.(string)
            this.req.Header.Add("Cookie", cookies)
        case map[string]string:
            cookies := this.opts.Cookies.(map[string]string)
            for k, v := range cookies {
                if strings.ReplaceAll(v, " ", "") != "" {
                    this.req.AddCookie(&http.Cookie{
                        Name:  k,
                        Value: v,
                    })
                }
            }
        case []*http.Cookie:
            cookies := this.opts.Cookies.([]*http.Cookie)
            for _, cookie := range cookies {
                if cookie != nil {
                    this.req.AddCookie(cookie)
                }
            }
    }
}

func (this *Request) parseHeaders() {
    if this.opts.Headers != nil {
        for k, v := range this.opts.Headers {
            if vv, ok := v.([]string); ok {
                for _, vvv := range vv {
                    if strings.ReplaceAll(vvv, " ", "") != "" {
                        this.req.Header.Add(k, vvv)
                    }
                }
                continue
            }

            vv := fmt.Sprintf("%v", v)
            this.req.Header.Set(k, vv)
        }
    }
}

func (this *Request) parseBody() {
    // application/x-www-form-urlencoded
    if this.opts.FormParams != nil {
        values := url.Values{}
        for k, v := range this.opts.FormParams {
            if vv, ok := v.([]string); ok {
                for _, vvv := range vv {
                    if strings.ReplaceAll(vvv, " ", "") != "" {
                        values.Add(k, vvv)
                    }
                }
                continue
            }
            vv := fmt.Sprintf("%v", v)
            values.Set(k, vv)
        }
        this.body = strings.NewReader(values.Encode())

        return
    }

    // application/json
    if this.opts.JSON != nil {
        b, err := json.Marshal(this.opts.JSON)
        if err == nil {
            this.body = bytes.NewReader(b)

            return
        }
    }

    // text/xml
    if this.opts.XML != "" {
        this.body = strings.NewReader(this.opts.XML)
        return
    }

    return
}

// 解析 get 方式传递的 formData(application/x-www-form-urlencoded)
func (this *Request) parseFormData() string {
    if this.opts.FormParams == nil {
        return ""
    }

    values := url.Values{}
    for k, v := range this.opts.FormParams {
        if vv, ok := v.([]string); ok {
            for _, vvv := range vv {
                if strings.ReplaceAll(vvv, " ", "") != "" {
                    values.Add(k, vvv)
                }
            }
            continue
        }
        vv := fmt.Sprintf("%v", v)
        values.Set(k, vv)
    }
    this.subFormDataParams = values.Encode()

    return this.subFormDataParams
}

//（接受到的）简体中文 转换为 utf-8
func (this *Request) SimpleChineseToUtf8(vBytes []byte) string {
    return mahonia.NewDecoder("GB18030").ConvertString(string(vBytes))
}

// （一般是go 语言发送的数据）utf-8 转换为  简体中文发出去
func (this *Request) Utf8ToSimpleChinese(vBytes []byte, charset ...string) string {
    if len(charset) == 0 {
        return mahonia.NewEncoder("GB18030").ConvertString(string(vBytes))
    } else {
        return mahonia.NewEncoder(charset[0]).ConvertString(string(vBytes))
    }
}
