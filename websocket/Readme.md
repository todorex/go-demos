# WebSocket

用go语言实现一个WebSocket demo，实现一个线程安全的WebSocket连接接口

## 参考
1. [GO实现千万级WebSocket消息推送服务](https://www.imooc.com/learn/1025)

## 前置条件

1. [WebSocket工具包](https://github.com/gorilla/websocket)

## 运行
1. 启动服务端
    ```
    // 在go-demos/websocket/server目录下运行主server.go
    go run server.go
    ```
2. 打开测试页面
    在go-demos/websocket/web目录下的client.html
    
关于微信公众号的配置请参考：
[微信公众平台技术文档](https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1445241432)Å