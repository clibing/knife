package pkg

import "log"

type Brew struct{}

/**
 * 安装应用
 */
func (v *Brew) Install(value *Package) bool {
	var err error
	log.Printf("[%s]下载homebrew安装包\n", value.Name)
	err = ExecuteCommand(value.Name, "git", []string{
		"clone",
		"--depth=1",
		"https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/install.git",
		"/tmp/brew-install",
	}, false)

	if err != nil {
		log.Printf("[%s]克隆homebrew镜像异常: %s\n", value.Name, err.Error())
		return false
	}

	log.Printf("[%s]授权homebrew可执行脚本\n", value.Name)
	// 授权
	ExecuteCommand(value.Name, "chmod", []string{"+x", "/tmp/brew-install/install.sh"}, false)

	log.Printf("[%s]安装homebrew\n", value.Name)
	err = ExecuteCommand(value.Name, "/bin/bash", []string{"/tmp/brew-install/install.sh"}, false)
	if err != nil {
		log.Printf("[%s]安装homebrew异常: %s\n", value.Name, err.Error())
		return false
	}

	export := FilterAppendEnv(value.Name, value.Env)
	result, e := SettingEnv(value.Name, export, value.Target)
	if e != nil {
		log.Printf("[%s]\n", value.Name)
	}
	return result
}

/**
 * 升级应用
 */
func (v *Brew) Upgrade(value *Package) bool {
	return true
}

func (v *Brew) Before(value *Package, overwrite bool) bool {
	has, _ := CheckCommand(value.Name, "git")
	if !has {
		log.Printf("[%s]需要先安装git\n", value.Name)
		return false
	}

	if overwrite {
		ExecuteCommand(value.Name, "rm", []string{"-rf", "/tmp/brew-install"}, false)
		log.Printf("[%s]强制安装\n", value.Name)
		return true
	}

	has, e := CheckCommand(value.Name, value.Bin)
	log.Printf("[%s]检查当前命令是否存在: %t, err: %s\n", value.Name, has, e)
	return !has
}

func (v *Brew) After(value *Package) {
	e := ExecuteCommand(value.Name, "rm", []string{"-rf", "/tmp/brew-install"}, false)
	if e != nil {
		log.Printf("[%s]删除brew-install目录失败\n", value.Name)
	}

	// ExecuteCommand(value.Name, "brew", []string{"update", "--verbose"}, false)

	// ExecuteCommand(value.Name, "brew", []string{"upgrade", "--verbose"}, false)
	ExecuteCommand(value.Name, "brew", []string{"help"}, false)
}

func (v *Brew) GetPackage() *Package {
	return &Package{
		Name:    v.Key(),
		Bin:     "brew",
		Version: "latest",
		Env: []*Env{
			{
				Key:   "HOMEBREW_API_DOMAIN",
				Value: "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles/api",
			},
			{
				Key:   "HOMEBREW_BOTTLE_DOMAIN",
				Value: "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles",
			},
			{
				Key:   "HOMEBREW_BREW_GIT_REMOTE",
				Value: "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/brew.git",
			},
			{
				Key:   "HOMEBREW_CORE_GIT_REMOTE",
				Value: "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-core.git",
			},
			{
				Key:   "HOMEBREW_PIP_INDEX_URL",
				Value: "https://pypi.tuna.tsinghua.edu.cn/simple",
			},
		},
	}
}

func (v *Brew) Key() string {
	return "homebrew"
}
