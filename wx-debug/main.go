package main

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	// 接受命令行参数（5个）
	arg_num := len(os.Args)
	if arg_num < 6 {
		fmt.Println("No Enough Args")
		return
	}
	// 拿到5个参数
	// 发送URL
	url := os.Args[1]
	// TOKEN
	token := os.Args[2]
	// FROM
	from := os.Args[3]
	// TO
	to := os.Args[4]
	// 内容
	text := os.Args[5]

	// 模拟微信请求参数
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	nonce := randStringRunes(8)
	sign := signature(timestampStr, nonce, token)
	url = fmt.Sprintf("%s?signature=%s&timestamp=%s&nonce=%s", url, sign, timestampStr, nonce)

	// 设置消息内容
	var message Text
	message.FromUserName = strToCDATA(from)
	message.ToUserName = strToCDATA(to)
	message.MsgType = strToCDATA("text")
	message.MsgId = rand.Int63()
	message.Content = strToCDATA(text)
	message.CreateTime = timestamp
	// 转成XML
	xml, err := xml.Marshal(message)

	if err != nil {
		fmt.Println(err)
		return
	}
	// 发送
	resp, err := send(url, string(xml))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("URL:", url)
	fmt.Println("--------------------")
	fmt.Println("Send Message:")

	x, err := formatXML(xml)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(x))
	}

	fmt.Println("--------------------")
	fmt.Println("Response:")

	if resp == nil {
		fmt.Println("--------------------")
		return
	}
	// 打印返回
	x, err = formatXML(resp)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(x))
	}

	fmt.Println("--------------------")
}
