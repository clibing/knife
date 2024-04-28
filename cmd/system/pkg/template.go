package pkg

type Template struct {
	Log
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

func (v *Template) Key() string {
	return "template"
}

func NewTemplate() *Template {
	tmp := &Template{}
	l := Log{Key: tmp.Key()}
	tmp.Log = l
	return tmp
}
