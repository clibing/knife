package pkg

type Env struct {
	Key       string `json:"key,omitempty" yaml:"key,omitempty"`               // 环境变量名字
	Value     string `json:"value,omitempty" yaml:"value,omitempty"`           // 环境变量值
	AppendKey string `json:"append_key,omitempty" yaml:"append_key,omitempty"` // 追加到指定变量
}

type Package struct {
	Name        string   `json:"name,omitempty" yaml:"name,omitempty"`               // 包名字
	Bin         string   `json:"bin,omitempty" yaml:"bin,omitempty"`                 // 安装后的名字
	Version     string   `json:"version,omitempty" yaml:"version,omitempty"`         // 安装包版本
	Env         []*Env   `json:"env,omitempty" yaml:"env,omitempty"`                 // 环境变量
	Shell       string   `json:"shell,omitempty" yaml:"shell,omitempty"`             // 安装脚本
	Compress    string   `json:"compress,omitempty" yaml:"compress,omitempty"`       // zip, tar.gz, tar, tar.bz2, rar, xz, tar.xz, git, dmg, txt
	Target      string   `json:"target,omitempty" yaml:"target,omitempty"`           // target 默认为 /Applications
	Description string   `json:"description,omitempty" yaml:"description,omitempty"` // 描述信息
	Source      []string `json:"source,omitempty" yaml:"source,omitempty"`           // 源文件 文件下载地址 | git clone 地址 | 备用 git clone地址
}

/**
 * 定义应用
 */
type Application interface {
	/**
	 * 获取安装
	 */
	GetPackage() *Package

	/**
	 * 安装前事件
	 * false: 不会执行后续事件
	 */
	Before(value *Package, overwrite bool) bool

	/**
	 * 安装应用
	 */
	Install(value *Package) bool

	/**
	 * 升级应用
	 */
	Upgrade(value *Package) bool

	/**
	 * 安装后事件
	 */
	After(value *Package)
}
