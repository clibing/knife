package pkg

type Brew struct {
	Log
}

/**
 * 安装应用
 */
func (v *Brew) Install(value *Package) bool {
	var err error
	v.Log.Println("下载homebrew安装包")

	err = ExecuteCommand(value.Name, "git", []string{
		"clone",
		"--depth=1",
		"https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/install.git",
		"/tmp/brew-install",
	}, false)

	if err != nil {
		v.Log.Println("克隆homebrew镜像异常: %s", err.Error())
		return false
	}

	v.Log.Println("授权homebrew可执行脚本")

	// 授权
	ExecuteCommand(value.Name, "chmod", []string{"+x", "/tmp/brew-install/install.sh"}, false)

	v.Log.Println("安装homebrew")

	err = ExecuteCommand(value.Name, "/bin/bash", []string{"/tmp/brew-install/install.sh"}, false)
	if err != nil {
		v.Log.Println("安装homebrew异常: %s", err.Error())
		return false
	}

	export := FilterAppendEnv(value.Name, value.Env)
	result, e := SettingEnv(value.Name, export, value.Target)
	if e != nil {
		v.Log.Println("设置环境失败")
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
		v.Log.Println("需要先安装git")
		return false
	}

	if overwrite {
		ExecuteCommand(value.Name, "rm", []string{"-rf", "/tmp/brew-install"}, false)
		v.Log.Println("强制安装")
		return true
	}

	has, e := CheckCommand(value.Name, value.Bin)
	v.Log.Println("检查当前命令是否存在: %t, err: %s", has, e)
	return !has
}

func (v *Brew) After(value *Package) {
	e := ExecuteCommand(value.Name, "rm", []string{"-rf", "/tmp/brew-install"}, false)
	if e != nil {
		v.Log.Println("删除brew-install目录失败")
	}

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

func NewBrew() *Brew {
	b := &Brew{}
	log := Log{Key: b.Key()}
	b.Log = log
	return b
}
