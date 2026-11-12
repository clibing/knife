package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

/**
 * 执行系统git命令
 * 安装hook
 */
func ExecGitV2(commands ...string) (err error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("git.exe", commands...)
	} else {
		cmd = exec.Command("git", commands...)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err = cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Error starting command:", err)
		return
	}

	if err = cmd.Wait(); err != nil {
		return errors.New(strings.TrimSpace(err.Error()))
	}
	return nil
}
