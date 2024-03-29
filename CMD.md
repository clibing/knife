# 封装一些常用的小工具

1. 时间格式化工具，实现提取当前系统的戳毫秒(13位), 接收一个时间戳按照指定的格式进行格式化。
   knife time 输出13为时间戳
   knife time -h 帮助
   knife time -d 1h 以当前时间点向后增加1小时, 1h 可以替换为 -2h, 3m, 1d, -1d,
   knife time -i 1636610579860 格式当前时间戳

2. URL编码解码

   编码：knife url "https://github.com/clibing/knife"
   解码：knife url -e "http%3A%2F%2Fgithub.com%2Fclibing%2Fknife"
   编码并解码：knife url "https://github.com/clibing/knife" | knife url -e

3. 加密计算， 默认接收字符串计算，支持计算指定的文件
   knife sign -t md5 "clibing"
   knife sign -t md5 -s /tmp/data.txt 注意文件签名与指定字符串签名不一致， 因为文件最后含有一个\r\n 、\r之类的换行符是隐藏的
   echo "clibing" | knife sign -t md5
   其中md5可以替换 sha1, sha256, sha512, base64

4. xml,json,yml 互转与美化
   4.1 美化json
   knife json "{\"id\":1,\"name\":\"clibing\"}"
   输出
   {
           "id": 1,
           "name": "clibing"
   }

   4.2 xml 转换 json
   knife json -c 1 "<?xml version=\"1.0\" encoding=\"UTF-8\"?><name>clibing</name>"
   输出
   {"name": "clibing"}
   
   4.3 json 转换为 xml
   knife json -c 2 "{\"id\":1,\"name\":\"clibing\"}" 
   输出
   <?xml version="1.0" encoding="UTF-8"?> <object><number name="id">1</number><string name="name">clibing</string></object>

   4.4 json 转换为 yml
   knife json -c 3 "{\"id\":1,\"name\":\"clibing\"}" 
   输出
   id: 1
   name: clibing
   
   4.5 yml 转换 json
   knife json -c 4 "id: 1"
   输出
   {"id":1}

5. 定时器cron表达式
   knife cron 这是常用的cron表达式

6. 图片生成 从Base64生成文件，根据文件生成Base64
   无

7. 证书pem生成器
   knife rsa -b 1024

8. 二维码生成器
   8.1 当前目录快速生成二维码, 名字默认为 output.png
   knife qrcode "https://clibing.com"

   8.2 有边框，大小512，recovery level 2 输出到 /tmp/512.png 二维码的内容是 "https://clibing.com"
   knife qrcode -l 2 -s 512 -f /tmp/512.png "https://clibing.com"

   8.3 无边框，大小512，recovery level 2 输出到 /tmp/512.png 二维码的内容是 "https://clibing.com"
   knife qrcode -d -l 2 -s 512 -f /tmp/512.png "https://clibing.com"

9. IP查询，支持本机ip、出口ip等
   9.1 查看本机ip
   knife ip

   9.2 查看出口ip
   knife ip -e

10. markdown处理，支持从HTML转Markdown

   10.1. html -> markdown
   knife md -s /tmp/source.html -t /tmp/target.md

   10.2  markdown -> html
   knife md -d -s /tmp/target.md -t /tmp/source.html

11. 正则表达式

   11.1 根据正则执行 查找String模式
   knife reg -e "H(.*)d" "HelloWorld message "

   11.2 根据正则执行 匹配String模式
   knife reg -d -e "H(.*)d" "HelloWorld message "

12. 硬件检测

   12.1 检查当前cpu和内存使用率
   knife monitor

   12.1 检查当前cpu和内存使用率 并重复检查10
   knife monitor -t 10

13. 随机数

   13.1 默认生成uuid

   knife random

   13.2 生成包含数字、大小写字母、标点符号，长度为6, 共生成3个
   knife random -n -c -p -l 6 -t 3

更新帮助文档详见
knife -h
knife <command> -h