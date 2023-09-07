
<a name="v0.0.9"></a>
## [v0.0.9](https://github.com/clibing/knife/compare/v0.0.7...v0.0.9) (2023-09-07)

### Fix

* **aes:** add aes cbc cfb ecb
* **aes:** support set iv, default cbc use key
* **base64:** 重构base64
* **base64:** support write file
* **carba:** many values
* **client:** http client format json
* **client:** add http client
* **cmd:** rename
* **code:** remove not used file
* **code:** go mod tidy error
* **code:** delete useless code
* **command:** add client command
* **command:** version 注释 增加换行
* **common:** 抽取common分类
* **debug:** remove short cmd
* **debug:** fix debug
* **debug:** add maven debug info
* **fmt:** 修复参数传入...
* **format:** [] 存在 fmt.println [] how to fix me
* **get:** return single value when sentinel is true
* **go:** go build -v $GOPATH
* **go:** update dep
* **help:** add help
* **image:** image base64互转
* **ip:** 调整url,允许指定
* **key:** key is empty
* **log:** add log debug info
* **markdown:** move markdown to convert
* **maven:** Clean Maven snapshot dependencies
* **md:** add changelog md
* **md5:** 重构md5
* **move:** move
* **note:** add note
* **note:** add redis help document
* **redis:** add set result info
* **redis:** add redis get
* **redis:** add get set del command
* **redis:** add redis -k -v -e
* **redis:** add redis role
* **root:** 修改说明
* **rsa:** size to 2048
* **sha1:** 修复fmt输出 去除换行
* **sign:** md5 sha1
* **sign:** add sign method
* **sync:** sync
* **tag:** add 0.0.8
* **tag:** tag
* **url:** move url to convert/
* **version:** update to v0.0.9

### Refactor

* **client:** move http redis websocket to client/
* **code:** 调整代码结构
* **code:** 调整代码结构
* **code:** 调整代码结构
* **code:** 重构code

### Style

* **format:** format


<a name="v0.0.7"></a>
## [v0.0.7](https://github.com/clibing/knife/compare/0.0.9...v0.0.7) (2023-06-09)

### Fix

* **code:** 优化代码
* **code:** 优化代码
* **sub:** sub command


<a name="0.0.9"></a>
## [0.0.9](https://github.com/clibing/knife/compare/v0.0.8...0.0.9) (2023-06-09)

### Feat

* **static:** add static web, support many computer download file

### Fix

* **build:** add build shell on linux amd64 docker
* **code:** code
* **download:** not test
* **fmt:** remove fmt println
* **mod:** update go.mod go.sum
* **name:** update name upper or lower
* **note:** note download command
* **revert:** revert code
* **tag:** update to 0.0.9


<a name="v0.0.8"></a>
## [v0.0.8](https://github.com/clibing/knife/compare/v0.0.6...v0.0.8) (2023-05-23)

### Feat

* **arch:** generator download url with os arch version
* **download:** add download with process

### Fix

* **cmd:** update version format
* **download:** check err
* **download:** debug
* **ip:** get remote ip use myself service
* **mod:** update mod
* **randome:** create times uuid
* **rsa:** add pkcs8 file gen
* **upgrade:** go upgrade to 1.20.x
* **version:** upgrade


<a name="v0.0.6"></a>
## [v0.0.6](https://github.com/clibing/knife/compare/0.0.5...v0.0.6) (2022-04-07)

### Feat

* **upgrade:** add bin auto upgrade

### Fix

* **wol:** add wol function


<a name="0.0.5"></a>
## [0.0.5](https://github.com/clibing/knife/compare/0.0.4...0.0.5) (2022-03-30)

### Feat

* **random:** 生成随机数
* **version:** update to 0.0.5

### Fix

* **note:** note
* **note:** readme remove statistics
* **readme.md:** 增加统计


<a name="0.0.4"></a>
## [0.0.4](https://github.com/clibing/knife/compare/0.0.3...0.0.4) (2022-03-15)

### Fix

* **time:** 重构time函数
* **time:** add string time to long mills


<a name="0.0.3"></a>
## [0.0.3](https://github.com/clibing/knife/compare/0.0.2...0.0.3) (2022-01-30)

### Docs

* **chglog:** add git chglog
* **note:** note

### Feat

* **monitor:** add monitor, cpu memory

### Fix

* **cpu:** add run times
* **cpu:** 增加获取温度
* **cpu:** add cpu temperature
* **readme:** add readme
* **unit:** fix unit


<a name="0.0.2"></a>
## [0.0.2](https://github.com/clibing/knife/compare/0.0.1...0.0.2) (2021-11-12)

### Fix

* **0.0.2:** 完善帮助文档
* **code:** 优化代码，增加测试
* **qrcode:** 修复二维码识别不准确
* **version:** add version


<a name="0.0.1"></a>
## 0.0.1 (2021-11-10)

### Docs

* **readme:** remove exist description
* **readme:** 更新功能的支持
* **readme:** update fromat

### Feat

* **json:** add yml json xml
* **markdown:** 支持markdown html互转
* **qrcode:** add qrcode tool
* **reg:** reg
* **release:** 发布第一个小版本
* **time:** add time format

### Fix

* **crypt:** support string， not supoort file
* **readme:** remove space
* **revert:** instal和uninstall代码混乱了

