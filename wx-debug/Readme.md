# wx-server

用go语言编写的一个简单的微信公众号后台**本地测试脚本**

## 参考
1. [Go语言开发微信公众号](https://www.imooc.com/learn/783)


## 运行
```
// 在go-demos/wx-debug目录下，生成可执行文件(目录中本身也有)
go build


// 运行可执行文件
// Url：后台URL，如：localhost:8080
// Token：后台设置的TOKEN（同微信设置），如todorex123456
// FromUserName：开发者微信号，测试可为假，如fakeUserName
// ToUserName：发送方帐号（一个OpenID），测试可为假，如fakeToUserName
// Text：文本消息内容
./wx-dubug Url Token FromUserName ToUserName Text
```
关于微信公众号的配置请参考：
[微信公众平台技术文档](https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1445241432)