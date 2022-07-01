package logger

import (
    "io"

    "github.com/rs/zerolog"
)

// 指定级别的日志写文件
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
