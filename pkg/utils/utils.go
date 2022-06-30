package utils

import (
    "io"
    "os"
)

// 文件是否存在
func FileExists(path string) bool {
    _, err := os.Stat(path)

    return err == nil || os.IsExist(err)
}

// 文件删除
func FileDelete(path string) error {
    return os.Remove(path)
}

// 获取数据
func FileRead(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer file.Close()

    data, err2 := io.ReadAll(file)
    if err2 != nil {
        return "", err2
    }

    return string(data), nil
}
