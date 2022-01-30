### 小工具

go get -u github.com/clibing/knife

or 

git clone github.com/clibing/knife

### 安装

下载二进制包，执行安装命令即可


* linux/macos 默认安装到 /usr/local/bin/
* win 默认安装到 C:/Windows/system32/

以macos为例
```shell
chmod a+x knife_darwin_amd64
./knife_darwin_amd64 install 
```


### 查看帮助

knife -h

主要是实现一些简单常用的方法

### 功能

* [x] 1.时间格式化工具，实现提取当前系统的戳毫秒(13位), 接收一个时间戳按照指定的格式进行格式化。
* [x] 2.URL编码解码
* [x] 3.加密计算， 默认接收字符串计算，支持计算指定的文件
  *  [x] 接收字符串
  *  [x] 加密文件
* [x] 4.xml,json,yml 互转与美化
  *  [x] 格式化json
  *  [x] json to xml
  *  [x] xml to json
  *  [x] yml to json
  *  [x] json to yml
* [x] 5.常用定时器cron表达式
* [ ] 6.图片生成 从Base64生成文件，根据文件生成Base64
* [x] 7.RSA公钥私钥生成器
* [x] 8.二维码生成器
* [x] 9.IP查询，支持本机ip、出口ip等
* [x] 10.markdown处理，支持从HTML转Markdown
* [x] 11.正则表达式
* [x] 12.CPU和内存使用率

