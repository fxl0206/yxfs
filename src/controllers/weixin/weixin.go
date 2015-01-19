package weixin

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

//weixin Controller 请求路径
//var Prefix string = "/sign"
var SignPrefix string = "/sign"
var ActionPrefix string = "/action"

//微信请求消息结构
type Request struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}

//微信响应头结构
type Response struct {
	ToUserName   string `xml:"xml>ToUserName"`
	FromUserName string `xml:"xml>FromUserName"`
	CreateTime   string `xml:"xml>CreateTime"`
	MsgType      string `xml:"xml>MsgType"`
	Content      string `xml:"xml>Content"`
	MsgId        int    `xml:"xml>MsgId"`
}

//服务器Token认证
func DoSign(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(SignPrefix):]
	log.Println(path)
	r.ParseForm()
	token := "my tocken"
	var signature string = strings.Join(r.Form["signature"], "")
	var timestamp string = strings.Join(r.Form["timestamp"], "")
	var nonce string = strings.Join(r.Form["nonce"], "")
	var echostr string = strings.Join(r.Form["echostr"], "")
	tmps := []string{token, timestamp, nonce}
	sort.Strings(tmps)
	tmpStr := tmps[0] + tmps[1] + tmps[2]
	tmp := str2sha1(tmpStr)
	if tmp == signature {
		log.Println("signature Success")
		fmt.Fprintf(w, echostr)
	} else {
		log.Println("signature Failed!")
	}
}

//消息处理
func DoAction(w http.ResponseWriter, r *http.Request) {
	log.Println("here")

	postedMsg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()
	v := Request{}
	xml.Unmarshal(postedMsg, &v)
	if v.MsgType == "text" {
		v := Request{v.ToUserName, v.FromUserName, v.CreateTime, v.MsgType, v.Content, v.MsgId}
		output, err := xml.MarshalIndent(v, " ", " ")
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
		fmt.Fprintf(w, string(output))
	} else if v.MsgType == "event" {
		Content := `"欢迎关注  我的微信"`
		v := Request{v.ToUserName, v.FromUserName, v.CreateTime, v.MsgType, Content, v.MsgId}
		output, err := xml.MarshalIndent(v, " ", " ")
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
		fmt.Fprintf(w, string(output))
	} else {
		log.Println("Something is Wrong!")
	}
}
func str2sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
