package k8s

import (
	"strings"
	"k8s.io/client-go/kubernetes"
	"github.com/astaxie/beego/logs"
	"fmt"
	"cloud/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

// 获取namespace
func GetNamespaces(clientset kubernetes.Clientset) ([]v1.Namespace, error) {
	opt := metav1.ListOptions{}
	data, err := clientset.CoreV1().Namespaces().List(opt)
	if err != nil {
		logs.Error("获取namespaces失败啦", err)
		return make([]v1.Namespace, 0), err
	}
	return data.Items, err
}

// 获取自己创建的namespace应用, 规则是app名加资源名区分
func GetNamespaceApp(clientset kubernetes.Clientset) []CloudApp {
	namespaces, _ := GetNamespaces(clientset)
	fmt.Println(namespaces)
	datas := make([]CloudApp, 0)
	for _, v := range namespaces {
		data := CloudApp{}
		data.Status = string(v.Status.Phase)
		data.CreateTime = util.ReplaceTime(v.CreationTimestamp.String())

		name := strings.Split(v.Name, "--")
		if len(name) < 2 {
			continue
		}
		pods := GetPods(v.Name, clientset)
		data.ContainerNumber = len(pods)
		if len(pods) > 0 {
			data.LastUpdateTime = util.ReplaceTime(pods[0].Status.StartTime.String())
		}else{
			data.LastUpdateTime = "未知"
		}
		data.AppLabels = "无"
		data.AppName = name[0]
		data.ClusterName = name[1]
		datas = append(datas, data)
	}
	return datas
}
