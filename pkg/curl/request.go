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
    "strconv"
    "reflect"
    "crypto/tls"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "strings"
    "encoding/json"
    "html/template"
    "mime/multipart"

    "github.com/axgle/mahonia"
)

// 请求
type Request struct {
    opts       *Options
    cli        *http.Client
    req        *http.Request
    body       io.Reader
    Params     string
    cookiesJar *cookiejar.Jar
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
func (this *Request) Request(method string, uri string, opts ...Opt) (*Response, error) {
    if len(opts) > 0 {
        for _, opt := range opts{
            opt(this.opts)
        }
    }

    // http.MethodGet, http.MethodDelete
    // http.MethodPost, http.MethodPut,
    // http.MethodPatch, http.MethodOptions

    // 解析链接
    url := this.parseUrl(uri, this.parseParams())

    // 解析内容
    this.parseBody()

    // 请求
    req, err := http.NewRequest(method, url, this.body)
    if err != nil {
        return nil, err
    }

    this.req = req

    this.opts.Headers["Host"] = fmt.Sprintf("%v", this.req.Host)

    // parseTimeout
    this.parseTimeout()

    // parseClient
    this.parseClient()

    // parse headers
    this.parseHeaders()

    // parse cookies
    this.parseCookies()

    // 执行请求
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

    // 额外设置
    tr.MaxIdleConns = this.opts.MaxIdleConns
    tr.MaxConnsPerHost = this.opts.MaxConnsPerHost
    tr.MaxIdleConnsPerHost = this.opts.MaxIdleConnsPerHost

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
    if this.opts.Body == nil {
        return
    }

    switch this.opts.Body.(type) {
        // application/x-www-form-urlencoded
        // application/json
        case map[string]any:
            data := this.opts.Body.(map[string]any)

            values := url.Values{}
            for k, v := range data {
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
        // 上传文件
        case string:
            data := this.opts.Body.(string)

            if strings.HasPrefix(data, "@") {
                newData := data[1:]

                newDatas := strings.Split(newData, "|")
                if len(newDatas) != 2 {
                    fmt.Println("parseBody err: 数据错误")
                    return
                }

                paramName := newDatas[0]
                filePath := newDatas[1]

                // 打开要上传的文件
                file, err := os.Open(filePath)
                if err != nil {
                    fmt.Println("parseBody err: ", err)
                    return
                }

                defer file.Close()

                body := &bytes.Buffer{}
                // 创建一个multipart类型的写文件
                writer := multipart.NewWriter(body)

                // 使用给出的属性名paramName和文件名filePath创建一个新的form-data头
                part, err := writer.CreateFormFile(paramName, filePath)
                if err != nil {
                    fmt.Println("parseBody err: ", err)
                    return
                }

                // 将源复制到目标，将file写入到part
                // 是按默认的缓冲区32k循环操作的，
                // 不会将内容一次性全写入内存中,
                // 这样就能解决大文件的问题
                _, err = io.Copy(part, file)
                err = writer.Close()
                if err != nil {
                    fmt.Println("parseBody err: ", err)
                    return
                }

                this.body = body

                this.opts.Headers["Content-Type"] = writer.FormDataContentType()
            } else {
                this.body = strings.NewReader(data)
            }

        // text/xml
        default:
            data := this.ToString(this.opts.Body)

            this.body = strings.NewReader(data)
    }
}

// 解析 get 方式传递的 formData(application/x-www-form-urlencoded)
func (this *Request) parseParams() string {
    if this.opts.Params == nil {
        return ""
    }

    values := url.Values{}
    for k, v := range this.opts.Params {
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
    this.Params = values.Encode()

    return this.Params
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

func (this *Request) indirectToStringerOrError(a interface{}) interface{} {
    if a == nil {
        return nil
    }

    var errorType = reflect.TypeOf((*error)(nil)).Elem()
    var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

    v := reflect.ValueOf(a)
    for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
        v = v.Elem()
    }

    return v.Interface()
}

// 转换为字符
func (this *Request) ToString(i interface{}) string {
    i = this.indirectToStringerOrError(i)

    switch s := i.(type) {
        case string:
            return s
        case bool:
            return strconv.FormatBool(s)
        case float64:
            return strconv.FormatFloat(s, 'f', -1, 64)
        case float32:
            return strconv.FormatFloat(float64(s), 'f', -1, 32)
        case int:
            return strconv.Itoa(s)
        case int64:
            return strconv.FormatInt(s, 10)
        case int32:
            return strconv.Itoa(int(s))
        case int16:
            return strconv.FormatInt(int64(s), 10)
        case int8:
            return strconv.FormatInt(int64(s), 10)
        case uint:
            return strconv.FormatUint(uint64(s), 10)
        case uint64:
            return strconv.FormatUint(uint64(s), 10)
        case uint32:
            return strconv.FormatUint(uint64(s), 10)
        case uint16:
            return strconv.FormatUint(uint64(s), 10)
        case uint8:
            return strconv.FormatUint(uint64(s), 10)
        case json.Number:
            return s.String()
        case []byte:
            return string(s)
        case template.HTML:
            return string(s)
        case template.URL:
            return string(s)
        case template.JS:
            return string(s)
        case template.CSS:
            return string(s)
        case template.HTMLAttr:
            return string(s)
        case nil:
            return ""
        case fmt.Stringer:
            return s.String()
        case error:
            return s.Error()
        default:
            return ""
    }
}
