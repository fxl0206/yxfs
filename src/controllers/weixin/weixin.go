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

type Request struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}

type Response struct {
	ToUserName   string `xml:"xml>ToUserName"`
	FromUserName string `xml:"xml>FromUserName"`
	CreateTime   string `xml:"xml>CreateTime"`
	MsgType      string `xml:"xml>MsgType"`
	Content      string `xml:"xml>Content"`
	MsgId        int    `xml:"xml>MsgId"`
}

var Prefix string = "/sign"

func DoSign(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(Prefix):]
	fmt.Println("你的token")
	log.Println(path)
	var token string = "你的token"
	log.Println(token)
	r.ParseForm()
	var signature string = strings.Join(r.Form["signature"], "")
	var timestamp string = strings.Join(r.Form["timestamp"], "")
	var nonce string = strings.Join(r.Form["nonce"], "")
	var echostr string = strings.Join(r.Form["echostr"], "")
	tmps := []string{token, timestamp, nonce}
	sort.Strings(tmps)
	tmpStr := tmps[0] + tmps[1] + tmps[2]
	log.Println(tmpStr)
	tmp := str2sha1(tmpStr)

	if tmp == signature {
		fmt.Fprintf(w, echostr)
	}
}
func action(w http.ResponseWriter, r *http.Request) {
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
		Content := `"欢迎关注  
                                我的微信"`
		v := Request{v.ToUserName, v.FromUserName, v.CreateTime, v.MsgType, Content, v.MsgId}
		output, err := xml.MarshalIndent(v, " ", " ")
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
		fmt.Fprintf(w, string(output))
	}
}
func str2sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
