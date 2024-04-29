package scan

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
)

/**
 * 默认调度器
 */
type DefaultScheduler struct {
	Processer []Processer // 分析器
	Pipline   []Pipline   // 保存 保存到sqlite3，log
	Done      chan bool   // signal 等待信号量
}

func (ds *DefaultScheduler) Log(format string, args ...interface{}) {
	fmt.Printf("[Scheduler] "+format+"\n", args...)
}

/**
 * 启动
 */
func (ds *DefaultScheduler) Start(ctx context.Context, path string) (err error) {
	ds.Log("系统运行中")

	// 提取完整路径
	scanDir, _ := filepath.Abs(path)
	ds.Log("scan dir: %s", scanDir)

	// 扫描
	filepath.Walk(scanDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		processers, e := ds.GetProcesser()
		if e == nil {
			for _, p := range processers {
				ds.Log("running processor")
				meta, ok, _ := p.Processer(path, info)
				if ok {
					piplines, _ := ds.GetPipeline()
					for _, p := range piplines {
						p.Pipline(meta)
					}
				}
			}

		}
		return nil
	})
	ds.Done <- true
	return
}

/**
 * 获取持久化处理
 */
func (ds *DefaultScheduler) GetPipeline() (values []Pipline, err error) {
	values = ds.Pipline
	return
}

/**
 * 分析器处理器
 */
func (ds *DefaultScheduler) GetProcesser() (values []Processer, err error) {
	values = ds.Processer
	return
}

/**
 * 绑定后置处理器
 */
func (ds *DefaultScheduler) BindPipline(value ...Pipline) *DefaultScheduler {
	if value != nil {
		ds.Pipline = append(ds.Pipline, value...)
	}
	return ds
}

/**
 * 处理器
 * 对文件分析， 不同的绑定
 */
func (ds *DefaultScheduler) BindProcesser(value ...Processer) *DefaultScheduler {
	if value != nil {
		ds.Processer = append(ds.Processer, value...)
	}
	return ds
}
