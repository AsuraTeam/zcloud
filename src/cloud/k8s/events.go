package k8s

import (
	"k8s.io/client-go/kubernetes"
	"github.com/astaxie/beego/logs"
	"cloud/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

// 2018-02-27 21:01
// 分析事件信息
func parseEvents(items []v1.Event) []EventData {
	data := make([]EventData, 0)
	for _, v := range items {
		t := EventData{}
		t.Reason = v.Reason
		t.Messages = v.Message
		t.EventTime = util.ReplaceTime(v.LastTimestamp.String())
		t.Host = v.Source.Host
		t.Type = v.Type
		data = append(data, t)
	}
	return data
}

// 2018-02-27 20:54
// 获取容器事件信息
func GetEvents(namespace string, podName string, clientSet kubernetes.Clientset) []EventData {
	opt := metav1.ListOptions{}
	opt.FieldSelector = "involvedObject.name="+podName+",involvedObject.namespace=" + namespace
	events, err := clientSet.CoreV1().Events(namespace).List(opt)
	if err != nil {
		logs.Error("获取Pods错误", err.Error())
	}
	return parseEvents(events.Items)
}
