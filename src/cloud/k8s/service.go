package k8s

import (
	"k8s.io/client-go/kubernetes"
	"github.com/astaxie/beego/logs"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"sort"
	"encoding/json"
	"cloud/util"
)

// 获取某个集群的服务信息
func GetServices(clientset kubernetes.Clientset, namespace string) ([]v1.Service, error) {
	opt := metav1.ListOptions{}
	data, err := clientset.CoreV1().Services(namespace).List(opt)
	if err != nil {
		logs.Error("获取service失败啦", err)
		return make([]v1.Service, 0), err
	}
	json.Marshal(data.Items)
	return data.Items, err
}

// 2018-02-13 19:27
// 获取服务信息
func GetAppService(clientset kubernetes.Clientset, namespace string, serviceName string) *v1.Service  {
	data, err := clientset.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{})
	if err == nil {
		return data
	}
	return &v1.Service{}
}

// 删除某个service
func DeleteService(clustername string,  namespace string,name string) (error) {
	cl, err := GetYamlClient(clustername, "", "v1", "/api")
	resource := &metav1.APIResource{Name: "Services", Namespaced: true}
	opt := metav1.DeleteOptions{}
	err = cl.Resource(resource, namespace).Delete(name, &opt)
	if err != nil {
		logs.Error("删除 Service 失败", err)
		return err
	}
	logs.Info("删除 Service 成功")
	return err
}

// 获取某个集群服务的数量
func GetServiceNumber(clientset kubernetes.Clientset,namespace string) int {
	data, _ := GetServices(clientset,namespace)
	return len(data)
}

// 获取自定义的namespace的service
func GetCusomService(clientset kubernetes.Clientset, namespace string) []string {
	data,_ := GetServices(clientset, namespace)
	servers := make([]string,0)
	for _, k := range data{
		// 自定义的namespace都包含 --
		if strings.Contains(k.Namespace, "--"){
			servers = append(servers, k.Name)
		}
	}
	return servers
}

// 2018-02-20 11:29
// 获取服务端口
func GetServerPort(clientset kubernetes.Clientset, namespace string, name string)  util.Lock{
	data, _ := clientset.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	temp := util.Lock{}
	for _, v := range data.Spec.Ports {
		temp.Put("nodePort", v.NodePort)
		temp.Put("name", v.Name)
		temp.Put("port", v.Port)
		temp.Put("targetPort", v.TargetPort)
		temp.Put("protocol", v.Protocol)
		return temp
	}
	return temp
}

// 2018-01-14 17:52
// 获取当前服务使用的端口,和集群地址
// 在更新service的时候使用
func GetCurrentPort(clientset kubernetes.Clientset, namespace string, name string) util.Lock {
	all := util.Lock{}
	data, err := clientset.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("获取服务端口错误", err)
	}
	if err == nil {
		for _, v := range data.Spec.Ports {
			temp := util.Lock{}
			temp.Put("nodePort", v.NodePort)
			temp.Put("name", v.Name)
			temp.Put("port", v.Port)
			temp.Put("targetPort", v.TargetPort)
			temp.Put("protocol", v.Protocol)
			all.Put(v.TargetPort.String(), temp)
		}
		all.Put("resourceVersion", data.ResourceVersion)
		all.Put("clusterIp", data.Spec.ClusterIP)
	}
	return all
}

// 获取某个服务使用的端口
// 2018-01-21 18:03
func GetServicePort(clientset kubernetes.Clientset, namespace string, name string) *v1.Service {
	opt := metav1.GetOptions{}
	data, err := clientset.CoreV1().Services(namespace).Get(name, opt)
	if err == nil {
		return data
	}
	return &v1.Service{}
}

// 获取k8s svc 已经使用的端口
// 应该任务计划去执行,数据放到库里
// 端口范围默认Wie
func GetServiceUsedPort(clientset kubernetes.Clientset) []int {
	data,_ := GetServices(clientset,"")
	usedPort := make([]int,0)
	for _, k := range data{
		servicePort := k.Spec.Ports
		for _, s := range servicePort {
			if s.NodePort >= 30000 {
				usedPort = append(usedPort, int(s.NodePort))
				usedPort = append(usedPort, int(s.Port))
			}
		}
	}
	// grafana监控
	usedPort = append(usedPort, 43000)
	// 给Glusterfs预留
	usedPort = append(usedPort, 48080)
	return usedPort
}

// 获取一个端口给service使用
// 获取最大不超过65535
func GetServiceFreePort(clientset kubernetes.Clientset)int32 {
	ports := GetServiceUsedPort(clientset)
	size := len(ports)
	if size < 1{
		return 30001
	}
	sort.Ints(ports)
	return int32(ports[size-1] + 1)
}


// 获取可用的端口,一次获取多个
// 2018-01-14 14:29
func GetServicePorts(clientset kubernetes.Clientset, size int, start int, end int) []int {
	ports := GetServiceUsedPort(clientset)
	logs.Info("获取到服务器已经使用的端口", ports)
	all := make([]int, 0)
	if start == 0 {
		start = 30000
	}
	if end == 0 {
		end = 49000
	}
	for i := start; i <= end; i++  {
		if !util.ListExistsInt(ports, i){
			all = append(all, i)
			if len(all) > size {
				return all
			}
		}
	}
	return all
}