# wx-server

用go语言编写的一个简单的微信公众号后台

## 参考
1. [Go语言开发微信公众号](https://www.imooc.com/learn/783)

## 开发总览
![微信公众号请求流程]()

## 前置条件

1. [解析XML为Map的包](https://github.com/clbanning/mxj)

## 运行
```
// 在go-demos/wx-serve目录下，生成可执行文件(目录中本身也有)
go build
// 运行可执行文件
./wx-server
```
关于微信公众号的配置请参考：
[微信公众平台技术文档](https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1445241432)