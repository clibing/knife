package model

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Info struct {
	FullName    string `json:"full_name"`    // 文件的完整路径+名字
	SimpleName  string `json:"simple_name"`  // 文件的名字
	GroupId     string `json:"group_id"`     // maven group id
	ArticleId   string `json:"article"`      // maven article id
	Version     string `json:"version"`      // maven vension
	Snapshot    bool   `json:"snapshot"`     // 是否为快照
	Classifier  bool   `json:"classifier"`   // 源码， 值为 "sources"
	Extension   string `json:"extension"`    // 文件类型 jar pom war ,当 Source
	Signature   bool   `json:"signature"`    // 文件是否为 .sha1 signature
	Date        string `json:"date"`         // 快照的yyyMMdd
	Time        string `json:"time"`         // 时间 hhmmss
	BuildNumber string `json:"build_number"` // 快照的编译版本号
	Basic       bool   `json:"basic"`        // 基础版本， 非多build版本
	Deleted     bool   `json:"deleted"`      // 是否需要删除
}

/**
 * maven
 * + - group
 *   - artifact
 *   - version
 */
type Gav struct {
	Release  []*Info // 基础版本 包括 jar source pom
	Snapshot []*Info // 快照，里面还有build date time tag
}

/**
 * 格式化 matadata xml文件
 */
type Metadata struct {
	XMLName      xml.Name `xml:"metadata"`
	Text         string   `xml:",chardata"`
	ModelVersion string   `xml:"modelVersion,attr"`
	GroupId      struct {
		Text string `xml:",chardata"`
	} `xml:"groupId"`
	ArtifactId struct {
		Text string `xml:",chardata"`
	} `xml:"artifactId"`
	Version struct {
		Text string `xml:",chardata"`
	} `xml:"version"`
	Versioning struct {
		Text     string `xml:",chardata"`
		Snapshot struct {
			Text      string `xml:",chardata"`
			LocalCopy struct {
				Text string `xml:",chardata"`
			} `xml:"localCopy"`
		} `xml:"snapshot"`
		LastUpdated struct {
			Text string `xml:",chardata"`
		} `xml:"lastUpdated"`
		SnapshotVersions struct {
			Text            string `xml:",chardata"`
			SnapshotVersion []struct {
				Text       string `xml:",chardata"`
				Classifier struct {
					Text string `xml:",chardata"`
				} `xml:"classifier"`
				Extension struct {
					Text string `xml:",chardata"`
				} `xml:"extension"`
				Value struct {
					Text string `xml:",chardata"`
				} `xml:"value"`
				Updated struct {
					Text string `xml:",chardata"`
				} `xml:"updated"`
			} `xml:"snapshotVersion"`
		} `xml:"snapshotVersions"`
	} `xml:"versioning"`
}

/**
 *  GAV 唯一标识
 */
func (p *Info) Key() string {
	if p.GroupId == "" && p.ArticleId == "" && p.Version == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s", p.GroupId, p.ArticleId, p.Version)
}

/**
 * 根据名字解析数据
 */

func (p *Info) IsJarOrPomOrWar() bool {
	if strings.LastIndex(p.FullName, ".jar.sha1") != -1 {
		return false
	}
	if strings.LastIndex(p.FullName, ".jar") != -1 {
		p.Extension = "jar"
		return true
	}

	if strings.LastIndex(p.FullName, ".pom") != -1 {
		p.Extension = "pom"
		return true
	}
	if strings.LastIndex(p.FullName, ".war") != -1 {
		p.Extension = "war"
		return true
	}
	return false
}

/**
 * 忽略掉的文件
 * 原始文件
 * 1. ${ArticleId}-${Version}.jar
 * 2. ${ArticleId}-${Version}-sources.jar
 * 3. ${ArticleId}-${Version}.pom
 * 3. ${ArticleId}-${Version}.war
 * 4. *.sha1 签名文件
 */
func (gav *Info) Support() bool {
	simpleName := gav.SimpleName
	if simpleName == gav.ArticleId+"-"+gav.Version+".war" || simpleName == gav.ArticleId+"-"+gav.Version+".jar" || simpleName == gav.ArticleId+"-"+gav.Version+".pom" || simpleName == gav.ArticleId+"-"+gav.Version+"-sources.jar" || strings.LastIndex(simpleName, ".sha1") != -1 {
		return true
	}
	return false
}
