package yxserver

import (
	"code.google.com/p/go.net/websocket"
	"controllers/json"
	"controllers/weixin"
	"controllers/zkvalue"
	"fmt"
	"github.com/widuu"
	"gmzoo.com/wsconn"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

//全局配置文路径
var ConfPath string

//websocket请求路径
var socketPath string

//系统配置初始化
func NewInstance(confPath string, serverName string) WbServer {
	conf := goini.SetConfig(confPath)
	port := conf.GetValue(serverName, "port")
	viewPath := conf.GetValue(serverName, "viewPath")
	ConfPath = confPath
	wsPath := conf.GetValue("websocket", "wsPath")
	wsPort := conf.GetValue("websocket", "port")
	wsIp := conf.GetValue("websocket", "wsIp")
	socketPath = "ws://" + wsIp + ":" + wsPort + "/" + wsPath
	return WbServer{port, viewPath, serverName}
}

type WbServer struct {
	port       string
	viewPath   string
	ServerName string
}

//启动服务器，绑定请求控制器
func (this *WbServer) Start() {
	//处理Websocket请求
	http.Handle("/ws", websocket.Handler(DealSocketReq))
	//请求静态资源文件
	http.HandleFunc("/"+this.viewPath+"/", StaticServer)
	//请求zookeeper菜单json数据
	http.HandleFunc("/getMenu", json.EchoMenu)
	//获取zookeeper节点数据
	http.HandleFunc(zkvalue.Prefix+"/", zkvalue.EchoValue)
	//微信token验证
	http.HandleFunc(weixin.SignPrefix+"/", weixin.DoSign)
	//微信消息处理
	http.HandleFunc(weixin.ActionPrefix+"/", weixin.DoAction)

	//监听端口
	err := http.ListenAndServe(":"+this.port, nil)
	//服务器启动错误处理
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

//处理静态文件请求
func StaticServer(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/"):]
	log.Println(r.URL.Path)
	if m, _ := regexp.MatchString("/json", r.URL.Path); m {
		w.Header().Add("Content-Type", "application/json")
		log.Println(w.Header().Get("Content-Type"))
	}
	if m, _ := regexp.MatchString("/menu.html", filePath); m {
		t := template.New("menu.html")
		t.ParseFiles(filePath)
		t.Execute(w, socketPath)
	} else {
		http.ServeFile(w, r, filePath)
	}
}

//websocket请求处理逻辑
func DealSocketReq(ws *websocket.Conn) {
	wsconn.PutConn(ws)
	var err error
	websocket.Message.Send(ws, "websocket Connection successfull!")
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}
		fmt.Println("Received back from client: " + reply)

		//msg := "Received from " + ws.Request().Host + "  " + reply
		msg := "welcome to websocket do by pp"
		fmt.Println("Sending to client: " + msg)
		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}
