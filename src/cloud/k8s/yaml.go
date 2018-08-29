package k8s

import (
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"github.com/astaxie/beego/logs"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"cloud/util"
	"k8s.io/client-go/dynamic"

	"encoding/json"
)


// 通过yaml文件创建服务

func CreateService() {

}

//apiVersion: apps/v1beta1
//kind: Deployment
//metadata:
//  name: deploymentexample
//spec:
//  replicas: 3
//  template:
//    metadata:
//     labels:
//        app: nginx
//    spec:
//      containers:
//      - name: nginx
//        image: nginx:1.10
// 通过yaml创建Deployment 服务
func YamlDeployment(clustername string, yaml []byte, namespace string, isService bool, uuid string) error {
	cl, err := GetYamlClient(clustername, "apps", "v1beta1", "/apis")
	resource := &v12.APIResource{Name: "Deployments", Namespaced: true}
	conf, err := util.Yaml2Json(yaml)

	podobj := unstructured.Unstructured{Object: conf}

	labels := podobj.GetLabels()
	if len(labels) < 1{
		labels = make(map[string]string)
	}
	labels["zcloud-app"] = namespace
	labels["space"] = namespace
	labels["release-version"] = "0"
	labels["uuid"] = uuid
	podobj.SetLabels(labels)


	d, err := cl.Resource(resource, namespace).Create(&podobj)
	if err != nil {
		logs.Error("创建yaml的deployment服务失败", err)
		return err
	}
	logs.Info("创建yaml的deployment服务成功", d)
	if !isService {
		logs.Info("该 deployment 不需要创建 Service ... ")
		return err
	}
	logs.Info("开始创建 Service ...")
	client, err := GetClient(clustername)
	if err != nil {
		return err
	}
	serviceData := GetPodsFromUUid(namespace, labels["uuid"], client)
	fmt.Println(serviceData)
	cl2, _ := GetYamlClient(clustername, "", "v1", "api")
	for _, s := range serviceData {
		for _, port := range s.ContainerPort {
			nodePort := GetServiceFreePort(client)
			logs.Info("获取到未使用 nodeport 未", nodePort)
			YamlCreateService(s.AppName, s.ResouceName, port, cl2, nodePort, s.Selector, namespace)
		}
	}
	return err
}

// 创建namespace
func YamlCreateNamespace(clustername string, namespace string) error {
	cl, err := GetYamlClient(clustername, "", "v1", "/api")
	yaml := "apiVersion: v1\n" +
		"kind: Namespace\n" +
		"metadata:\n" +
		"   name: {0}\n" +
		"   labels:\n" +
		"      name: {0}"
	yaml = strings.Replace(yaml, "{0}", namespace, -1)
	conf, err := util.Yaml2Json([]byte(yaml))
	fmt.Println(err)
	resource := &v12.APIResource{Name: "Namespaces", Namespaced: false}
	podobj := unstructured.Unstructured{Object: conf}
	d, err := cl.Resource(resource, namespace).Create(&podobj)
	if err != nil {
		logs.Error("创建yaml的namespace服务失败", err)
		return err
	}
	logs.Info("创建yaml的namespace服务成功", d)
	return err
}

// 通过yaml创建DaemonSet服务
func YamlDaemonSet(master string, port string, yaml []byte, namespace string) {
}

// 通过yaml创建StatefulSets服务
func YamlStatefulSets(master string, port string, yaml []byte, namespace string) {
}

// 部署一个nginx的负载均衡服务器
// 使用容器方式运行nginx,并使用k8s管理
// 在容器内部跑agent，去配置nginx，主要
// 一个集群一个nginx,或haproxy提供服务,内部服务注册发现的不需要配置
// 信息放到redis里面，agent监控redis变化实现配置更新
// 由master服务往redis里面放置数据
// 1、对k8s service
func YamlLbNgxin() {

}

// 部署一个haproxy的负载均衡
func YamlLbHaproxy() {

}


// --service-node-port-range=20000-65535
// NodePort 对外部可见的
// 创建一个service, 在应用创建完成后,自动创建一个service
// 通过配置开关来定义是否创建service，如果打开就创建,主要用于http非服务化的服务
// c,_ := k8s.GetYamlClient("10.16.55.6","8080", "","v1", "api")
// b := k8s.YamlCreateService("my-nginx","dddd",80, c, 50000, "cccccc")
//  uuid 为创建pods时候生成的
func YamlCreateService(appname string, resourceName string, containerPort int32, cl *dynamic.Client, nodePort int32, selector map[string]string, namespace string) error {
	name := appname + "--" + resourceName
	conf := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]interface{}{
			"name": name,
			"labels": map[string]interface{}{
				"app": name,
			},
		},
		"spec": map[string]interface{}{
			"type": "NodePort",
			"selector": selector,
			"ports": []map[string]interface{}{
				map[string]interface{}{
					"port":       nodePort,
					"targetPort": containerPort,
					"nodePort":   nodePort,
				},
			},
		},
	}
	d1 ,_ := json.Marshal(conf)
	fmt.Println(string(d1))
	resource := &v12.APIResource{Name: "Services", Namespaced: true}
	obj := unstructured.Unstructured{Object: conf}
	d, err := cl.Resource(resource, namespace).Create(&obj)
	if err != nil {
		logs.Error("创建yaml的service失败", err)
		return err
	}
	logs.Info("创建yaml的service服务成功", d)
	return err
}