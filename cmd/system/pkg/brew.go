package pkg

import "log"

type Brew struct{}

/**
 * 安装应用
 */
func (v *Brew) Install(value *Package) bool {

	err := ExecuteCommand(value.Name, "git", []string{
		"clone",
		"--depth=1",
		"https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/install.git",
		"/tmp/brew-install",
	}, false)
	if err != nil {
		log.Printf("[%s]克隆homebrew镜像异常: %s", value.Name, err.Error())
		return false
	}
	err = ExecuteCommand(value.Name, "/bin/bash", []string{"brew-install/install.sh"}, false)
	if err != nil {
		log.Printf("[%s]安装homebrew异常: %s", value.Name, err.Error())
		return false
	}
	export := FilterAppendEnv(value.Name, value.Env)
	result, e := SettingEnv(value.Name, export, value.Target)
	if e != nil {
		log.Printf("[%s]", value.Name)
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
		log.Printf("[%s]需要先安装git", value.Name)
		return false
	}

	has, e := CheckCommand(value.Name, value.Bin)
	if e != nil {
		return false
	}
	return has
}

func (v *Brew) After(value *Package) {
	e := ExecuteCommand(value.Name, "rm", []string{"-rf", "brew-install"}, false)
	if e != nil {
		log.Printf("[%s]删除brew-install目录失败", value.Name)
	}
}

func (v *Brew) GetPackage() *Package {
	return &Package{
		Name:    "Homebrew",
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
		Shell:       "",
		Compress:    "",
		Target:      "",
		Description: "",
		Source:      []string{},
	}
}
func InstallByBrew(prefix, cmd string) (err error) {
	if len(prefix) == 0 {
		prefix = "brew"
	}
	err = ExecuteCommand(prefix, "brew", []string{"install", cmd}, false)
	return
}
