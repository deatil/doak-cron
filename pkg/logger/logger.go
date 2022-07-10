package logger

import (
    "os"
    "sync"

    "github.com/rs/zerolog"
)

var (
    log  *zerolog.Logger
    once sync.Once

    mu   sync.RWMutex
    file string = "./cron.log"
)

// 设置日志存储文件
func WithLogFile(f string) {
    mu.Lock()
    defer mu.Unlock()

    file = f
}

// 获取日志存储文件
func GetLogFile() string {
    mu.RLock()
    defer mu.RUnlock()

    return file
}

// Log().Trace().Msg("test TRACE")
// Log().Debug().Msg("test DEBUG")
// Log().Info().Msg("test INFO")
// Log().Warn().Msg("test WARN")
// Log().Error().Msg("test ERROR")
// Log().Fatal().Msg("test FATAL")
func Log() *zerolog.Logger {
    once.Do(func() {
        file := GetLogFile()

        log = Manager(file)
    })

    return log
}

// 日志管理
func Manager(file string) *zerolog.Logger {
    errorFile, _ := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

    newLog := zerolog.New(zerolog.MultiLevelWriter(
            // Trace 日志写入 cron.log
            &LevelFileWriter{
                lw: errorFile,
                lv: zerolog.TraceLevel,
            },
            // Debug 日志写入 cron.log
            &LevelFileWriter{
                lw: errorFile,
                lv: zerolog.DebugLevel,
            },
            // Info 日志写入 cron.log
            &LevelFileWriter{
                lw: errorFile,
                lv: zerolog.InfoLevel,
            },
            // Warn 日志写入 cron.log
            &LevelFileWriter{
                lw: errorFile,
                lv: zerolog.WarnLevel,
            },
            // Error 日志写入 cron.log
            &LevelFileWriter{
                lw: errorFile,
                lv: zerolog.ErrorLevel,
            },
            // Fatal 日志写入 cron.log
            &LevelFileWriter{
                lw: errorFile,
                lv: zerolog.FatalLevel,
            },
            // Panic 日志写入 cron.log
            &LevelFileWriter{
                lw: errorFile,
                lv: zerolog.PanicLevel,
            },
            // Debug, Fatal 日志显示在控制台
            &LevelConsoleWriter{
                lw: zerolog.ConsoleWriter{Out: os.Stdout},
                lv: []zerolog.Level{
                    zerolog.DebugLevel,
                },
            },
        )).
        With().
        Timestamp().
        Caller().
        Logger()

    return &newLog
}
