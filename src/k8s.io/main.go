package main

import (
	"log"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/dynamic"
	//"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	//"fmt"
	"fmt"
	"io/ioutil"
	"cloud/util"
)


func main() {
	log.SetFlags(log.Llongfile)
	//flag.Parse()
	////获取Config
	//
	//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	//if err != nil {
	//	log.Println(err)
	//}
	config := restclient.Config{}
	config.Host = "http://10.16.55.6:8080"
	//指定gv
	gv := &schema.GroupVersion{"apps", "v1beta1"}
	//指定resource
	//resource := &v12.APIResource{Name: "pods", Namespaced: true}
	resource := &v12.APIResource{Name: "Deployments", Namespaced: true}

	//指定GroupVersion
	config.ContentConfig = rest.ContentConfig{GroupVersion: gv}
	//默认的是/api 需要手动指定
	config.APIPath = "/apis"
	//创建新的dynamic client
	cl, err := dynamic.NewClient(&config)
	if err != nil {
		log.Println(err)
	}

	//根据APIResource获取
	//
	//obj, err := cl.Resource(resource, "default").Get("zhaoyun-rc-2fb49",v12.GetOptions{})
	//if err != nil {
	//	log.Println(err)
	//}
	//pod := v1.Pod{}
	//b, err := json.Marshal(obj.Object)
	//if err != nil {
	//	log.Println(err)
	//}
	//json.Unmarshal(b, &pod)
	//log.Println(pod.Name)

	//创建pod
	//conf := make(map[string]interface{})
	//conf = map[string]interface{}{
	//	"apiVersion": "v1",
	//	"kind":       "Pod",
	//	"metadata": map[string]interface{}{
	//		"name": "golang1",
	//	},
	//	"spec": map[string]interface{}{
	//		"containers": []map[string]interface{}{
	//			map[string]interface{}{
	//				"image": "golang",
	//				"command": []string{
	//					"sleep",
	//					"3600",
	//				},
	//				"name": "golang1",
	//			},
	//		},
	//	},
	//}
	y,_ := ioutil.ReadFile("D:\\F\\1.yaml")
	conf,err := util.Yaml2Json(y)
	podobj := unstructured.Unstructured{Object: conf}
	data, err := cl.Resource(resource, "default").Create(&podobj)
	fmt.Println(data.Object)
	if err != nil {
		log.Println(err)
	}
	//// 删除一个pod,删除资源前最好获取UUID
	err = cl.Resource(resource, "default").Delete("zhaoyun-rc-2fb49", &metav1.DeleteOptions{})
	fmt.Println("删除信息...", err)

	//// 获取列表
	//got, err := cl.Resource(resource, "default").List(metav1.ListOptions{})
	//if err != nil {
	//	log.Println(err)
	//}
	//js, err := json.Marshal(reflect.ValueOf(got).Elem().Interface())
	//if err != nil {
	//	log.Println(err)
	//}
	//podlist := v1.PodList{}
	//err = json.Unmarshal(js, &podlist)
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println(podlist.Items[0].Name)
	//
	//// 获取thirdpart resource
	//gvthird := &schema.GroupVersion{"test.io", "v1"}
	//thirdpartresource := metav1.APIResource{Name: "podtoservices", Namespaced: true}
	//config.ContentConfig = rest.ContentConfig{GroupVersion: gvthird}
	//config.APIPath = "/apis"
	//clthird, err := dynamic.NewClient(&config)
	//if err != nil {
	//	log.Println(err)
	//}
	//objthird, err := clthird.Resource(&thirdpartresource, "default").Get("redis-slave-360xf",  metav1.GetOptions{})
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println(objthird)
	//
	//
	//
	////watch一个resource
	//watcher, err := clthird.Resource(&thirdpartresource, "").Watch(metav1.ListOptions{})
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//c := watcher.ResultChan()
	//for {
	//	select {
	//	case e := <-c:
	//		getptrstring(e)
	//	}
	//}
}
