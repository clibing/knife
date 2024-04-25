package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"
)

func exitStatus(state *os.ProcessState) int {
	status, ok := state.Sys().(syscall.WaitStatus)
	if !ok {
		return -1
	}
	return status.ExitStatus()
}

func processState(e error) *os.ProcessState {
	err, ok := e.(*exec.ExitError)
	if !ok {
		return nil
	}
	return err.ProcessState
}

func CheckCommand(prefix, value string) (has bool, err error) {
	err = ExecuteCommand(prefix, "which", []string{"-s", value})
	if err != nil {
		return
	}
	has = true
	return
}

func ExecuteCommand(prefix, bin string, value []string) (err error) {
	cmd := exec.Command(bin, value...)
	stdOut, err := cmd.StdoutPipe()
	cmd.Stdin = os.Stdin
	cmd.Stderr = cmd.Stdout
	go func(out io.ReadCloser) {
		for {
			tmp := make([]byte, 1024)
			n, err := out.Read(tmp)
			log.Printf("[%s]%s", prefix, string(tmp[:n]))
			if err != nil {
				break
			}
		}
	}(stdOut)

	if err = cmd.Start(); err != nil {
		return
	}

	if err = cmd.Wait(); err != nil {
		return
	}
	defer stdOut.Close()
	return
}

func GetCurrentProfile(prefix string) (profile string, err error) {
	var homeDir string
	homeDir, err = GetHomeDir(prefix)
	if err != nil {
		return
	}

	sh := os.Getenv("SHELL")
	if strings.HasSuffix(sh, "bash") {
		profile, err = fmt.Sprintf("%s/.bashrc", homeDir), nil
		return
	} else if strings.HasSuffix(sh, "zsh") {
		profile, err = fmt.Sprintf("%s/.zshrc", homeDir), nil
		return
	}
	err = fmt.Errorf("暂不支持当前shell: %s", sh)
	return
}

func GetHomeDir(prefix string) (homeDir string, e error) {
	var u *user.User
	u, e = user.Current()
	if e != nil {
		log.Printf("[%s]获取当前用户名异常, %s", prefix, e)
	}
	homeDir = u.HomeDir
	return
}

/**
 * 检查 env, 并生成待追加的环境变量
 */
func FilterAppendEnv(prefix string, envs []*Env) (export []string) {
	export = make([]string, 0)
	for _, env := range envs {
		current := os.Getenv(env.Key)
		if len(current) > 0 {
			if strings.ContainsAny(current, env.Value) {
				log.Printf("[%s]已经存在配置, current: %s, check: %s\n", prefix, current, env.Value)
				continue
			}
		}
		m := fmt.Sprintf("export %s=%s", env.Key, env.Value)
		export = append(export, m)
		log.Printf("[%s]export %s=%s\n", prefix, env.Key, env.Value)
		if len(env.AppendKey) > 0 {
			existValue := os.Getenv(env.AppendKey)
			latestValue := fmt.Sprintf("export %s=%s:$%s", env.AppendKey, existValue, env.Key)
			export = append(export, latestValue)
		}
	}
	return
}

/**
 * prefix: 日志的前缀， 建议 Package.Name
 * value: export key=value;
 * target: 将export配置文件写入到目标文件
 */
func SettingEnv(prefix string, value []string, target string) (result bool, err error) {
	if len(value) == 0 {
		log.Printf("[%s]不需要追加配置", prefix)
		result = true
		return
	}
	var pd *os.File
	if len(target) > 0 {
		pd, err = os.OpenFile(target, os.O_RDWR|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Printf("[%s]打开配置文件异常: %s", prefix, err.Error())
			return
		}
		defer pd.Close()
		pd.Seek(0, io.SeekEnd)

		for _, line := range value {
			log.Printf("[%s]追加配置 %s", prefix, line)
			_, e := pd.WriteString(line + "\n")
			if e != nil {
				log.Printf("[%s]写入配置异常 %s", prefix, e)
				break
			}
		}
		result = true
	} else {
		result = false
		err = fmt.Errorf("[%s]写入配置文件不存在", prefix)
	}
	return
}

/**
 * 检查目录是否存在
 */
func ExistPath(prefix, base, value string) (has bool, err error) {
	_, err = os.Stat(fmt.Sprintf("%s/%s", base, value))
	if os.IsNotExist(err) {
		has = false
		return
	}
	has = true
	return
}
func ExistApplications(prefix, value string) (has bool, err error) {
	has, err = ExistPath(prefix, "/Applications", value)
	return
}
