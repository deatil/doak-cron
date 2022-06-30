package cmd

import (
    "os"
    "os/exec"
    "bytes"
)

// 脚本文件
func Command(commandName string, params ...string) (string, error) {
    cmd := exec.Command(commandName, params...)

    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr

    // 重定目录
    // cmd.Dir = commandDir

    err := cmd.Start()
    if err != nil {
        return "", err
    }

    err = cmd.Wait()

    return out.String(), err
}

