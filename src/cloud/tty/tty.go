package tty

import (
	//"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	//"strings"
	"io"
	//"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"fmt"
	"strings"
	"cloud/util"
	"strconv"
	"time"
	"cloud/k8s"
)

// 获取url的参数
// 2017-01-15 10:29
func getParam(param []string, key string) string {
	for _, v := range param {
		vs := strings.Split(v, "=")
		if vs[0] == key && len(vs) == 2 {
			return vs[1]
		}
	}
	return ""
}

func Handler(r io.Reader, w io.Writer, containername string, podname string, namespace string, clustername string) {
	restclient,config, err := k8s.GetRestlient(clustername)
	if err != nil {
		log.Fatalln(err)
	}

	req := restclient.Post().
		Resource("pods").
		Name(podname).
		Namespace(namespace).
		SubResource("exec").Timeout(time.Second * 1)
	names := strings.Split(podname, "--")
	if len(names) > 1 {
		version := strings.Split(strings.Split(podname, "--")[1], "-")[0]
		containername = util.Namespace(names[0], version)
	}
	req.VersionedParams(
		&v1.PodExecOptions{
			Container: containername,
			Command:   []string{"sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		},
		scheme.ParameterCodec,
	)

	log.Println(req.URL().String())
	executor, err := remotecommand.NewSPDYExecutor(
		&config, http.MethodPost, req.URL(),
	)

	if err != nil {
		log.Println(err)
	}
	//strings.NewReader("touc /aa.txt")
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:  r,
		Stdout: w,
		Stderr: w,
		Tty:    true,
		//TerminalSizeQueue:  Terminal.s,
	})

	if err != nil {
		log.Println(err)
	}
}

//Handler(ws, ws,"auto-3","auto-3-7b45c8757c-bngtf","auto-3--dfsad")
func TtyHandler(ws *websocket.Conn) {
	p := strings.Split(ws.Request().RequestURI, "&")
	pod := getParam(p, "pod")
	container := getParam(p, "container")
	namespace := getParam(p, "namespace")
	token := getParam(p, "token")
	username := getParam(p, "username")
	timestamp := getParam(p, "timestamp")
	cluster := getParam(p, "cluster")
	fmt.Println(pod, token, namespace, container)

	d := make([]string, 0)
	d = append(d, username)
	d = append(d, namespace)
	d = append(d, pod)
	d = append(d, container)
	d = append(d, timestamp)
	d = append(d, cluster)

	times,err := strconv.ParseInt(timestamp,10, 64)
	if err != nil {
		ws.Write([]byte("时间戳格式错误"))
		return
	}

	if time.Now().Unix() - times > 1800 {
		ws.Write([]byte("该会话已经超时"))
		return
	}

	pass := beego.AppConfig.String("ttysecurity")
	tokenSecurity := util.Md5String(strings.Join(d, pass))
	if pod == "" || container == "" || namespace == "" || token == "" || username == "" || cluster == "" {
		ws.Write([]byte("参数错误"))
		return
	}
	if tokenSecurity != token {
		ws.Write([]byte("验证失败"))
		return
	}
	Handler(ws, ws, container, pod, namespace, cluster)
}

func TtyStart() {
	http.Handle("/tty", websocket.Handler(TtyHandler))
	addr := beego.AppConfig.String("ttyport")
	if addr == "" {
		addr = "8999"
	}
	logs.Info("web terminal port:", addr)
	err := http.ListenAndServe(":"+addr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
