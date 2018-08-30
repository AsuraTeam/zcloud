package k8s

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"cloud/sql"
	"time"
	"strings"
	"github.com/garyburd/redigo/redis"
	"cloud/cache"
	"cloud/util"
	"github.com/astaxie/beego/logs"
)

const SelectCloudClusterHosts = "select host_ip,host_type,cluster_name,api_port from cloud_cluster_hosts"

func GetMasterIp(cluster string) (string, string) {
	//logs.Info(masterCache, masterErr, cluster)
	if cache.MasterCache != nil {
		r := cache.MasterCache.Get(cluster)
		if r != nil {
			masterData, _ := redis.String(r, nil)
			if strings.Contains(masterData, ",") {
				masterDatas := strings.Split(masterData, ",")
				if len(masterDatas) > 1 {
					if len(masterDatas[0]) > 0 {
						return masterDatas[0], masterDatas[1]
					}
				}
			}
		}
	}

	data := make([]CloudClusterHosts, 0)
	searchMap := sql.SearchMap{}
	searchMap.Put("ClusterName", cluster)
	searchMap.Put("HostType", "master")
	q := sql.SearchSql(CloudClusterHosts{}, SelectCloudClusterHosts, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, d := range data {
		if d.HostType == "master" {
			if cache.MasterCache != nil {
				cache.MasterCache.Put(cluster, d.HostIp+","+d.ApiPort, time.Second*600)
			}
			return d.HostIp, d.ApiPort
		}
	}
	return "", ""
}

// 2018-03-01 14:24
// 获取证书信息
func getCertFile(name string) CertData {
	data := CertData{}
	q := "select ca_data,cert_data,key_data from cloud_cluster where cluster_name=?"
	sql.GetOrm().Raw(q, name).QueryRow(&data)
	return data
}

// 2018-03-01 14:50
// 获取客户端证书配置
func getTnlCfg(cluster string) restclient.Config {
	master, port := GetMasterIp(cluster)
	logs.Info("获取集群地址", master, port, cluster)
	caData := getCertFile(cluster)
	config := restclient.Config{}
	config.CAData = []byte(caData.CaData)
	tlsCfg := restclient.TLSClientConfig{
		Insecure: false,
		CAData:config.CAData,
		KeyData:[]byte(caData.KeyData),
		CertData:[]byte(caData.CertData),
	}
	config.TLSClientConfig = tlsCfg
	port = "6443"
	config.Host = "https://" + master + ":" + port
	return config
}

// client信息缓存
var clientPool = util.Lock{}
func GetClient(cluster string) (kubernetes.Clientset, error) {
	key := cluster + "clientSet"
	c, ok := clientPool.Get(key)
	if ok && c != nil {
		return c.(kubernetes.Clientset), nil
	}

	config := getTnlCfg(cluster)

	config.Timeout = time.Second * 3
	client, err := kubernetes.NewForConfig(&config)
	if err != nil {
		logs.Error("GetClient Error", err.Error())
		return kubernetes.Clientset{}, err
	}
	clientPool.Put(key, *client)
	return *client, err
}

// 通过yaml方式部署服务
func GetYamlClient(cluster string, groups string, version string, api string) (*dynamic.Client, error) {
	config := getTnlCfg(cluster)
	//指定gv
	// k8s 1.8.2
	gv := &schema.GroupVersion{groups, version}
	//指定resource
	//resource := &v12.APIResource{Name: "pods", Namespaced: true}
	//指定GroupVersion
	config.ContentConfig = rest.ContentConfig{GroupVersion: gv}
	config.Timeout = time.Second * 3
	//默认的是/api 需要手动指定
	config.APIPath = api
	//创建新的dynamic client
	cl, err := dynamic.NewClient(&config)
	return cl, err
}

// 2018-02-28 09:26
// 获取用来执行命令和websocke使用的client
func GetRestlient(cluster string) (*rest.RESTClient, restclient.Config, error) {
	groupversion := schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	config := getTnlCfg(cluster)
	config.GroupVersion = &groupversion
	config.APIPath = "/api"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	restclient, err := rest.RESTClientFor(&config)
	return restclient, config, err
}
