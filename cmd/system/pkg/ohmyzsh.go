package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Ohmyzsh struct {
	Plugin []*Plugin
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
func (v *Ohmyzsh) Install(value *Package) bool {
	err := ExecuteCommand(value.Name, "sh", value.Source)
	if err != nil {
		log.Printf("[%s]安装失败%s", value.Name, err)
		return false
	}

	for _, plugin := range v.Plugin {
		err = ExecuteCommand(value.Name, "git", []string{"clone", plugin.Source, plugin.Target})
		if err != nil {
			log.Printf("[%s]安装插件[%s]失败%s", value.Name, plugin.Name, err)
			return false
		}
	}

	profile, e := GetCurrentProfile(value.Name)
	if e != nil {
		log.Printf("[%s]读取配置文件异常%s", value.Name, e)
	}

	file, err := os.Open(profile) // 打开文件
	if err != nil {
		log.Printf("[%s]打开配置文件异常%s", value.Name, err)
		return false
	}
	defer file.Close() // 确保文件在函数结束时关闭

	target := fmt.Sprintf("%s.%d", profile, time.Now().Unix())
	upgrade, err := os.Create(target)
	if err != nil {
		log.Printf("[%s]创建临时配置文件[%s], 发生异常%s", value.Name, target, err)
		return false
	}
	defer upgrade.Close() // 升级配置文件 结束时关闭

	// 创建新的Writer
	writer := bufio.NewWriter(upgrade)

	// 创建新的Scanner
	scanner := bufio.NewScanner(file)

	// 逐行扫描
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "plugins=") {
			plugins := v.GetPlugin()
			plugins = append(plugins, "git", "docker", "last-working-dir")
			line = fmt.Sprintf("plugins=(%s)", strings.Join(plugins, " "))
		}
		writer.WriteString(line)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[%s]逐行扫描配置文件异常%s", value.Name, err)
		return false
	}
	return true
}

/**
 * 升级应用
 */
func (v *Ohmyzsh) Upgrade(value *Package) bool {
	return true
}

/**
 * 之后事件
 */
func (v *Ohmyzsh) After(value *Package) {
}

/**
 * 获取插件 name []string
 */
func (o *Ohmyzsh) GetPlugin() (key []string) {
	for _, v := range o.Plugin {
		key = append(key, v.Name)
	}
	return
}

func (v *Ohmyzsh) Before(value *Package, overwrite bool) bool {
	homeDir, _ := GetHomeDir(value.Name)
	has, e := ExistPath(value.Name, homeDir, ".oh-my-zsh")
	if e != nil {
		log.Printf("[%s]安装ohmyzsh检查异常%s", value.Name, e)
	}
	return !has
}

func (v *Ohmyzsh) GetPackage() *Package {
	return &Package{
		Name:        "ohmyzsh",
		Bin:         "ohmyzsh",
		Version:     "latest",
		Env:         []*Env{},
		Shell:       "zsh",
		Compress:    "git",
		Target:      "",
		Description: "ohmyzsh源码安装",
		Source: []string{
			"-c",
			"\"$(curl -fsSL https://install.ohmyz.sh/)\"",
		},
	}
}

func NewOhmyzsh() *Ohmyzsh {
	return &Ohmyzsh{
		Plugin: []*Plugin{
			{
				Name:   "zsh-syntax-highlighting",
				Source: "https://github.com/zsh-users/zsh-syntax-highlighting.git",
				Target: "$ZSH_CUSTOM/plugins/zsh-syntax-highlighting",
			},
			{

				Name:   "zsh-autosuggestions",
				Source: "https://github.com/zsh-users/zsh-autosuggestions",
				Target: "$ZSH_CUSTOM/plugins/zsh-autosuggestions",
			},
		},
	}
}
