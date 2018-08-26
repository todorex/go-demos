package main

import (
	"net/http"

	"go-demos/websocket/impl"

	"time"

	"log"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("localhost:8080", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 将加入 upgrade: websocket 请求头
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	log.Println("connection is open")
	conn, err := impl.InitConnection(wsConn)
	if err != nil {
		goto ERR
	}

	go func() {
		for {
			// 发送心跳包
			log.Println("heart beat")
			err := conn.WriteMessage([]byte("heart beat"))
			if err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		data, err := conn.ReadMessage()
		if err != nil {
			goto ERR
		}

		err = conn.WriteMessage(data)
		if err != nil {
			goto ERR
		}
	}

ERR:
	// 可重入的关闭
	conn.Close()
}
