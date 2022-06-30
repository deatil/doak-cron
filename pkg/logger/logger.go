package logger

import (
    "io"
    "os"
    "sync"

    "github.com/rs/zerolog"
)

var (
    log *zerolog.Logger

    once sync.Once
)

// LevelFileWriter 指定级别的日志写文件
type LevelFileWriter struct {
    lw io.Writer
    lv zerolog.Level
}

func (w *LevelFileWriter) Write(p []byte) (n int, err error) {
    return w.lw.Write(p)
}

func (w *LevelFileWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
    if l == w.lv {
        return w.lw.Write(p)
    }

    return len(p), nil
}

// LevelConsoleWriter 特定级别日志写到控制台
type LevelConsoleWriter struct {
    lw zerolog.ConsoleWriter
    lv []zerolog.Level
}

func (w *LevelConsoleWriter) Write(p []byte) (n int, err error) {
    return w.lw.Write(p)
}

func (w *LevelConsoleWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
    for _, v := range w.lv {
        if v == l {
            return w.lw.Write(p)
        }
    }

    return len(p), nil
}

// Log().Trace().Msg("test TRACE")
// Log().Debug().Msg("test DEBUG")
// Log().Info().Msg("test INFO")
// Log().Warn().Msg("test WARN")
// Log().Error().Msg("test ERROR")
// Log().Fatal().Msg("test FATAL")
func Log() *zerolog.Logger {
    once.Do(func() {
        file := "./cron.log"

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
