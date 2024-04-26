package pkg

type Software struct {
}

func (v *Software) Install(value *Package) bool {
	for _, sf := range value.Source {
		ExecuteCommand(value.Name, "brew", []string{"install", sf}, false)
	}
	return true
}

func (v *Software) Upgrade(value *Package) bool {
	return true
}

func (v *Software) Before(value *Package, overwrite bool) bool {
	return true
}

func (v *Software) After(value *Package) {
}

func (v *Software) GetPackage() *Package {
	return &Package{
		Name:        v.Key(),
		Bin:         "",
		Version:     "latest",
		Env:         []*Env{},
		Shell:       "",
		Compress:    "",
		Target:      "",
		Description: "安装一些常用软件",
		Source: []string{
			"hugo",
			"wget",
			"lux",
			"dive",
			"ctop",
			"restic",
			"noti",
			"ffmpeg",
			"git-chglog",
		},
	}
}

func (v *Software) Key() string {
	return "Software"
}
