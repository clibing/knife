package maven

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/clibing/knife/cmd/debug"
	"github.com/clibing/knife/internal/model"
)

// parse name
const REG = `(?P<Prefix>[\w\-\._]+)\-(?P<Date>[0-9]{8})\.?(?P<Time>[0-9]{6})\-?(?P<BuildNumber>[0-9]+)\-?(sources)?\.(?P<Extension>(jar|pom|war))\.?(?P<Signature>sha1)?`
const SEP_KEY = "/"
const MAVEN_METADATA_LOCAL_NAME = "maven-metadata-local.xml"
const MAVEN_METADATA_LOCAL_NAME_2 = "maven-metadata-evernote-mirror.xml"

// 核心入口处理文件
func Doing(debug *debug.Debug, path string) (data map[string]*model.Gav) {
	// scan(scanPath)
	result, metadata := scanning(path)

	gavList := make([]*model.Info, 0)

	for _, value := range result {
		ok, gav := parse(value)
		if ok {
			m, ok := metadata[filepath.Dir(value)]
			if ok {
				gav.GroupId = m.GroupId.Text
				gav.ArticleId = m.ArtifactId.Text
				gav.Version = m.Version.Text
			}
			gavList = append(gavList, &gav)
		}
	}

	// init map
	data = make(map[string]*model.Gav)
	for _, item := range gavList {
		key := item.Key()
		if key == "" {
			continue
		}

		mgav, ok := data[key]
		// new
		if !ok {
			mgav = &model.Gav{
				Release:  make([]*model.Info, 0),
				Snapshot: make([]*model.Info, 0),
			}
			if item.Basic {
				mgav.Release = append(mgav.Release, item)
			} else {
				mgav.Snapshot = append(mgav.Snapshot, item)
			}
			data[key] = mgav
			continue
		}
		// exist
		if item.Basic {
			mgav.Release = append(mgav.Release, item)
			continue
		}
		debug.Debug("key: %s ;snapshot version: %s", key, item.SimpleName)
		mgav.Snapshot = append(mgav.Snapshot, item)
		// 多版本处理
		if len(mgav.Snapshot) <= 3 {
			continue
		}

	}

	for _, mgav := range data {
		sort.Slice(mgav.Snapshot, func(i, j int) bool {
			s1b, _ := strconv.Atoi(mgav.Snapshot[i].BuildNumber)
			s1 := fmt.Sprintf("%s.%s-%2d", mgav.Snapshot[i].Date, mgav.Snapshot[i].Time, s1b)
			s2b, _ := strconv.Atoi(mgav.Snapshot[j].BuildNumber)
			s2 := fmt.Sprintf("%s.%s-%2d", mgav.Snapshot[j].Date, mgav.Snapshot[j].Time, s2b)
			return strings.Compare(s1, s2) == 1
		})

		currentTag := 0
		for _, mv := range mgav.Snapshot {
			number, _ := strconv.Atoi(mv.BuildNumber)
			if currentTag == 0 {
				currentTag = number
				continue
			}
			if number < currentTag {
				mv.Deleted = true
			}
		}
	}
	return
}

/**
 * 扫描指定目录
 */
func scanning(path string) (result []string, metadata map[string]model.Metadata) {
	result = make([]string, 0)
	metadata = make(map[string]model.Metadata)
	// 提取完整路径
	scanDir, _ := filepath.Abs(path)
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
				ok, key, meta := scanMetadataLocal(path)
				if ok {
					metadata[key] = meta
				}
			}
		}
		return nil
	})
	return
}

func scanMetadataLocal(name string) (status bool, dir string, metadata model.Metadata) {
	simple := filepath.Base(name)
	// 跳过 不是 以 metadata prefix file
	if !strings.HasPrefix(simple, "maven-metadata") {
		return
	}
	// 不是 xml文件
	if strings.LastIndex(simple, ".xml") == -1 {
		return
	}

	v, err := os.ReadFile(name)
	if err != nil {
		return
	}

	xml.Unmarshal(v, &metadata)
	dir = filepath.Dir(name)
	status = true
	return
}

/**
 * 格式化文件名字为maven
 * Parse full name
 */
func parse(fullName string) (status bool, gav model.Info) {
	// 是否为快照
	snapshot := strings.Contains(fullName, "-SNAPSHOT")

	var compRegEx = regexp.MustCompile(REG)
	match := compRegEx.FindStringSubmatch(fullName)

	arr := strings.Split(fullName, SEP_KEY)
	size := len(arr)

	simepleName := arr[size-1]

	gav.FullName = fullName
	gav.Snapshot = snapshot
	gav.Classifier = false
	gav.SimpleName = simepleName

	if !gav.IsJarOrPomOrWar() {
		status = false
		return
	}

	status = true
	if strings.LastIndex(fullName, "-SNAPSHOT-sources.jar") != -1 {
		gav.Basic = true
		gav.Classifier = true
		gav.Extension = "jar"
		return
	} else if strings.LastIndex(fullName, "-SNAPSHOT.jar") != -1 {
		gav.Basic = true
		gav.Extension = "jar"
		return
	} else if strings.LastIndex(fullName, "-SNAPSHOT.pom") != -1 {
		gav.Basic = true
		gav.Extension = "pom"
		return
	} else if strings.LastIndex(fullName, "-SNAPSHOT.war") != -1 {
		gav.Basic = true
		gav.Extension = "war"
		return
	}

	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			if name == "Date" {
				gav.Date = match[i]
			} else if name == "Time" {
				gav.Time = match[i]
			} else if name == "BuildNumber" {
				gav.BuildNumber = match[i]
				gav.Basic = false
			} else if name == "Extension" {
				gav.Extension = match[i]
			} else if name == "Signature" {
				gav.Signature = true
			}
		}
	}

	if strings.LastIndex(fullName, "-sources.jar") != -1 && gav.Signature {
		gav.Classifier = true
	}

	return
}
