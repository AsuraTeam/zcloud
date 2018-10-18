package k8s

import (
	"io/ioutil"
	"time"
	"net/http"
	"strings"
	"path/filepath"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"github.com/astaxie/beego/logs"
	"os"
	"cloud/util"
	"k8s.io/client-go/rest"
)

var restclientLock = util.Lock{}

// 2018-02-27 17:35
// 执行命令
func Exec(clustername string, podname string, namespace string, containername string, cmd []string) string {
	key := clustername
	configKey := key + "config"
	var restclient *rest.RESTClient
	var config rest.Config
	var err error
	if _, ok := restclientLock.Get(key); ok {
		config = restclientLock.GetV(configKey).(rest.Config)
		restclient = restclientLock.GetV(key).(*rest.RESTClient)
		logs.Info("从缓存中获取连接")
	} else {
		restclient, config, err = GetRestlient(clustername)
		restclientLock.Put(key, restclient)
		restclientLock.Put(configKey, config)
	}
	if err != nil {
		logs.Error("执行docker命令失败", err)
		return ""
	}
	names := strings.Split(containername, "--")
	if len(names) > 1 {
		version := strings.Split(strings.Split(podname, "--")[1], "-")[0]
		containername = util.Namespace(names[0], version)
	}
	logs.Info("执行容器命令", containername, util.ObjToString(cmd))
	req := restclient.Post().
		Resource("pods").
		Name(podname).
		Namespace(namespace).
		SubResource("exec").Timeout(time.Second * 1)

	req.VersionedParams(
		&v1.PodExecOptions{
			Container: containername,
			Command:   cmd,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		},
		scheme.ParameterCodec,
	)

	executor, err := remotecommand.NewSPDYExecutor(
		&config, http.MethodPost, req.URL(),
	)

	if err != nil {
		logs.Error("NewSPDYExecutor", err)
		return ""
	}
	var dir string
	cwd, _ := os.Getwd()
	if strings.Contains(cwd, "/") {
		dir = "/dev/shm/"
	} else {
		dir = "./"
	}
	if dir == "" {
		return ""
	}
	file := filepath.Join(dir, namespace+podname+containername)
	r, err := os.Create(file)
	w := r
	c1 := make(chan error)
	go func() {
		c1 <- executor.Stream(remotecommand.StreamOptions{
			Stdin:  r,
			Stdout: w,
			Stderr: w,
			Tty:    false,
		})
	}()
	select {
	case err := <-c1:
		logs.Info(err)
		break
	case <-time.After(time.Second * 1):
		logs.Info("执行超时")
		break
	}
	w.Close()
	r.Close()
	t, err := ioutil.ReadFile(file)
	os.Remove(file)
	if err == nil {
		return string(t)
	}
	return ""
}
