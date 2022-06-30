package curl

import (
    "io"
    "net"
    "net/http"
    "net/http/cookiejar"
    "fmt"
    "errors"
    "strings"

    "github.com/axgle/mahonia"
)

// 响应
type Response struct {
    resp       *http.Response
    req        *http.Request
    cookiesJar *cookiejar.Jar
    err        error
    charset    string
}

// 获取服务端生成的全部cookies
func (this *Response) GetCookies() []*http.Cookie {
    return this.cookiesJar.Cookies(this.req.URL)
}

// 通过键获取相关的cookie值
func (this *Response) GetCookie(cookieName string) *http.Cookie {
    cookies := this.cookiesJar.Cookies(this.req.URL)
    if len(cookies) > 0 {
        for i := 0; i < len(cookies); i++ {
            if cookies[i].Name == cookieName {
                return cookies[i]
            }
        }
    }

    return nil
}

// GetRequest get request object
func (this *Response) GetRequest() *http.Request {
    return this.req
}

// GetRequest get request object
func (this *Response) GetResponse() *http.Response {
    return this.resp
}

// GetBody parse response body
func (this *Response) GetContents() (bodyStr string, err error) {
    defer func() {
        _ = this.resp.Body.Close()
    }()

    body, err := io.ReadAll(this.resp.Body)
    if err != nil {
        return "", err
    }

    temp := strings.ReplaceAll(fmt.Sprintf("%v", this.resp.Header["Content-Type"]), " ", "")

    // utf 系列直接返回
    if strings.Contains(strings.ToLower(temp), "charset=utf") {
        bodyStr = string(body)

        // gb 系列当做简体中文处理
    } else if strings.Contains(strings.ToLower(temp), "gb") {
        bodyStr = mahonia.NewDecoder("GB18030").ConvertString(string(body))
    } else {
        //程序没有从对方响应 Header["Content-Type"] 检测到编码类型，那么需要请求者手动设置对方的站点编码
        if decoder := mahonia.NewDecoder(this.charset); decoder != nil {
            bodyStr = decoder.ConvertString(string(body))
        } else {
            return "", errors.New(charsetDecoderError)
        }

    }

    return bodyStr, nil
}

// Get Response ContentLength
func (this *Response) GetContentLength() int64 {
    return this.resp.ContentLength
}

// GetBody parse response body
func (this *Response) GetBody() io.ReadCloser {
    //defer this.resp.Body.Close()

    return this.resp.Body
}

// GetStatusCode get response status code
func (this *Response) GetStatusCode() int {
    return this.resp.StatusCode
}

// GetReasonPhrase get response reason phrase
func (this *Response) GetReasonPhrase() string {
    status := this.resp.Status
    arr := strings.Split(status, " ")

    return arr[1]
}

// IsTimeout get if request is timeout
func (this *Response) IsTimeout() bool {
    if this.err == nil {
        return false
    }
    netErr, ok := this.err.(net.Error)
    if !ok {
        return false
    }
    if netErr.Timeout() {
        return true
    }

    return false
}

// GetHeaders get response headers
func (this *Response) GetHeaders() map[string][]string {
    return this.resp.Header
}

// HasHeader get if header exsits in response headers
func (this *Response) HasHeader(name string) bool {
    headers := this.GetHeaders()
    for k := range headers {
        if strings.ToLower(name) == strings.ToLower(k) {
            return true
        }
    }

    return false
}

// GetHeader get response header
func (this *Response) GetHeader(name string) []string {
    headers := this.GetHeaders()
    for k, v := range headers {
        if strings.ToLower(name) == strings.ToLower(k) {
            return v
        }
    }

    return nil
}

// GetHeaderLine get a single response header
func (this *Response) GetHeaderLine(name string) string {
    header := this.GetHeader(name)
    if len(header) > 0 {
        return header[0]
    }

    return ""
}
