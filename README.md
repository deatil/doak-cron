## doak-cron

使用 go 实现的计划任务管理，支持脚本和 `curl` 请求两种方式


### 项目介绍

*  使用 go 实现的计划任务管理
*  最小设置单位为秒


### 使用方法

1. 构建对应系统的文件。交叉编译可查看文档 `/docs/go-build.md`

```go
go build main.go
```

2. 使用

执行计划任务。加 `--debug` 会记录返回的数据
```go
main.exe cron --conf="./cron.json" --debug
```

查看当前版本号
```go
main.exe cron version
```

3. 配置使用

`cron.json` 为计划任务配置文件。当前支持脚本和 `curl` 两种方式


### 使用预览

![doak-cron](https://user-images.githubusercontent.com/24578855/178781346-af72bea7-3210-4138-840c-3138408147ef.jpg)


### 特别鸣谢

感谢以下的项目,排名不分先后

 - github.com/urfave/cli

 - github.com/robfig/cron

 - github.com/go-resty/resty

 - github.com/rs/zerolog

 - github.com/spf13/cast

 - github.com/jedib0t/go-pretty


### 开源协议

*  `doak-cron` 遵循 `Apache2` 开源协议发布，在保留本系统版权的情况下提供个人及商业免费使用。


### 版权

*  该系统所属版权归 deatil(https://github.com/deatil) 所有。
