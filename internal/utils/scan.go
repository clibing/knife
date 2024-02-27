package utils
	/**
 * 扫描指定目录
 */


func Scan(){

	result = make([]string, 0)
	metadata = make(map[string]model.Metadata)
	// 提取完整路径
	scanDir, _ := filepath.Abs(path)
	debug.ShowSame("🟣 scan path: %s", scanDir)

	// 扫描
	filepath.Walk(scanDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只保存非目录
		if !info.IsDir() {
			result = append(result, path)
			key := filepath.Dir(path)
			_, ok := metadata[key]
			if !ok {
				ok, key, meta := scanMetadataLocal(debug, path)
				if ok {
					metadata[key] = meta
				}
			}
		}
		return nil
	})
	return
}