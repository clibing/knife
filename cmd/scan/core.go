/**
 * 扫描
 * git项目并记录到本地 ~/.knife/scan_record.db
 * 1. 整理目录
 * 2. 根据项目结构，清理编译输出。
 */
package scan

import (
	"context"
	"io/fs"
)

type Meta struct {
}

/**
 * 调度器
 */
type Scheduler interface {
	/**
	 * 启动
	 */
	Start(ctx context.Context, path string) (err error)

	/**
	 * 获取持久化处理
	 */
	GetPipeline() (values []Pipline, err error)

	/**
	 * 分析器处理器
	 */
	GetProcesser() (values []Processer, err error)
}

/**
 * 分析后置事件
 * 持久化 处理器
 */
type Pipline interface {
	/**
	 * 持久化处理
	 * 1. 是否保存
	 */
	Pipline(meta *Meta) (err error)
}

/**
 * 对不同目录进行不同的分析器
 * - 1. 是否为Github
 * - 2. 是否为Java Maven项目
 * - 3. ESP32, ESP8266
 * - 4. Rust项目
 */
type Processer interface {
	/**
	 * 分析处理
	 */
	Processer(path string, info fs.FileInfo) (meta *Meta, ok bool, err error)
}
