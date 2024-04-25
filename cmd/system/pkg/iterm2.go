package pkg

import "log"

type ITerm2 struct {
}

func (v *ITerm2) Install(value *Package) bool {
	tmp := "/tmp/iterm2.zip"

	url := value.Source[0]
	err := ExecuteCommand(value.Name, "curl", []string{"-fsSL", "-o", tmp, url}, false)
	if err != nil {
		log.Printf("[%s]下载文件异常", value.Name)
		return false
	}

	err = ExecuteCommand(value.Name, "unzip", []string{tmp, "-d", value.Target}, false)
	if err != nil {
		log.Printf("[%s]解药到%s失败", value.Name, err.Error())
		return false
	}
	return true
}

func (v *ITerm2) Upgrade(value *Package) bool {
	return true
}

func (v *ITerm2) Before(value *Package, overwrite bool) bool {
	has, e := ExistApplications(value.Name, value.Bin)
	if e != nil {
		log.Printf("[%s]检查应用异常: %s", value.Name, e.Error())
	}
	// 不存在需要执行安装
	return !has
}

func (v *ITerm2) After(value *Package) {

}

func (v *ITerm2) GetPackage() *Package {
	return &Package{
		Name:        "iTerm2",
		Bin:         "iTerm.app",
		Version:     "latest",
		Compress:    "zip",
		Target:      "/Applications",
		Description: "iTerm2是一个开源的 macOS 终端模拟器。",
		Source:      []string{"https://iterm2.com/downloads/stable/latest"},
	}
}
