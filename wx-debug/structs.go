package main

import "encoding/xml"

type CDATAText struct {
	Text string `xml:",innerxml"`
}

type Base struct {
	FromUserName CDATAText
	ToUserName   CDATAText
	MsgType      CDATAText
	CreateTime   int64
	MsgId        int64 `xml:",omitempty"`
}

// 文本消息
type Text struct {
	XMLName xml.Name `xml:"xml"`
	Base
	Content CDATAText
}

// 图片消息
type Image struct {
	XMLName xml.Name `xml:"xml"`
	Base
	PicUrl  CDATAText
	MediaId CDATAText
}

// 语音消息
type Voice struct {
	XMLName xml.Name `xml:"xml"`
	Base
	MediaId     CDATAText
	Format      CDATAText
	Recognition CDATAText `xml:",omitempty"`
}

// 视频消息
type Video struct {
	XMLName xml.Name `xml:"xml"`
	Base
	MediaId      CDATAText
	ThumbMediaId CDATAText
}

// 小视频消息
type ShortVideo struct {
	XMLName xml.Name `xml:"xml"`
	Base
	MediaId      CDATAText
	ThumbMediaId CDATAText
}

// 地理位置消息
type Location struct {
	XMLName xml.Name `xml:"xml"`
	Base
	Location_X float64
	Location_Y float64
	Scale      int
	Label      CDATAText
}
