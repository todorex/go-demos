package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"flag"

	"encoding/json"
	"net/http"

	"github.com/influxdata/influxdb/client/v2"
)

// 抽象一下
type Reader interface {
	// 从channel读
	Read(rc chan []byte)
}

type Writer interface {
	// 写入channel
	Write(wc chan *Message)
}

type ReadFromFile struct {
	// 文件路径
	path string
}

// 读取模块
func (r *ReadFromFile) Read(rc chan []byte) {
	// 打开文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error: %s", err.Error()))
	}

	// 从文件末尾开始逐行读取文件内容
	f.Seek(0, 2)
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes error: %s", err.Error()))
		}
		TypeMonitorChan <- TypeHandleLine
		rc <- line[:len(line)-1]
	}

}

type WriteToInfluxDB struct {
	// influxDB信息
	influxDBDsn string
}

// 写入模块
func (w *WriteToInfluxDB) Write(wc chan *Message) {
	// 解析influxDBDsn
	infSli := strings.Split(w.influxDBDsn, "@")

	// github 示例
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     infSli[0],
		Username: infSli[1],
		Password: infSli[2],
	})
	if err != nil {
		log.Fatal("client error:", err)

	}
	defer c.Close()

	for v := range wc {
		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  infSli[3],
			Precision: infSli[4],
		})
		if err != nil {
			log.Fatal(err)
		}

		// Create a point and add to batch
		// Tags(索引): Path, Method, Scheme, Status
		// fields: UpstreamTime, RequestTime, BytesSent
		tags := map[string]string{"Path": v.Path, "Method": v.Method, "Scheme": v.Scheme, "Status": v.Status}
		fields := map[string]interface{}{
			"UpstreamTime": v.UpstreamTime,
			"RequestTime":  v.RequestTime,
			"BytesSent":    v.BytesSent,
		}

		pt, err := client.NewPoint("nginx_log", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		// Write the batch
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}

		// Close client resources
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}

		log.Println("write success!!")
	}

}

// 日志处理模块
type LogProcess struct {
	// 读取模块channel
	rc chan []byte
	// 写入模块channel
	wc     chan *Message
	reader Reader
	writer Writer
}

// 根据分析需求确定，本次将分析某个协议下的某个请求在某个请求方法的QPS、响应时间、流量
type Message struct {
	TimeLocal time.Time
	// 流量
	BytesSent int
	Path      string
	Method    string
	// 协议
	Scheme       string
	Status       string
	UpstreamTime float64
	RequestTime  float64
}

// 解析模块
func (lp *LogProcess) Process() {

	// 解析日志数据
	reg := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Panicf("LoadLocation failure: %s", err.Error())
	}
	for v := range lp.rc {
		result := reg.FindStringSubmatch(string(v))
		if len(result) != 14 {
			TypeMonitorChan <- TypeErrNum
			log.Printf("FindStringSubmatch fail: %s \n", string(v))
			continue
		}

		message := &Message{}
		// 解析TimeLocal
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", result[4], location)
		if err != nil {
			TypeMonitorChan <- TypeErrNum
			log.Println("ParseInLocation failure:", err.Error(), result[4])
			continue
		}
		message.TimeLocal = t

		// 解析流量
		bytesSent, _ := strconv.Atoi(result[8])
		message.BytesSent = bytesSent

		// 解析方法、路径
		reqSli := strings.Split(result[6], " ")
		if len(reqSli) != 3 {
			TypeMonitorChan <- TypeErrNum
			log.Println("strings.Split failure:", result[6])
			continue
		}
		message.Method = reqSli[0]

		u, err := url.Parse(reqSli[1])
		if err != nil {
			TypeMonitorChan <- TypeErrNum
			log.Println("url.Parse fail:", err.Error())
			continue
		}
		message.Path = u.Path

		message.Scheme = result[5]
		message.Status = result[7]

		upstreamTime, _ := strconv.ParseFloat(result[12], 64)
		requestTime, _ := strconv.ParseFloat(result[13], 64)
		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime

		lp.wc <- message
	}
}

// 系统状态监控（这些监控数据可以通过HTTP接口暴露出去）
type SystemInfo struct {
	// 总处理日志行数
	HandleLine int `json:"handleLine"`
	// 系统吞出量
	Tps float64 `json:"tps"`
	// read channel 长度
	ReadChanLen int `json:"readChanLen"`
	// write channel 长度
	WriteChanLen int `json:"writeChanLen"`
	// 运行总时间
	RunTime string `json:"runTime"`
	// 错误数
	ErrNum int `json:"errNum"`
}

// 枚举
const (
	// 总处理日志行数
	TypeHandleLine = 0
	// 错误数
	TypeErrNum = 1
)

var TypeMonitorChan = make(chan int, 200)

type Monitor struct {
	// 系统运行时间
	startTime time.Time
	// 系统信息
	data SystemInfo
	// 5s的处理函数，模拟tps
	tpsSli []int
}

// 开启监控
func (m *Monitor) start(lp *LogProcess) {
	// 消费TypeMonitorChan数据
	go func() {
		for n := range TypeMonitorChan {
			switch n {
			case TypeHandleLine:
				m.data.HandleLine += 1
			case TypeErrNum:
				m.data.ErrNum += 1
			}
		}
	}()
	// 每5秒产生一次tps数据
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-ticker.C
			m.tpsSli = append(m.tpsSli, m.data.HandleLine)
			// 维护最新的tps（每次只保存2个数据）
			if len(m.tpsSli) > 2 {
				m.tpsSli = m.tpsSli[1:]
			}
		}
	}()

	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		m.data.RunTime = time.Now().Sub(m.startTime).String()
		m.data.ReadChanLen = len(lp.rc)
		m.data.WriteChanLen = len(lp.wc)

		if len(m.tpsSli) >= 2 {
			m.data.Tps = float64(m.tpsSli[1]-m.tpsSli[0]) / 5
		}

		// 转成json
		result, _ := json.MarshalIndent(m.data, "", "\t")
		// 返回
		io.WriteString(w, string(result))
	})

	http.ListenAndServe(":9193", nil)
}

func main() {

	var path, influxDsn string

	flag.StringVar(&path, "path", "/Users/rex/GoglandProjects/src/go-demos/log-monitor/access.log", "log file path")

	flag.StringVar(&influxDsn, "influxDsn", "http://127.0.0.1:8086@root@root@log_process@s", "influx data source")
	reader := &ReadFromFile{
		path: path,
	}

	flag.Parse()

	writer := &WriteToInfluxDB{
		// @后面分别代表了用户名，密码，数据库，精度
		influxDBDsn: influxDsn,
	}

	lp := &LogProcess{
		rc:     make(chan []byte, 200),
		wc:     make(chan *Message, 200),
		reader: reader,
		writer: writer,
	}
	go lp.reader.Read(lp.rc)
	// 由于处理模块和写模块比较慢，所以可以多开几个goroutine加快处理
	for i := 0; i < 2; i++ {
		go lp.Process()
	}

	for i := 0; i < 4; i++ {
		go lp.writer.Write(lp.wc)
	}

	monitor := &Monitor{
		startTime: time.Now(),
		data:      SystemInfo{},
	}
	monitor.start(lp)
}
