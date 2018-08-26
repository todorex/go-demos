package impl

import (
	"sync"

	"github.com/go-errors/errors"
	"github.com/gorilla/websocket"
)

// 封装websocket连接
type Connection struct {
	// 连接
	wsConn *websocket.Conn
	// 读取的消息channel
	inChan chan []byte
	// 写入的消息channel
	outChan chan []byte
	// 协助关闭的channel
	closeChan chan byte
	// 锁
	mutex sync.Mutex
	// 关闭标志位
	isClosed bool
}

// 初始化（包装）连接
func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}
	// 启动读协程
	go conn.readLoop()

	// 启动写协程
	go conn.writeLoop()

	return
}

// 读取消息
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}

	return
}

// 写入消息
func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")

	}

	return
}

// 关闭连接
func (conn *Connection) Close() {
	// 线程安全，可重入
	conn.wsConn.Close()

	// 如果有一个操作触发了关闭连接操作，则也将另一个操作关闭，但是该操作不可重入，需保证只执行一次
	// 在并发情况下需要加锁
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

// 准备读协程
func (conn *Connection) readLoop() {

	for {
		_, data, err := conn.wsConn.ReadMessage()
		if err != nil {
			goto ERR
		}

		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			goto ERR
		}

	}

ERR:
	conn.Close()
}

// 准备写协程
func (conn *Connection) writeLoop() {
	var data []byte
	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR

		}

		err := conn.wsConn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}
