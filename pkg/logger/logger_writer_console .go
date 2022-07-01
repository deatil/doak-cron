package logger

import (
    "github.com/rs/zerolog"
)

// 特定级别日志写到控制台
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
