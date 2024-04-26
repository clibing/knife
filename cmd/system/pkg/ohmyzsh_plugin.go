package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type OhmyzshPlugin struct {
	Plugin []*Plugin
	Log
}

type Plugin struct {
	Name   string   // 名字
	Source string   // 源码
	Target string   // 目标
	Shell  []string // 自定义 shell
}

/**
 * 安装应用
 */
func (v *OhmyzshPlugin) Install(value *Package) bool {
	var err error
	for _, plugin := range v.Plugin {
		log.Printf("[%s]安装插件:%s\n", value.Name, plugin.Name)
		err = ExecuteCommand(value.Name, "git", []string{"clone", plugin.Source, plugin.Target}, false)
		if err != nil {
			log.Printf("[%s]安装插件[%s]失败%s\n", value.Name, plugin.Name, err)
			return false
		}
	}

	profile, e := GetCurrentProfile(value.Name)
	if e != nil {
		log.Printf("[%s]读取配置文件异常%s\n", value.Name, e)
	}

	log.Printf("[%s]配置文件\n", value.Name)
	source, err := os.Open(profile) // 打开文件
	if err != nil {
		log.Printf("[%s]打开配置文件异常%s\n", value.Name, err)
		return false
	}
	defer source.Close() // 确保文件在函数结束时关闭
	// 创建新的Scanner
	scanner := bufio.NewScanner(source)

	newFile := fmt.Sprintf("%s.%d", profile, time.Now().Unix())
	target, err := os.Create(newFile)
	if err != nil {
		log.Printf("[%s]创建临时配置文件[%s], 发生异常%s\n", value.Name, newFile, err)
		return false
	}
	defer target.Close() // 升级配置文件 结束时关闭

	// 创建新的Writer
	writer := bufio.NewWriter(target)

	// 逐行扫描
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "plugins=") {
			plugins := v.GetPlugin()
			plugins = append(plugins, "git", "docker", "last-working-dir")
			line = fmt.Sprintf("plugins=(%s)", strings.Join(plugins, " "))
		}
		writer.WriteString(line + "\n")
	}
	writer.Flush()
	target.Close()
	source.Close()

	if err := scanner.Err(); err != nil {
		log.Printf("[%s]逐行扫描配置文件异常%s\n", value.Name, err)
		return false
	}
	backupFile := fmt.Sprintf("%s.backup.%d", profile, time.Now().Unix())
	log.Printf("[%s]备份配置文件%s\n", value.Name, backupFile)
	ExecuteCommand(value.Name, "mv", []string{profile, backupFile}, false)
	ExecuteCommand(value.Name, "mv", []string{newFile, profile}, false)

	return true
}

/**
 * 升级应用
 */
func (v *OhmyzshPlugin) Upgrade(value *Package) bool {
	ExecuteCommand(value.Name, "omz", []string{"update"}, false)
	return true
}

/**
 * 之后事件
 */
func (v *OhmyzshPlugin) After(value *Package) {
}

/**
 * 获取插件 name []string
 */
func (o *OhmyzshPlugin) GetPlugin() (key []string) {
	for _, v := range o.Plugin {
		key = append(key, v.Name)
	}
	return
}

func (v *OhmyzshPlugin) Before(value *Package, overwrite bool) bool {
	homeDir, _ := GetHomeDir(value.Name)
	has, _ := ExistPath(value.Name, homeDir, ".oh-my-zsh")
	if overwrite {
		log.Printf("[%s]强制安装\n", value.Name)
		ExecuteCommand(value.Name, "rm", []string{"-rf", fmt.Sprintf("%s/%s", homeDir, ".oh-my-zsh/custom/plugins/zsh-syntax-highlighting")}, false)
		ExecuteCommand(value.Name, "rm", []string{"-rf", fmt.Sprintf("%s/%s", homeDir, ".oh-my-zsh/custom/plugins/zsh-autosuggestions")}, false)
		return true
	}
	return has
}

func (v *OhmyzshPlugin) GetPackage() *Package {
	return &Package{
		Name:        v.Key(),
		Version:     "latest",
		Env:         []*Env{},
		Compress:    "git",
		Target:      "",
		Description: "ohmyzsh源码安装",
	}
}

func (v *OhmyzshPlugin) Key() string {
	return "ohmyzsh-plugin"
}

func NewOhmyzshPlugin() *OhmyzshPlugin {
	homeDir, _ := GetHomeDir("ohmyzsh")
	o := &OhmyzshPlugin{
		Plugin: []*Plugin{
			{
				Name:   "zsh-syntax-highlighting",
				Source: "https://github.com/zsh-users/zsh-syntax-highlighting.git",
				Target: fmt.Sprintf("%s/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting", homeDir),
			},
			{

				Name:   "zsh-autosuggestions",
				Source: "https://github.com/zsh-users/zsh-autosuggestions",
				Target: fmt.Sprintf("%s/.oh-my-zsh/custom/plugins/zsh-autosuggestions", homeDir),
			},
		},
	}
	o.Log = Log{Key: o.Key()}
	return o
}
