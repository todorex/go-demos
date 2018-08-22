package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func writeLog(r *http.Request, t time.Time, match string, pattern string) {

	// 不是产品级别
	if logLevel != "prod" {
		// 拿到时间间隔
		d := time.Now().Sub(t)

		l := fmt.Sprintf("[ACCESS] | % -10s | % -40s | % -16s | % -10s | % -40s |", r.Method, r.URL.Path, d.String(), match, pattern)

		log.Println(l)
	}
}
