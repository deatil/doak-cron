Go交叉编译(Go语言Mac/Linux/Windows下交叉编译)

~~~cmd
GOOS：目标平台的操作系统（darwin、freebsd、linux、windows）
GOARCH：目标平台的体系架构（386、amd64、arm）
交叉编译不支持 CGO 所以要禁用它
~~~

## 1.Mac下编译Linux, Windows

### Linux
~~~cmd
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
~~~

### Windows
~~~cmd
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
~~~
如: `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o main-windows main.go`

## 2.Linux下编译Mac, Windows

### Mac
~~~cmd
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
~~~

### Windows
~~~cmd
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
~~~

## 3.Windows下编译Mac, Linux

### Mac
~~~cmd
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build main.go
~~~

### Linux
~~~cmd
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go
~~~

### 编译为Window可运行文件
~~~cmd
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build main.go
~~~

## 4.参数说明

查看环境：

~~~cmd
$> go env
GO111MODULE=""
GOARCH="amd64"
GOBIN=""
GOCACHE="/Users/usr/Library/Caches/go-build"
GOENV="/Users/usr/Library/Application Support/go/env"
GOEXE=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="darwin"
GONOPROXY=""
GONOSUMDB=""
GOOS="darwin"
GOPATH="/go"
GOPRIVATE=""
GOPROXY="https://proxy.golang.org,direct"
GOROOT="/usr/local/go"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/usr/local/go/pkg/tool/darwin_amd64"
GCCGO="gccgo"
AR="ar"
CC="clang"
CXX="clang++"
CGO_ENABLED="1"
GOMOD=""
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS=""
~~~
