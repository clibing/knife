package pkg

type Template struct {
}

func (v *Template) Install(value *Package) bool {
	return true
}

func (v *Template) Upgrade(value *Package) bool {
	return true
}

func (v *Template) Before(value *Package, overwrite bool) bool {
	return true
}

func (v *Template) After(value *Package) {
}

func (v *Template) GetPackage() *Package {
	return &Package{}
}
