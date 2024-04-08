package utils

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

/**
 * 执行系统git命令
 * 安装hook
 */
func ExecGit(commands ...string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("git.exe", commands...)
	} else {
		cmd = exec.Command("git", commands...)
	}

	bs, err := cmd.CombinedOutput()
	if err != nil {
		if bs != nil {
			return "", errors.New(strings.TrimSpace(string(bs)))
		}
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}

