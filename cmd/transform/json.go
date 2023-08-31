package transform

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
	"vimagination.zapto.org/json2xml"
)

var (
	prefix, indent string
	convert        int
)

const (
	formatPrint = iota
	xmlToJson
	jsonToXml
	jsonToYml
	ymlToJson
)

// convertCmd represents the json command
var convertCmd = &cobra.Command{
	Use:   "json",
	Short: "json xml yml转换器",
	Long: `美化json，json与xml转换, json与yml转换:

0. 美化json
   knife json "{\"id\":1,\"name\":\"clibing\"}" 
   输出
   {
           "id": 1,
           "name": "clibing"
   }
1. xml 转换 json
   knife json -c 1 "<?xml version=\"1.0\" encoding=\"UTF-8\"?><name>clibing</name>"
   输出
   {"name": "clibing"}
   
2. json 转换为 xml
   knife json -c 2 "{\"id\":1,\"name\":\"clibing\"}" 
   输出
   <?xml version="1.0" encoding="UTF-8"?> <object><number name="id">1</number><string name="name">clibing</string></object>

3. json 转换为 yml
   knife json -c 3 "{\"id\":1,\"name\":\"clibing\"}" 
   输出
   id: 1
   name: clibing
   
4. yml 转换 json
   knife json -c 4 "id: 1" 
   输出
   {"id":1}

-----------------------------------
话外篇 可以使用 jq c语言实现的格式化工具

如果此格式化不能满足 可以安装jq进行格式化
jq: https://github.com/stedolan/jq.git
使用:
echo '{"id":1,"name":"clibing"}' | jq.`,
	Run: func(_ *cobra.Command, args []string) {
		// 使用的参数传递直接返回
		if len(args) > 0 {
			doJsonFunc(args)
			return
		}

		fileInfo, _ := os.Stdin.Stat()
		if (fileInfo.Mode() & os.ModeNamedPipe) != os.ModeNamedPipe {
			return
		}
		var buf strings.Builder
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			buf.WriteString(s.Text())
		}
		pip := []string{buf.String()}
		doJsonFunc(pip)
	},
}

func doJsonFunc(content []string) {
	for _, value := range content {
		switch convert {
		case xmlToJson:
			xmlToJsonFunc(value)
			return
		case jsonToXml:
			jsonToXmlFunc(value)
			return
		case jsonToYml:
			jsonToYmlFunc(value)
			return
		case ymlToJson:
			ymlToJsonFunc(value)
			return
		default:
			befaultJson(value)
			return
		}
	}
}

func befaultJson(content string) (err error) {
	// 接收管道
	option := &pretty.Options{Width: 80, Prefix: prefix, Indent: indent, SortKeys: false}
	fmt.Printf("%s\n", pretty.PrettyOptions([]byte(content), option))
	return
}

func xmlToJsonFunc(content string) (err error) {
	json, err := xj.Convert(strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("xml to json error, %s", err)
	}
	fmt.Println(json.String())
	return
}

func jsonToXmlFunc(content string) (err error) {
	var buf strings.Builder
	x := xml.NewEncoder(&buf)
	if err := json2xml.Convert(json.NewDecoder(strings.NewReader(content)), x); err != nil {
		return fmt.Errorf("json to xml error, %s", err)
	}
	x.Flush()
	output := buf.String()
	fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`, output)
	return
}

func jsonToYmlFunc(content string) (err error) {
	yml, err := yaml.JSONToYAML([]byte(content))
	if err != nil {
		return fmt.Errorf("json to yml err: %s", err)
	}
	fmt.Println(string(yml))
	return
}

func ymlToJsonFunc(content string) (err error) {
	json, err := yaml.YAMLToJSON([]byte(content))
	if err != nil {
		return fmt.Errorf("yml to json err: %s", err)
	}
	fmt.Println(string(json))
	return
}

func init() {
	// Here you will define your flags and configuration settings.
	convertCmd.Flags().StringVarP(&indent, "indent", "i", "\t", "json格式化缩进标记，默认制表符")
	convertCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "json格式化前缀")
	convertCmd.Flags().IntVarP(&convert, "convert", "c", 0, "转换模式\n0: json格式化\n1: xml to json\n2: json to xml(建议使用struct Tag)\n3: json to yml\n4: yml to json, 默认为0美化")
}

func NewJsonConvert() *cobra.Command {
	return convertCmd
}
