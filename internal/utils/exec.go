package utils

import (
	"bufio"
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
func ExecGit(commands ...string) (err error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("git.exe", commands...)
	} else {
		cmd = exec.Command("git", commands...)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating stdout pipe:", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating stderr pipe:", err)
		return
	}

	if err = cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Error starting command:", err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintln(os.Stderr, scanner.Text())
		}
	}()

	if err = cmd.Wait(); err != nil {
		return errors.New(strings.TrimSpace(err.Error()))
	}
	return nil
}
