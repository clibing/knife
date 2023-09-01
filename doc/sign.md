# 参考

[参考](https://baijiahao.baidu.com/s?id=1607127562017641658&wfr=spider&for=pco)

怎么下载 https 网站的 公钥证书呢

* [Golang里的AES加密、解密](https://blog.51cto.com/u_13460811/5785309)
* [Golang之编码的使用](https://nljb.gitee.io/default/Golang%E4%B9%8B%E7%BC%96%E7%A0%81%E7%9A%84%E4%BD%BF%E7%94%A8/)
* [ECB CBC CFB](https://zhuanlan.zhihu.com/p/335091384)
* [对字符串进行AES加解密，CFB模式，128/192/256](https://www.golangcodes.com/forum.php?mod=viewthread&tid=66)

## AES

* aes是基于数据块的加密方式，也就是说，每次处理的数据时一块（16字节），当数据不是16字节的倍数时填充，这就是所谓的分组密码（区别于基于比特位的流密码），16字节是分组长度

分组加密的几种模式：

* ECB：是一种基础的加密方式，密文被分割成分组长度相等的块（不足补齐），然后单独一个个加密，一个个输出组成密文。
* CBC：是一种循环模式，前一个分组的密文和当前分组的明文异或或操作后再加密，这样做的目的是增强破解难度。
* CFB/OFB：实际上是一种反馈模式，目的也是增强破解的难度。
* FCB和CBC的加密结果是不一样的，两者的模式不同，而且CBC会在第一个密码块运算时加入一个初始化向量。

## AES CFB/OFB/ECB/CBC/CTR优缺点

AES常见加密模式有CFB/OFB/ECB/CBC/CTR，本文概述这些算法特点，让大家更快的了解AES，当然天缘也不是专业做算法的，工作中也只是使用到才会学习一点，如有错误，欢迎指出。

### 一、Cipher feedback（CFB）

CFB算法优点：

同明文不同密文，分组密钥转换为流密码。

CFB算法缺点：

串行运算不利并行，传输错误可能导致后续传输块错误。

### 二、Output feedback（OFB）

OFB算法优点：

同明文不同密文，分组密钥转换为流密码。

OFB算法缺点：

串行运算不利并行，传输错误可能导致后续传输块错误。

### 三、Electronic codebook(ECB)

ECB算法优点：

简单、孤立，每个块单独运算。适合并行运算。传输错误一般只影响当前块。

ECB算法缺点：

同明文输出同密文，可能导致明文攻击。我们平时用的AES加密很多都是ECB模式的，此模式加密不需要向量IV。

### 四、Cipher-block chaining（CBC）

CBC算法优点：

串行化运算，相同明文不同密文

CBC算法缺点：

需要初始向量，不过这其实不算缺点，下文的CTR也是需要随机数的。如果出现传输错误，那么后续结果解密后可能全部错误。

此外，还有Propagating cipher-block chaining（PCBC）加密模式，

### 五、Counter mode（CTR）

CTR算法优点：

无填充，同明文不同密文，每个块单独运算，适合并行运算。

CTR算法缺点：

可能导致明文攻击。

补充：

关于Padding补位问题，上文加密模式中，比如CBC等对输入块是有要求的，必须是块的整数倍，对不是整块的数据，要求进行填充，填充的方法有很多种，常见的有PKCS5和PKCS7、ISO10126等。

例如按照16字节分组的话：

对不足16字节部分（假设差n个满16字节），填充n个字节（n范围(1,15]），且每字节值均为n。
对正好16字节部分，则填充一个block，也就是补16个字节，每字节值为16
参考1：PKCS #7: Cryptographic Message Syntax

参考2：PKCS #5: Password-Based Cryptography Specification

所以上述算法中，默认：

需要Padding的有：CBC（，PCBC也需要，本文未涉及该加密模式）、ECB。

不需要Padding的有：CFB、OFB、CTR。

### 参考资料

* http://hi.baidu.com/doomsword/blog/item/ec73eb2b18f95435d52af120.html
* http://hi.baidu.com/tweetyf/item/36d2f94a8639320ae8350480
* http://zh.wikipedia.org/wiki/块密码的工作模式