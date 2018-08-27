# 日志监控系统

用go语言实现一个简单的日志监控系统

## 参考
1. [Go并发编程案例解析](https://www.imooc.com/learn/982)

## 流程

![log-monitor](https://github.com/todorex/go-demos/raw/master/log-monitor/image/log-monitor.png)

## 前置条件

1. [influxdb工具包](https://github.com/influxdata/influxdb)

2. docker安装influxdb
    ```
    #  docker run -idt --name influxdb -p 8086:8086 influxdb
    ```
    具体的数据库配置可以参考后面的[博客](https://studygolang.com/articles/10990?fr=sidebar)
3. docker安装grafana
    ```
    #  docker run \
         -d \
         -p 3000:3000 \
         --name=grafana \
         -e "GF_SERVER_ROOT_URL=http://grafana.server.name" \
         -e "GF_SECURITY_ADMIN_PASSWORD=secret" \
         grafana/grafana
    ```
    密码就是admin和你配置的secret，具体的时候可以参考后面提到的[官方文档](http://docs.grafana.org/)
    
## 运行
1. 启动日志生成程序
    ```
    # go run mock_data.go
    ```
2. 启动日志监控程序
   ```
   # go run log_monitor.go -path -influxDsn
   ```
   
   * path: 日志文件路径，如：/Users/rex/GoglandProjects/src/go-demos/log-monitor/access.log
   * influxDsn: influxdb数据库信息，如: http://127.0.0.1:8086@root@root@log_process@s
3. 查看监控信息
    * grafana的使用可以参考文档
    * 监控命令使用
        ```
        # curl 127.0.0.1:9193/monitor
        ```
        返回：
        ```
        {
                "handleLine": 183, // 处理日志行数
                "tps": 13.8, 
                "readChanLen": 0, // 读取channel的大小
                "writeChanLen": 0, // 写入channel的大小
                "runTime": "46.71648257s", // 程序运行时间
                "errNum": 0 // 错误处理行数
        }  
        ```
## 参考
1. [时序数据库InfluxDB使用详解](https://studygolang.com/articles/10990?fr=sidebar)
2. [grafana文档](http://docs.grafana.org/)