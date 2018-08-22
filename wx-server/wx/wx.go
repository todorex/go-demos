package wx

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/mxj"
)

// 微信请求字段
type weixinQuery struct {
	// `` 表示tag，类似java的注解，在json转化的时候会用到
	Signature    string `json:"signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	EncryptType  string `json:"encrypt_type"`
	MsgSignature string `json:"msg_signature"`
	Echostr      string `json:"echostr"`
}

type WeixinClient struct {
	Token          string
	Query          weixinQuery
	Message        map[string]interface{}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Methods        map[string]func() bool
}

// 构造客户端
func NewClient(r *http.Request, w http.ResponseWriter, token string) (*WeixinClient, error) {

	weixinClient := new(WeixinClient)

	weixinClient.Token = token
	weixinClient.Request = r
	weixinClient.ResponseWriter = w

	weixinClient.initWeixinQuery()

	// 校验token是否正确
	if weixinClient.Query.Signature != weixinClient.signature() {
		return nil, errors.New("Invalid Signature.")
	}

	return weixinClient, nil
}

func (this *WeixinClient) initWeixinQuery() {

	var q weixinQuery

	q.Nonce = this.Request.URL.Query().Get("nonce")
	q.Echostr = this.Request.URL.Query().Get("echostr")
	q.Signature = this.Request.URL.Query().Get("signature")
	q.Timestamp = this.Request.URL.Query().Get("timestamp")
	q.EncryptType = this.Request.URL.Query().Get("encrypt_type")
	q.MsgSignature = this.Request.URL.Query().Get("msg_signature")

	this.Query = q
}

// 计算签名
func (this *WeixinClient) signature() string {

	// 排序字符串数组并拼接为一个字符串
	strs := sort.StringSlice{this.Token, this.Query.Timestamp, this.Query.Nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	// 利用sha1生成签名
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 得到消息内容
func (this *WeixinClient) initMessage() error {

	body, err := ioutil.ReadAll(this.Request.Body)

	if err != nil {
		return err
	}

	// 将xml解析成map
	m, err := mxj.NewMapXml(body)

	if err != nil {
		return err
	}

	// 断言xml这个key是否存在
	if _, ok := m["xml"]; !ok {
		return errors.New("Invalid Message.")
	}
	// 断言xml这个key对应的值是否为一个map
	message, ok := m["xml"].(map[string]interface{})

	if !ok {
		return errors.New("Invalid Field `xml` Type.")
	}
	// 拿到消息这个map
	this.Message = message

	log.Println(this.Message)

	return nil
}

// 构造要回复的内容XML
func (this *WeixinClient) text() {
	// 拿到context标签里的值
	inMsg, ok := this.Message["Content"].(string)

	if !ok {
		return
	}
	// 要回复的内容XML
	var reply TextMessage

	reply.InitBaseData(this, "text")
	// 转换content内容
	//if inMsg == "孙佳琦" {
	//	inMsg += "：超级无敌优秀"
	//} else {
	//	inMsg += "：垃圾"
	//}

	reply.Content = value2CDATA(fmt.Sprintf("我收到是：%s", inMsg))

	// 转化为XML
	replyXml, err := xml.Marshal(reply)

	if err != nil {
		log.Println(err)
		// 返回403
		this.ResponseWriter.WriteHeader(http.StatusForbidden)
		return
	}
	// 响应头
	this.ResponseWriter.Header().Set("Content-Type", "text/xml")
	this.ResponseWriter.Write(replyXml)
}

// 启动客户端，执行处理逻辑
func (this *WeixinClient) Run() {

	err := this.initMessage()

	if err != nil {

		log.Println(err)
		this.ResponseWriter.WriteHeader(http.StatusForbidden)
		return
	}
	// 拿到消息种类，这里处理只处理文本消息
	MsgType, ok := this.Message["MsgType"].(string)

	if !ok {
		this.ResponseWriter.WriteHeader(http.StatusForbidden)
		return
	}

	// 处理不同的消息种类
	switch MsgType {
	case "text":
		this.text()
	default:
	}

	return
}
