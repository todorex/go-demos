package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
)

// 随机字符串
func randStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 格式化返回的XML
func formatXML(data []byte) ([]byte, error) {
	b := &bytes.Buffer{}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	encoder := xml.NewEncoder(b)
	encoder.Indent("", "	")
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			encoder.Flush()
			return b.Bytes(), nil
		}
		if err != nil {
			return nil, err
		}
		err = encoder.EncodeToken(token)
		if err != nil {
			return nil, err
		}
	}
}

// 字符串转成<![CDATA[" + str + "]]>"
func strToCDATA(str string) CDATAText {
	return CDATAText{"<![CDATA[" + str + "]]>"}
}

// 模拟签名
func signature(timestamp string, nonce string, token string) string {
	strs := sort.StringSlice{token, timestamp, nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 发送消息
func send(url string, message string) ([]byte, error) {
	// 客户端
	client := &http.Client{}
	// 发送POST请求
	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		return nil, err
	}
	// 请求头
	req.Header.Set("Content-Type", "text/xml")
	// 请求
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	// 注意关闭
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
