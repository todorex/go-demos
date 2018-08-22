package wx

import (
	"encoding/xml"
	"strconv"
	"time"
)

type Base struct {
	FromUserName CDATAText
	ToUserName   CDATAText
	MsgType      CDATAText
	CreateTime   CDATAText
}

func (b *Base) InitBaseData(w *WeixinClient, msgtype string) {

	b.FromUserName = value2CDATA(w.Message["ToUserName"].(string))
	b.ToUserName = value2CDATA(w.Message["FromUserName"].(string))
	b.CreateTime = value2CDATA(strconv.FormatInt(time.Now().Unix(), 10))
	b.MsgType = value2CDATA(msgtype)
}

// CDATA会被浏览器解析器忽略
// CDATA里的内容
type CDATAText struct {
	Text string `xml:",innerxml"`
}

// 文本内容
// 微信文本请求的内容为xml
// <xml>
// 	<ToUserName>< ![CDATA[toUser] ]></ToUserName>
// 	<FromUserName>< ![CDATA[fromUser] ]></FromUserName>
// 	<CreateTime>1348831860</CreateTime>
//  <MsgType>< ![CDATA[text] ]></MsgType>
//  <Content>< ![CDATA[this is a test] ]></Content>
//  <MsgId>1234567890123456</MsgId>
// </xml>

type TextMessage struct {
	XMLName xml.Name `xml:"xml"`
	Base
	Content CDATAText
}
