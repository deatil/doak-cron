package cron

import (
    "github.com/deatil/doak-cron/pkg/logger"
)

// 构造函数
func NewLogger() Logger {
    return Logger{}
}

// 日志文件地址
var logFile string = "./cron-error.log"

/**
 * 日志
 *
 * @create 2022-12-2
 * @author deatil
 */
type Logger struct {}

// 实现接口
func (this Logger) Printf(format string, v ...any) {
    format = "[cron-recover] " + format

    logger.Manager(logFile).Error().Msgf(format, v...)
}
