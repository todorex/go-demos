package main

import (
	"go-demos/wx-server/wx"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

// 处理器
type httpHandler struct{}

// controller处理器
type WebController struct {
	// 处理方法名
	Function func(http.ResponseWriter, *http.Request)
	// Get/Post...
	Method string
	// url模式
	Pattern string
}

// controller处理器数组
var controllers []WebController

// 初始化controller处理器数组
func init() {
	controllers = append(controllers, WebController{post, "POST", "^/"})
	controllers = append(controllers, WebController{get, "GET", "^/"})
}

// 处理get请求（用于微信公众平台验证）
func get(w http.ResponseWriter, r *http.Request) {

	client, err := wx.NewClient(r, w, token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	if len(client.Query.Echostr) > 0 {
		// 直接返回原字符串就可以
		w.Write([]byte(client.Query.Echostr))
		return
	}

	w.WriteHeader(http.StatusForbidden)
	return
}

// 处理post请求
func post(w http.ResponseWriter, r *http.Request) {

	client, err := wx.NewClient(r, w, token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}
	// 执行我们的处理逻辑
	client.Run()
	return
}

// 实现Handler接口
func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 用于日志记录
	t := time.Now()
	// 循环controller处理器处理对应的URL
	for _, webController := range controllers {
		// 匹配路径
		if m, _ := regexp.MatchString(webController.Pattern, r.URL.Path); m {
			// 匹配方法GET/POST
			if r.Method == webController.Method {
				// 处理请求
				webController.Function(w, r)
				// 写日志出去了
				go writeLog(r, t, "match", webController.Pattern)

				return
			}
		}
	}

	go writeLog(r, t, "unmatch", "")

	io.WriteString(w, "")
	return
}
