package hosts

import (
	"cloud/models/hosts"
	"k8s.io/client-go/kubernetes"
	"cloud/k8s"
	"github.com/astaxie/beego/logs"
	"time"
	"cloud/sql"
	"cloud/util"
	"k8s.io/api/core/v1"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"cloud/cache"
)

// 2018-02-13 09:30
// 将数据库有的数据写入到缓存
// 没有的写入到数据库
func writeHostToDb(lock util.Lock, host v1.Node, cl kubernetes.Clientset, nodeData []v1.Node, v hosts.CloudClusterHosts)  {
	hostDataStr, ok := lock.Get(host.Name)
	if ok {
		setHostCache(hostDataStr.(hosts.CloudClusterHosts), cl, nodeData)
	} else {
		// 新的主机可能没在数据库,写入到数据库
		hostD := hosts.CloudClusterHosts{}
		hostD.HostIp = host.Name
		hostD.ClusterName = v.ClusterName
		hostD.HostType = "slave"
		hostD.ApiPort = "0"
		q := sql.InsertSql(hostD, hosts.InsertCloudClusterHosts)
		sql.Raw(q).Exec()
	}
}

var CacheLock = util.Lock{}

func CacheNodeStatus(data []hosts.CloudClusterHosts, c kubernetes.Clientset) {
	if !util.WriteLock("last_update", &CacheLock, 10) {
		return
	}
	nodeData := k8s.GetNodes(c, "")
	logs.Info(util.ObjToString(nodeData))
	for _, hostData := range data {
		setHostCache(hostData, c, nodeData)
	}
}

// 2018-02-12 21:34
// 写入redis
func setHostCache(hostData hosts.CloudClusterHosts, c kubernetes.Clientset, nodeData []v1.Node) {
	nodeStatus := k8s.GetNodesFromIp(hostData.HostIp, c, nodeData)
	nodeStatus.HostType = hostData.HostType
	nodeStatus.HostId = hostData.HostId
	nodeStatus.HostLabel = hostData.HostLabel
	nodeStatus.HostIp = hostData.HostIp
	nodeStatus.HostType = hostData.HostType

	if cache.HostCacheErr == nil {
		cache.HostCache.Put(hostData.HostIp, util.ObjToString(nodeStatus), time.Second*86400*5)
	}
}

// 获取单个主机的信息
// 2018-01-18 21:11
func getHostData(id interface{}) hosts.CloudClusterHosts {
	searchMap := sql.SearchMap{}
	searchMap.Put("HostId", id)
	cloudCluster := hosts.CloudClusterHosts{}
	q := sql.SearchSql(
		cloudCluster,
		hosts.SelectCloudClusterHosts,
		searchMap)
	sql.Raw(q).QueryRow(&cloudCluster)
	return cloudCluster
}



// 将数据写入到map,避免重复查库
func setHostsMap(data []hosts.CloudClusterHosts) util.Lock {
	lock := util.Lock{}
	for _, v := range data {
		lock.Put(v.HostIp, v)
	}
	return lock
}

// 2018-02-12 21;37
// 任务计划设置缓存
func CronCache() {
	logs.Info("开始写入Node缓存")
	data := make([]hosts.CloudClusterHosts, 0)
	q := sql.SearchSql(
		hosts.CloudClusterHosts{},
		hosts.SelectCloudClusterHosts, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)

	lock := setHostsMap(data)
	for _, v := range data {
		if v.HostType == "master" {
			cl, _ := k8s.GetClient(v.ClusterName)
			nodeData := k8s.GetNodes(cl, "")
			for _, host := range nodeData {
				writeHostToDb(lock, host, cl, nodeData, v)
			}
		}
	}
}


// 获取某个集群里面集群的数量,删除集群时做验证
func GetClusterHosts(cluster string) []hosts.CloudClusterHosts {
	data := make([]hosts.CloudClusterHosts, 0)
	searchMap := sql.SearchMap{}
	searchMap.Put("ClusterName", cluster)
	searchSql := sql.SearchSql(
		hosts.CloudClusterHosts{},
		hosts.SelectCloudClusterHosts,
		searchMap)
	sql.Raw(searchSql).QueryRows(&data)
	return data
}

// 2018-02-13 09:37
// 将redis的数据读取出来
func getRedisNodeData(data []hosts.CloudClusterHosts) []k8s.NodeStatus  {
	returnData := make([]k8s.NodeStatus, 0)
	for _, hostData := range data {
		r := cache.HostCache.Get(hostData.HostIp)
		if r != nil {
			redisR, err := redis.String(r, nil)
			if err == nil {
				nodeStatus := k8s.NodeStatus{}
				json.Unmarshal([]byte(redisR), &nodeStatus)
				returnData = append(returnData, nodeStatus)
			}
		}
	}
	return returnData
}

//获取集群里master的地址
func GetMaster(cluster string) (string, string) {
	return k8s.GetMasterIp(cluster)
}


func getHostUser(this *HostsController) string {
	return util.GetUser(this.GetSession("username"))
}

func setHostJson(this *HostsController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}
